package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/nibble/mock-fps/internal/fpsid"
	"github.com/nibble/mock-fps/internal/jsonapi"
	"github.com/nibble/mock-fps/internal/models"
	"github.com/nibble/mock-fps/internal/store"
)

type PaymentHandler struct {
	store store.Store
}

func NewPaymentHandler(s store.Store) *PaymentHandler {
	return &PaymentHandler{store: s}
}

func (h *PaymentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req jsonapi.DataEnvelope[models.Payment]
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonapi.BadRequest(w, "invalid JSON: "+err.Error())
		return
	}

	p := req.Data
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	p.Type = models.ResourceTypePayment
	now := time.Now().UTC()
	p.CreatedOn = now
	p.ModifiedOn = now

	// Default FPS attributes
	if p.Attributes.PaymentScheme == "" {
		p.Attributes.PaymentScheme = "FPS"
	}
	if p.Attributes.SchemePaymentType == "" {
		p.Attributes.SchemePaymentType = "ImmediatePayment"
	}
	if p.Attributes.NumericReference == "" {
		p.Attributes.NumericReference = fpsid.NumericReference()
	}
	if p.Attributes.EndToEndReference == "" {
		p.Attributes.EndToEndReference = fpsid.EndToEndReference()
	}
	if p.Attributes.ProcessingDate == "" {
		p.Attributes.ProcessingDate = fpsid.ProcessingDate()
	}

	if err := h.store.CreatePayment(p); err != nil {
		if errors.Is(err, store.ErrConflict) {
			jsonapi.Conflict(w, "payment already exists")
			return
		}
		jsonapi.InternalError(w)
		return
	}

	w.Header().Set("Content-Type", jsonapi.ContentType)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(jsonapi.DataEnvelope[models.Payment]{Data: p})
}

func (h *PaymentHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("paymentID")
	p, err := h.store.GetPayment(id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			jsonapi.NotFound(w, "payment", id)
			return
		}
		jsonapi.InternalError(w)
		return
	}

	// Build relationships
	p.Relationships = h.buildRelationships(id)

	w.Header().Set("Content-Type", jsonapi.ContentType)
	json.NewEncoder(w).Encode(jsonapi.DataEnvelope[models.Payment]{Data: p})
}

func (h *PaymentHandler) List(w http.ResponseWriter, r *http.Request) {
	payments := h.store.ListPayments()

	// Parse filter params
	q := r.URL.Query()
	subCycle := q.Get("filter[submission.settlement_cycle]")
	subDate := q.Get("filter[submission.settlement_date]")
	admCycle := q.Get("filter[admission.settlement_cycle]")
	admDate := q.Get("filter[admission.settlement_date]")

	hasFilter := subCycle != "" || subDate != "" || admCycle != "" || admDate != ""

	if hasFilter {
		var filtered []models.Payment
		for _, p := range payments {
			if h.matchesFilters(p.ID, subCycle, subDate, admCycle, admDate) {
				filtered = append(filtered, p)
			}
		}
		payments = filtered
	}

	// Pagination
	pageSize := 100
	pageNum := 1
	if v := q.Get("page[size]"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			pageSize = n
		}
	}
	if v := q.Get("page[number]"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			pageNum = n
		}
	}

	start := (pageNum - 1) * pageSize
	if start >= len(payments) {
		payments = nil
	} else {
		end := start + pageSize
		if end > len(payments) {
			end = len(payments)
		}
		payments = payments[start:end]
	}

	if payments == nil {
		payments = []models.Payment{}
	}

	w.Header().Set("Content-Type", jsonapi.ContentType)
	json.NewEncoder(w).Encode(jsonapi.ListEnvelope[models.Payment]{Data: payments})
}

func (h *PaymentHandler) matchesFilters(paymentID, subCycle, subDate, admCycle, admDate string) bool {
	if subCycle != "" || subDate != "" {
		subs := h.store.ListPaymentSubmissions(paymentID)
		if !matchSubmissions(subs, subCycle, subDate) {
			return false
		}
	}
	if admCycle != "" || admDate != "" {
		adms := h.store.ListPaymentAdmissions(paymentID)
		if !matchAdmissions(adms, admCycle, admDate) {
			return false
		}
	}
	return true
}

func matchSubmissions(subs []models.PaymentSubmission, cycle, date string) bool {
	for _, s := range subs {
		if cycle != "" {
			if n, err := strconv.Atoi(cycle); err != nil || s.Attributes.SettlementCycle != n {
				continue
			}
		}
		if date != "" && s.Attributes.SettlementDate != date {
			continue
		}
		return true
	}
	return false
}

func matchAdmissions(adms []models.PaymentAdmission, cycle, date string) bool {
	for _, a := range adms {
		if cycle != "" {
			if n, err := strconv.Atoi(cycle); err != nil || a.Attributes.SettlementCycle != n {
				continue
			}
		}
		if date != "" && a.Attributes.SettlementDate != date {
			continue
		}
		return true
	}
	return false
}

func (h *PaymentHandler) buildRelationships(paymentID string) *models.PaymentRelationships {
	rel := &models.PaymentRelationships{}
	hasAny := false

	if subs := h.store.ListPaymentSubmissions(paymentID); len(subs) > 0 {
		data := make([]models.RelationshipData, len(subs))
		for i, s := range subs {
			data[i] = models.RelationshipData{Type: models.ResourceTypePaymentSubmission, ID: s.ID}
		}
		rel.PaymentSubmissions = &models.Relationship{Data: data}
		hasAny = true
	}

	if adms := h.store.ListPaymentAdmissions(paymentID); len(adms) > 0 {
		data := make([]models.RelationshipData, len(adms))
		for i, a := range adms {
			data[i] = models.RelationshipData{Type: models.ResourceTypePaymentAdmission, ID: a.ID}
		}
		rel.PaymentAdmissions = &models.Relationship{Data: data}
		hasAny = true
	}

	if rets := h.store.ListReturns(paymentID); len(rets) > 0 {
		data := make([]models.RelationshipData, len(rets))
		for i, rt := range rets {
			data[i] = models.RelationshipData{Type: models.ResourceTypeReturnPayment, ID: rt.ID}
		}
		rel.PaymentReturns = &models.Relationship{Data: data}
		hasAny = true
	}

	if recs := h.store.ListRecalls(paymentID); len(recs) > 0 {
		data := make([]models.RelationshipData, len(recs))
		for i, rc := range recs {
			data[i] = models.RelationshipData{Type: models.ResourceTypeRecall, ID: rc.ID}
		}
		rel.PaymentRecalls = &models.Relationship{Data: data}
		hasAny = true
	}

	if revs := h.store.ListReversals(paymentID); len(revs) > 0 {
		data := make([]models.RelationshipData, len(revs))
		for i, rv := range revs {
			data[i] = models.RelationshipData{Type: models.ResourceTypeReversal, ID: rv.ID}
		}
		rel.PaymentReversals = &models.Relationship{Data: data}
		hasAny = true
	}

	if !hasAny {
		return nil
	}
	return rel
}
