package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/nibble/mock-fps/internal/jsonapi"
	"github.com/nibble/mock-fps/internal/models"
	"github.com/nibble/mock-fps/internal/store"
)

type RecallDecisionHandler struct {
	store store.Store
}

func NewRecallDecisionHandler(s store.Store) *RecallDecisionHandler {
	return &RecallDecisionHandler{store: s}
}

func (h *RecallDecisionHandler) Create(w http.ResponseWriter, r *http.Request) {
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

	var req jsonapi.DataEnvelope[models.RecallDecision]
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonapi.BadRequest(w, "invalid JSON: "+err.Error())
		return
	}

	d := req.Data
	if d.ID == "" {
		d.ID = uuid.New().String()
	}
	d.Type = models.ResourceTypeRecallDecision
	now := time.Now().UTC()
	d.CreatedOn = now
	d.ModifiedOn = now

	if err := h.store.CreateRecallDecision(paymentID, recallID, d); err != nil {
		if errors.Is(err, store.ErrConflict) {
			jsonapi.Conflict(w, "recall decision already exists")
			return
		}
		jsonapi.InternalError(w)
		return
	}

	w.Header().Set("Content-Type", jsonapi.ContentType)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(jsonapi.DataEnvelope[models.RecallDecision]{Data: d})
}

func (h *RecallDecisionHandler) Get(w http.ResponseWriter, r *http.Request) {
	paymentID := r.PathValue("paymentID")
	recallID := r.PathValue("recallID")
	decisionID := r.PathValue("decisionID")

	d, err := h.store.GetRecallDecision(paymentID, recallID, decisionID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			jsonapi.NotFound(w, "recall_decision", decisionID)
			return
		}
		jsonapi.InternalError(w)
		return
	}

	w.Header().Set("Content-Type", jsonapi.ContentType)
	json.NewEncoder(w).Encode(jsonapi.DataEnvelope[models.RecallDecision]{Data: d})
}

func (h *RecallDecisionHandler) List(w http.ResponseWriter, r *http.Request) {
	paymentID := r.PathValue("paymentID")
	recallID := r.PathValue("recallID")
	decisions := h.store.ListRecallDecisions(paymentID, recallID)
	w.Header().Set("Content-Type", jsonapi.ContentType)
	json.NewEncoder(w).Encode(jsonapi.ListEnvelope[models.RecallDecision]{Data: decisions})
}
