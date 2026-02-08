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

type ReversalSubmissionHandler struct {
	store  store.Store
	engine *lifecycle.Engine
}

func NewReversalSubmissionHandler(s store.Store, e *lifecycle.Engine) *ReversalSubmissionHandler {
	return &ReversalSubmissionHandler{store: s, engine: e}
}

func (h *ReversalSubmissionHandler) Create(w http.ResponseWriter, r *http.Request) {
	paymentID := r.PathValue("paymentID")
	reversalID := r.PathValue("reversalID")

	if _, err := h.store.GetReversal(paymentID, reversalID); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			jsonapi.NotFound(w, "reversal", reversalID)
			return
		}
		jsonapi.InternalError(w)
		return
	}

	var req jsonapi.DataEnvelope[models.ReversalSubmission]
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonapi.BadRequest(w, "invalid JSON: "+err.Error())
		return
	}

	s := req.Data
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	s.Type = models.ResourceTypeReversalSubmission
	now := time.Now().UTC()
	s.CreatedOn = now
	s.ModifiedOn = now
	s.Attributes.Status = lifecycle.SimpleSubmissionChain[0]
	s.Attributes.SubmissionDate = now.Format(time.DateOnly)

	if err := h.store.CreateReversalSubmission(paymentID, reversalID, s); err != nil {
		if errors.Is(err, store.ErrConflict) {
			jsonapi.Conflict(w, "reversal submission already exists")
			return
		}
		jsonapi.InternalError(w)
		return
	}

	submissionID := s.ID
	h.engine.StartTransition(models.ResourceTypeReversalSubmission, submissionID, lifecycle.SimpleSubmissionChain, func(newStatus string) error {
		sub, err := h.store.GetReversalSubmission(paymentID, reversalID, submissionID)
		if err != nil {
			return err
		}
		sub.Attributes.Status = newStatus
		sub.ModifiedOn = time.Now().UTC()
		return h.store.UpdateReversalSubmission(paymentID, reversalID, sub)
	})

	w.Header().Set("Content-Type", jsonapi.ContentType)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(jsonapi.DataEnvelope[models.ReversalSubmission]{Data: s})
}

func (h *ReversalSubmissionHandler) Get(w http.ResponseWriter, r *http.Request) {
	paymentID := r.PathValue("paymentID")
	reversalID := r.PathValue("reversalID")
	submissionID := r.PathValue("submissionID")

	s, err := h.store.GetReversalSubmission(paymentID, reversalID, submissionID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			jsonapi.NotFound(w, "reversal_submission", submissionID)
			return
		}
		jsonapi.InternalError(w)
		return
	}

	w.Header().Set("Content-Type", jsonapi.ContentType)
	json.NewEncoder(w).Encode(jsonapi.DataEnvelope[models.ReversalSubmission]{Data: s})
}
