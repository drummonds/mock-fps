package store

import (
	"testing"
	"time"

	"github.com/nibble/mock-fps/internal/models"
)

func newPayment(id string) models.Payment {
	return models.Payment{
		Resource: models.Resource{
			Type:      models.ResourceTypePayment,
			ID:        id,
			CreatedOn: time.Now().UTC(),
		},
		Attributes: models.PaymentAttributes{
			Amount:   "100.00",
			Currency: "GBP",
		},
	}
}

func TestCreateAndGetPayment(t *testing.T) {
	s := NewMemoryStore()
	p := newPayment("p1")

	if err := s.CreatePayment(p); err != nil {
		t.Fatalf("CreatePayment: %v", err)
	}

	got, err := s.GetPayment("p1")
	if err != nil {
		t.Fatalf("GetPayment: %v", err)
	}
	if got.ID != "p1" {
		t.Errorf("expected id p1, got %s", got.ID)
	}
	if got.Attributes.Amount != "100.00" {
		t.Errorf("expected amount 100.00, got %s", got.Attributes.Amount)
	}
}

func TestCreatePaymentConflict(t *testing.T) {
	s := NewMemoryStore()
	p := newPayment("p1")
	s.CreatePayment(p)

	err := s.CreatePayment(p)
	if err != ErrConflict {
		t.Errorf("expected ErrConflict, got %v", err)
	}
}

func TestGetPaymentNotFound(t *testing.T) {
	s := NewMemoryStore()
	_, err := s.GetPayment("nonexistent")
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestListPayments(t *testing.T) {
	s := NewMemoryStore()
	s.CreatePayment(newPayment("p1"))
	s.CreatePayment(newPayment("p2"))

	payments := s.ListPayments()
	if len(payments) != 2 {
		t.Errorf("expected 2 payments, got %d", len(payments))
	}
}

func TestPaymentSubmissionCRUD(t *testing.T) {
	s := NewMemoryStore()
	s.CreatePayment(newPayment("p1"))

	sub := models.PaymentSubmission{
		Resource: models.Resource{ID: "s1", Type: models.ResourceTypePaymentSubmission},
		Attributes: models.PaymentSubmissionAttributes{Status: "accepted"},
	}

	if err := s.CreatePaymentSubmission("p1", sub); err != nil {
		t.Fatalf("CreatePaymentSubmission: %v", err)
	}

	got, err := s.GetPaymentSubmission("p1", "s1")
	if err != nil {
		t.Fatalf("GetPaymentSubmission: %v", err)
	}
	if got.Attributes.Status != "accepted" {
		t.Errorf("expected status accepted, got %s", got.Attributes.Status)
	}

	// Update
	got.Attributes.Status = "delivery_confirmed"
	if err := s.UpdatePaymentSubmission("p1", got); err != nil {
		t.Fatalf("UpdatePaymentSubmission: %v", err)
	}

	updated, _ := s.GetPaymentSubmission("p1", "s1")
	if updated.Attributes.Status != "delivery_confirmed" {
		t.Errorf("expected delivery_confirmed, got %s", updated.Attributes.Status)
	}

	// List
	subs := s.ListPaymentSubmissions("p1")
	if len(subs) != 1 {
		t.Errorf("expected 1 submission, got %d", len(subs))
	}
}

func TestSubscriptionMatchAndDelete(t *testing.T) {
	s := NewMemoryStore()
	sub := models.Subscription{
		Resource: models.Resource{ID: "sub1", Type: models.ResourceTypeSubscription},
		Attributes: models.SubscriptionAttributes{
			CallbackURI: "http://example.com/hook",
			EventType:   "updated",
			RecordType:  "payment_submissions",
			IsActive:    true,
		},
	}

	s.CreateSubscription(sub)

	matches := s.MatchSubscriptions("payment_submissions", "updated")
	if len(matches) != 1 {
		t.Errorf("expected 1 match, got %d", len(matches))
	}

	// No match for different event
	matches = s.MatchSubscriptions("payment_submissions", "created")
	if len(matches) != 0 {
		t.Errorf("expected 0 matches, got %d", len(matches))
	}

	// Delete
	if err := s.DeleteSubscription("sub1"); err != nil {
		t.Fatalf("DeleteSubscription: %v", err)
	}

	_, err := s.GetSubscription("sub1")
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestNestedResourceKeys(t *testing.T) {
	s := NewMemoryStore()
	s.CreatePayment(newPayment("p1"))

	ret := models.ReturnPayment{
		Resource: models.Resource{ID: "r1", Type: models.ResourceTypeReturnPayment},
	}
	s.CreateReturn("p1", ret)

	retSub := models.ReturnSubmission{
		Resource: models.Resource{ID: "rs1", Type: models.ResourceTypeReturnSubmission},
		Attributes: models.ReturnSubmissionAttributes{Status: "accepted"},
	}
	s.CreateReturnSubmission("p1", "r1", retSub)

	got, err := s.GetReturnSubmission("p1", "r1", "rs1")
	if err != nil {
		t.Fatalf("GetReturnSubmission: %v", err)
	}
	if got.Attributes.Status != "accepted" {
		t.Errorf("expected accepted, got %s", got.Attributes.Status)
	}

	// Recall -> Decision -> Decision Submission (4-level nesting)
	rec := models.Recall{Resource: models.Resource{ID: "rec1", Type: models.ResourceTypeRecall}}
	s.CreateRecall("p1", rec)

	dec := models.RecallDecision{Resource: models.Resource{ID: "d1", Type: models.ResourceTypeRecallDecision}}
	s.CreateRecallDecision("p1", "rec1", dec)

	ds := models.RecallDecisionSubmission{
		Resource: models.Resource{ID: "ds1", Type: models.ResourceTypeRecallDecisionSubmission},
		Attributes: models.RecallDecisionSubmissionAttributes{Status: "accepted"},
	}
	s.CreateRecallDecisionSubmission("p1", "rec1", "d1", ds)

	gotDS, err := s.GetRecallDecisionSubmission("p1", "rec1", "d1", "ds1")
	if err != nil {
		t.Fatalf("GetRecallDecisionSubmission: %v", err)
	}
	if gotDS.Attributes.Status != "accepted" {
		t.Errorf("expected accepted, got %s", gotDS.Attributes.Status)
	}
}
