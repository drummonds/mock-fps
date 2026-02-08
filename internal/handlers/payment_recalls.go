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

type PaymentRecallHandler struct {
	store store.Store
}

func NewPaymentRecallHandler(s store.Store) *PaymentRecallHandler {
	return &PaymentRecallHandler{store: s}
}

func (h *PaymentRecallHandler) Create(w http.ResponseWriter, r *http.Request) {
	paymentID := r.PathValue("paymentID")

	if _, err := h.store.GetPayment(paymentID); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			jsonapi.NotFound(w, "payment", paymentID)
			return
		}
		jsonapi.InternalError(w)
		return
	}

	var req jsonapi.DataEnvelope[models.Recall]
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonapi.BadRequest(w, "invalid JSON: "+err.Error())
		return
	}

	rec := req.Data
	if rec.ID == "" {
		rec.ID = uuid.New().String()
	}
	rec.Type = models.ResourceTypeRecall
	now := time.Now().UTC()
	rec.CreatedOn = now
	rec.ModifiedOn = now

	if err := h.store.CreateRecall(paymentID, rec); err != nil {
		if errors.Is(err, store.ErrConflict) {
			jsonapi.Conflict(w, "recall already exists")
			return
		}
		jsonapi.InternalError(w)
		return
	}

	w.Header().Set("Content-Type", jsonapi.ContentType)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(jsonapi.DataEnvelope[models.Recall]{Data: rec})
}

func (h *PaymentRecallHandler) Get(w http.ResponseWriter, r *http.Request) {
	paymentID := r.PathValue("paymentID")
	recallID := r.PathValue("recallID")

	rec, err := h.store.GetRecall(paymentID, recallID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			jsonapi.NotFound(w, "recall", recallID)
			return
		}
		jsonapi.InternalError(w)
		return
	}

	w.Header().Set("Content-Type", jsonapi.ContentType)
	json.NewEncoder(w).Encode(jsonapi.DataEnvelope[models.Recall]{Data: rec})
}

func (h *PaymentRecallHandler) List(w http.ResponseWriter, r *http.Request) {
	paymentID := r.PathValue("paymentID")
	recalls := h.store.ListRecalls(paymentID)
	w.Header().Set("Content-Type", jsonapi.ContentType)
	json.NewEncoder(w).Encode(jsonapi.ListEnvelope[models.Recall]{Data: recalls})
}
