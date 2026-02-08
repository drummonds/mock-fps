package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/nibble/mock-fps/internal/jsonapi"
	"github.com/nibble/mock-fps/internal/lifecycle"
	"github.com/nibble/mock-fps/internal/models"
	"github.com/nibble/mock-fps/internal/store"
)

type RecallSubmissionHandler struct {
	store  store.Store
	engine *lifecycle.Engine
}

func NewRecallSubmissionHandler(s store.Store, e *lifecycle.Engine) *RecallSubmissionHandler {
	return &RecallSubmissionHandler{store: s, engine: e}
}

func (h *RecallSubmissionHandler) Create(w http.ResponseWriter, r *http.Request) {
	paymentID := r.PathValue("paymentID")
	recallID := r.PathValue("recallID")

	if _, err := h.store.GetRecall(paymentID, recallID); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			jsonapi.NotFound(w, "recall", recallID)
			return
		}
		jsonapi.InternalError(w)
		return
	}

	var req jsonapi.DataEnvelope[models.RecallSubmission]
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonapi.BadRequest(w, "invalid JSON: "+err.Error())
		return
	}

	s := req.Data
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	s.Type = models.ResourceTypeRecallSubmission
	now := time.Now().UTC()
	s.CreatedOn = now
	s.ModifiedOn = now
	s.Attributes.Status = lifecycle.SimpleSubmissionChain[0]
	s.Attributes.SubmissionDate = now.Format(time.DateOnly)

	if err := h.store.CreateRecallSubmission(paymentID, recallID, s); err != nil {
		if errors.Is(err, store.ErrConflict) {
			jsonapi.Conflict(w, "recall submission already exists")
			return
		}
		jsonapi.InternalError(w)
		return
	}

	submissionID := s.ID
	h.engine.StartTransition(models.ResourceTypeRecallSubmission, submissionID, lifecycle.SimpleSubmissionChain, func(newStatus string) error {
		sub, err := h.store.GetRecallSubmission(paymentID, recallID, submissionID)
		if err != nil {
			return err
		}
		sub.Attributes.Status = newStatus
		sub.ModifiedOn = time.Now().UTC()
		return h.store.UpdateRecallSubmission(paymentID, recallID, sub)
	})

	w.Header().Set("Content-Type", jsonapi.ContentType)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(jsonapi.DataEnvelope[models.RecallSubmission]{Data: s})
}

func (h *RecallSubmissionHandler) Get(w http.ResponseWriter, r *http.Request) {
	paymentID := r.PathValue("paymentID")
	recallID := r.PathValue("recallID")
	submissionID := r.PathValue("submissionID")

	s, err := h.store.GetRecallSubmission(paymentID, recallID, submissionID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			jsonapi.NotFound(w, "recall_submission", submissionID)
			return
		}
		jsonapi.InternalError(w)
		return
	}

	w.Header().Set("Content-Type", jsonapi.ContentType)
	json.NewEncoder(w).Encode(jsonapi.DataEnvelope[models.RecallSubmission]{Data: s})
}
