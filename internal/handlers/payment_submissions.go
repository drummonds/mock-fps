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

type PaymentSubmissionHandler struct {
	store  store.Store
	engine *lifecycle.Engine
}

func NewPaymentSubmissionHandler(s store.Store, e *lifecycle.Engine) *PaymentSubmissionHandler {
	return &PaymentSubmissionHandler{store: s, engine: e}
}

func (h *PaymentSubmissionHandler) Create(w http.ResponseWriter, r *http.Request) {
	paymentID := r.PathValue("paymentID")

	// Verify payment exists
	if _, err := h.store.GetPayment(paymentID); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			jsonapi.NotFound(w, "payment", paymentID)
			return
		}
		jsonapi.InternalError(w)
		return
	}

	var req jsonapi.DataEnvelope[models.PaymentSubmission]
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonapi.BadRequest(w, "invalid JSON: "+err.Error())
		return
	}

	s := req.Data
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	s.Type = models.ResourceTypePaymentSubmission
	now := time.Now().UTC()
	s.CreatedOn = now
	s.ModifiedOn = now
	s.Attributes.Status = lifecycle.PaymentSubmissionChain[0]
	s.Attributes.SubmissionDate = now.Format(time.DateOnly)

	if err := h.store.CreatePaymentSubmission(paymentID, s); err != nil {
		if errors.Is(err, store.ErrConflict) {
			jsonapi.Conflict(w, "submission already exists")
			return
		}
		jsonapi.InternalError(w)
		return
	}

	// Start async lifecycle
	submissionID := s.ID
	h.engine.StartTransition(models.ResourceTypePaymentSubmission, submissionID, lifecycle.PaymentSubmissionChain, func(newStatus string) error {
		sub, err := h.store.GetPaymentSubmission(paymentID, submissionID)
		if err != nil {
			return err
		}
		sub.Attributes.Status = newStatus
		sub.ModifiedOn = time.Now().UTC()
		return h.store.UpdatePaymentSubmission(paymentID, sub)
	})

	w.Header().Set("Content-Type", jsonapi.ContentType)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(jsonapi.DataEnvelope[models.PaymentSubmission]{Data: s})
}

func (h *PaymentSubmissionHandler) Get(w http.ResponseWriter, r *http.Request) {
	paymentID := r.PathValue("paymentID")
	submissionID := r.PathValue("submissionID")

	s, err := h.store.GetPaymentSubmission(paymentID, submissionID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			jsonapi.NotFound(w, "payment_submission", submissionID)
			return
		}
		jsonapi.InternalError(w)
		return
	}

	w.Header().Set("Content-Type", jsonapi.ContentType)
	json.NewEncoder(w).Encode(jsonapi.DataEnvelope[models.PaymentSubmission]{Data: s})
}
