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

type PaymentReversalHandler struct {
	store store.Store
}

func NewPaymentReversalHandler(s store.Store) *PaymentReversalHandler {
	return &PaymentReversalHandler{store: s}
}

func (h *PaymentReversalHandler) Create(w http.ResponseWriter, r *http.Request) {
	paymentID := r.PathValue("paymentID")

	if _, err := h.store.GetPayment(paymentID); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			jsonapi.NotFound(w, "payment", paymentID)
			return
		}
		jsonapi.InternalError(w)
		return
	}

	var req jsonapi.DataEnvelope[models.Reversal]
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonapi.BadRequest(w, "invalid JSON: "+err.Error())
		return
	}

	rev := req.Data
	if rev.ID == "" {
		rev.ID = uuid.New().String()
	}
	rev.Type = models.ResourceTypeReversal
	now := time.Now().UTC()
	rev.CreatedOn = now
	rev.ModifiedOn = now

	if err := h.store.CreateReversal(paymentID, rev); err != nil {
		if errors.Is(err, store.ErrConflict) {
			jsonapi.Conflict(w, "reversal already exists")
			return
		}
		jsonapi.InternalError(w)
		return
	}

	w.Header().Set("Content-Type", jsonapi.ContentType)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(jsonapi.DataEnvelope[models.Reversal]{Data: rev})
}

func (h *PaymentReversalHandler) Get(w http.ResponseWriter, r *http.Request) {
	paymentID := r.PathValue("paymentID")
	reversalID := r.PathValue("reversalID")

	rev, err := h.store.GetReversal(paymentID, reversalID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			jsonapi.NotFound(w, "reversal", reversalID)
			return
		}
		jsonapi.InternalError(w)
		return
	}

	w.Header().Set("Content-Type", jsonapi.ContentType)
	json.NewEncoder(w).Encode(jsonapi.DataEnvelope[models.Reversal]{Data: rev})
}

func (h *PaymentReversalHandler) List(w http.ResponseWriter, r *http.Request) {
	paymentID := r.PathValue("paymentID")
	reversals := h.store.ListReversals(paymentID)
	w.Header().Set("Content-Type", jsonapi.ContentType)
	json.NewEncoder(w).Encode(jsonapi.ListEnvelope[models.Reversal]{Data: reversals})
}
