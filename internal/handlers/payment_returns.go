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

type PaymentReturnHandler struct {
	store store.Store
}

func NewPaymentReturnHandler(s store.Store) *PaymentReturnHandler {
	return &PaymentReturnHandler{store: s}
}

func (h *PaymentReturnHandler) Create(w http.ResponseWriter, r *http.Request) {
	paymentID := r.PathValue("paymentID")

	if _, err := h.store.GetPayment(paymentID); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			jsonapi.NotFound(w, "payment", paymentID)
			return
		}
		jsonapi.InternalError(w)
		return
	}

	var req jsonapi.DataEnvelope[models.ReturnPayment]
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonapi.BadRequest(w, "invalid JSON: "+err.Error())
		return
	}

	ret := req.Data
	if ret.ID == "" {
		ret.ID = uuid.New().String()
	}
	ret.Type = models.ResourceTypeReturnPayment
	now := time.Now().UTC()
	ret.CreatedOn = now
	ret.ModifiedOn = now

	if err := h.store.CreateReturn(paymentID, ret); err != nil {
		if errors.Is(err, store.ErrConflict) {
			jsonapi.Conflict(w, "return already exists")
			return
		}
		jsonapi.InternalError(w)
		return
	}

	w.Header().Set("Content-Type", jsonapi.ContentType)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(jsonapi.DataEnvelope[models.ReturnPayment]{Data: ret})
}

func (h *PaymentReturnHandler) Get(w http.ResponseWriter, r *http.Request) {
	paymentID := r.PathValue("paymentID")
	returnID := r.PathValue("returnID")

	ret, err := h.store.GetReturn(paymentID, returnID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			jsonapi.NotFound(w, "return_payment", returnID)
			return
		}
		jsonapi.InternalError(w)
		return
	}

	w.Header().Set("Content-Type", jsonapi.ContentType)
	json.NewEncoder(w).Encode(jsonapi.DataEnvelope[models.ReturnPayment]{Data: ret})
}

func (h *PaymentReturnHandler) List(w http.ResponseWriter, r *http.Request) {
	paymentID := r.PathValue("paymentID")
	returns := h.store.ListReturns(paymentID)
	w.Header().Set("Content-Type", jsonapi.ContentType)
	json.NewEncoder(w).Encode(jsonapi.ListEnvelope[models.ReturnPayment]{Data: returns})
}
