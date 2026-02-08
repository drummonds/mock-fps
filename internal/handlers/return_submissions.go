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

type ReturnSubmissionHandler struct {
	store  store.Store
	engine *lifecycle.Engine
}

func NewReturnSubmissionHandler(s store.Store, e *lifecycle.Engine) *ReturnSubmissionHandler {
	return &ReturnSubmissionHandler{store: s, engine: e}
}

func (h *ReturnSubmissionHandler) Create(w http.ResponseWriter, r *http.Request) {
	paymentID := r.PathValue("paymentID")
	returnID := r.PathValue("returnID")

	if _, err := h.store.GetReturn(paymentID, returnID); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			jsonapi.NotFound(w, "return_payment", returnID)
			return
		}
		jsonapi.InternalError(w)
		return
	}

	var req jsonapi.DataEnvelope[models.ReturnSubmission]
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonapi.BadRequest(w, "invalid JSON: "+err.Error())
		return
	}

	s := req.Data
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	s.Type = models.ResourceTypeReturnSubmission
	now := time.Now().UTC()
	s.CreatedOn = now
	s.ModifiedOn = now
	s.Attributes.Status = lifecycle.SimpleSubmissionChain[0]
	s.Attributes.SubmissionDate = now.Format(time.DateOnly)

	if err := h.store.CreateReturnSubmission(paymentID, returnID, s); err != nil {
		if errors.Is(err, store.ErrConflict) {
			jsonapi.Conflict(w, "return submission already exists")
			return
		}
		jsonapi.InternalError(w)
		return
	}

	submissionID := s.ID
	h.engine.StartTransition(models.ResourceTypeReturnSubmission, submissionID, lifecycle.SimpleSubmissionChain, func(newStatus string) error {
		sub, err := h.store.GetReturnSubmission(paymentID, returnID, submissionID)
		if err != nil {
			return err
		}
		sub.Attributes.Status = newStatus
		sub.ModifiedOn = time.Now().UTC()
		return h.store.UpdateReturnSubmission(paymentID, returnID, sub)
	})

	w.Header().Set("Content-Type", jsonapi.ContentType)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(jsonapi.DataEnvelope[models.ReturnSubmission]{Data: s})
}

func (h *ReturnSubmissionHandler) Get(w http.ResponseWriter, r *http.Request) {
	paymentID := r.PathValue("paymentID")
	returnID := r.PathValue("returnID")
	submissionID := r.PathValue("submissionID")

	s, err := h.store.GetReturnSubmission(paymentID, returnID, submissionID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			jsonapi.NotFound(w, "return_submission", submissionID)
			return
		}
		jsonapi.InternalError(w)
		return
	}

	w.Header().Set("Content-Type", jsonapi.ContentType)
	json.NewEncoder(w).Encode(jsonapi.DataEnvelope[models.ReturnSubmission]{Data: s})
}
