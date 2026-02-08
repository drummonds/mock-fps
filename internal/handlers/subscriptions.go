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

type SubscriptionHandler struct {
	store store.Store
}

func NewSubscriptionHandler(s store.Store) *SubscriptionHandler {
	return &SubscriptionHandler{store: s}
}

func (h *SubscriptionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req jsonapi.DataEnvelope[models.Subscription]
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonapi.BadRequest(w, "invalid JSON: "+err.Error())
		return
	}

	s := req.Data
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	s.Type = models.ResourceTypeSubscription
	now := time.Now().UTC()
	s.CreatedOn = now
	s.ModifiedOn = now
	if !s.Attributes.IsActive {
		s.Attributes.IsActive = true
	}

	if err := h.store.CreateSubscription(s); err != nil {
		if errors.Is(err, store.ErrConflict) {
			jsonapi.Conflict(w, "subscription already exists")
			return
		}
		jsonapi.InternalError(w)
		return
	}

	w.Header().Set("Content-Type", jsonapi.ContentType)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(jsonapi.DataEnvelope[models.Subscription]{Data: s})
}

func (h *SubscriptionHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("subscriptionID")

	s, err := h.store.GetSubscription(id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			jsonapi.NotFound(w, "subscription", id)
			return
		}
		jsonapi.InternalError(w)
		return
	}

	w.Header().Set("Content-Type", jsonapi.ContentType)
	json.NewEncoder(w).Encode(jsonapi.DataEnvelope[models.Subscription]{Data: s})
}

func (h *SubscriptionHandler) List(w http.ResponseWriter, r *http.Request) {
	subs := h.store.ListSubscriptions()
	w.Header().Set("Content-Type", jsonapi.ContentType)
	json.NewEncoder(w).Encode(jsonapi.ListEnvelope[models.Subscription]{Data: subs})
}

func (h *SubscriptionHandler) Patch(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("subscriptionID")

	existing, err := h.store.GetSubscription(id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			jsonapi.NotFound(w, "subscription", id)
			return
		}
		jsonapi.InternalError(w)
		return
	}

	var req jsonapi.DataEnvelope[models.Subscription]
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonapi.BadRequest(w, "invalid JSON: "+err.Error())
		return
	}

	patch := req.Data
	if patch.Attributes.CallbackURI != "" {
		existing.Attributes.CallbackURI = patch.Attributes.CallbackURI
	}
	if patch.Attributes.EventType != "" {
		existing.Attributes.EventType = patch.Attributes.EventType
	}
	if patch.Attributes.RecordType != "" {
		existing.Attributes.RecordType = patch.Attributes.RecordType
	}
	// Allow setting is_active to false explicitly via the raw JSON
	existing.Attributes.IsActive = patch.Attributes.IsActive
	existing.ModifiedOn = time.Now().UTC()
	existing.Version++

	if err := h.store.UpdateSubscription(existing); err != nil {
		jsonapi.InternalError(w)
		return
	}

	w.Header().Set("Content-Type", jsonapi.ContentType)
	json.NewEncoder(w).Encode(jsonapi.DataEnvelope[models.Subscription]{Data: existing})
}

func (h *SubscriptionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("subscriptionID")

	if err := h.store.DeleteSubscription(id); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			jsonapi.NotFound(w, "subscription", id)
			return
		}
		jsonapi.InternalError(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
