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

type PaymentAdmissionHandler struct {
	store  store.Store
	engine *lifecycle.Engine
}

func NewPaymentAdmissionHandler(s store.Store, e *lifecycle.Engine) *PaymentAdmissionHandler {
	return &PaymentAdmissionHandler{store: s, engine: e}
}

func (h *PaymentAdmissionHandler) Create(w http.ResponseWriter, r *http.Request) {
	paymentID := r.PathValue("paymentID")

	if _, err := h.store.GetPayment(paymentID); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			jsonapi.NotFound(w, "payment", paymentID)
			return
		}
		jsonapi.InternalError(w)
		return
	}

	var req jsonapi.DataEnvelope[models.PaymentAdmission]
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonapi.BadRequest(w, "invalid JSON: "+err.Error())
		return
	}

	a := req.Data
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	a.Type = models.ResourceTypePaymentAdmission
	now := time.Now().UTC()
	a.CreatedOn = now
	a.ModifiedOn = now
	a.Attributes.Status = lifecycle.AdmissionChain[0]
	a.Attributes.AdmissionDate = now.Format(time.DateOnly)

	if err := h.store.CreatePaymentAdmission(paymentID, a); err != nil {
		if errors.Is(err, store.ErrConflict) {
			jsonapi.Conflict(w, "admission already exists")
			return
		}
		jsonapi.InternalError(w)
		return
	}

	// Start async lifecycle
	admissionID := a.ID
	h.engine.StartTransition(models.ResourceTypePaymentAdmission, admissionID, lifecycle.AdmissionChain, func(newStatus string) error {
		adm, err := h.store.GetPaymentAdmission(paymentID, admissionID)
		if err != nil {
			return err
		}
		adm.Attributes.Status = newStatus
		adm.ModifiedOn = time.Now().UTC()
		return h.store.UpdatePaymentAdmission(paymentID, adm)
	})

	w.Header().Set("Content-Type", jsonapi.ContentType)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(jsonapi.DataEnvelope[models.PaymentAdmission]{Data: a})
}

func (h *PaymentAdmissionHandler) Get(w http.ResponseWriter, r *http.Request) {
	paymentID := r.PathValue("paymentID")
	admissionID := r.PathValue("admissionID")

	a, err := h.store.GetPaymentAdmission(paymentID, admissionID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			jsonapi.NotFound(w, "payment_admission", admissionID)
			return
		}
		jsonapi.InternalError(w)
		return
	}

	w.Header().Set("Content-Type", jsonapi.ContentType)
	json.NewEncoder(w).Encode(jsonapi.DataEnvelope[models.PaymentAdmission]{Data: a})
}

func (h *PaymentAdmissionHandler) PatchTask(w http.ResponseWriter, r *http.Request) {
	paymentID := r.PathValue("paymentID")
	admissionID := r.PathValue("admissionID")
	taskID := r.PathValue("taskID")

	// Get existing task or create if first access
	t, err := h.store.GetAdmissionTask(paymentID, admissionID, taskID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			// Auto-create the task
			t = models.AdmissionTask{
				Resource: models.Resource{
					Type:      "admission_tasks",
					ID:        taskID,
					CreatedOn: time.Now().UTC(),
					ModifiedOn: time.Now().UTC(),
				},
			}
			if createErr := h.store.CreateAdmissionTask(paymentID, admissionID, t); createErr != nil {
				if !errors.Is(createErr, store.ErrConflict) {
					jsonapi.InternalError(w)
					return
				}
				// Another goroutine created it, re-fetch
				t, err = h.store.GetAdmissionTask(paymentID, admissionID, taskID)
				if err != nil {
					jsonapi.InternalError(w)
					return
				}
			}
		} else {
			jsonapi.InternalError(w)
			return
		}
	}

	var req jsonapi.DataEnvelope[models.AdmissionTask]
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonapi.BadRequest(w, "invalid JSON: "+err.Error())
		return
	}

	patch := req.Data
	if patch.Attributes.Status != "" {
		t.Attributes.Status = patch.Attributes.Status
	}
	if patch.Attributes.Assignee != "" {
		t.Attributes.Assignee = patch.Attributes.Assignee
	}
	if patch.Attributes.Name != "" {
		t.Attributes.Name = patch.Attributes.Name
	}
	t.ModifiedOn = time.Now().UTC()

	if err := h.store.UpdateAdmissionTask(paymentID, admissionID, t); err != nil {
		jsonapi.InternalError(w)
		return
	}

	w.Header().Set("Content-Type", jsonapi.ContentType)
	json.NewEncoder(w).Encode(jsonapi.DataEnvelope[models.AdmissionTask]{Data: t})
}
