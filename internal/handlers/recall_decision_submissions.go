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

type RecallDecisionSubmissionHandler struct {
	store  store.Store
	engine *lifecycle.Engine
}

func NewRecallDecisionSubmissionHandler(s store.Store, e *lifecycle.Engine) *RecallDecisionSubmissionHandler {
	return &RecallDecisionSubmissionHandler{store: s, engine: e}
}

func (h *RecallDecisionSubmissionHandler) Create(w http.ResponseWriter, r *http.Request) {
	paymentID := r.PathValue("paymentID")
	recallID := r.PathValue("recallID")
	decisionID := r.PathValue("decisionID")

	if _, err := h.store.GetRecallDecision(paymentID, recallID, decisionID); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			jsonapi.NotFound(w, "recall_decision", decisionID)
			return
		}
		jsonapi.InternalError(w)
		return
	}

	var req jsonapi.DataEnvelope[models.RecallDecisionSubmission]
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonapi.BadRequest(w, "invalid JSON: "+err.Error())
		return
	}

	s := req.Data
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	s.Type = models.ResourceTypeRecallDecisionSubmission
	now := time.Now().UTC()
	s.CreatedOn = now
	s.ModifiedOn = now
	s.Attributes.Status = lifecycle.SimpleSubmissionChain[0]
	s.Attributes.SubmissionDate = now.Format(time.DateOnly)

	if err := h.store.CreateRecallDecisionSubmission(paymentID, recallID, decisionID, s); err != nil {
		if errors.Is(err, store.ErrConflict) {
			jsonapi.Conflict(w, "recall decision submission already exists")
			return
		}
		jsonapi.InternalError(w)
		return
	}

	submissionID := s.ID
	h.engine.StartTransition(models.ResourceTypeRecallDecisionSubmission, submissionID, lifecycle.SimpleSubmissionChain, func(newStatus string) error {
		sub, err := h.store.GetRecallDecisionSubmission(paymentID, recallID, decisionID, submissionID)
		if err != nil {
			return err
		}
		sub.Attributes.Status = newStatus
		sub.ModifiedOn = time.Now().UTC()
		return h.store.UpdateRecallDecisionSubmission(paymentID, recallID, decisionID, sub)
	})

	w.Header().Set("Content-Type", jsonapi.ContentType)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(jsonapi.DataEnvelope[models.RecallDecisionSubmission]{Data: s})
}

func (h *RecallDecisionSubmissionHandler) Get(w http.ResponseWriter, r *http.Request) {
	paymentID := r.PathValue("paymentID")
	recallID := r.PathValue("recallID")
	decisionID := r.PathValue("decisionID")
	submissionID := r.PathValue("submissionID")

	s, err := h.store.GetRecallDecisionSubmission(paymentID, recallID, decisionID, submissionID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			jsonapi.NotFound(w, "recall_decision_submission", submissionID)
			return
		}
		jsonapi.InternalError(w)
		return
	}

	w.Header().Set("Content-Type", jsonapi.ContentType)
	json.NewEncoder(w).Encode(jsonapi.DataEnvelope[models.RecallDecisionSubmission]{Data: s})
}
