package store

import (
	"testing"
	"time"

	"github.com/nibble/mock-fps/internal/models"
)

func TestPurgeOlderThan(t *testing.T) {
	s := NewMemoryStore()

	now := time.Now().UTC()
	old := now.Add(-4 * 24 * time.Hour)
	recent := now.Add(-1 * 24 * time.Hour)

	// Create old payment with children
	s.CreatePayment(models.Payment{
		Resource: models.Resource{ID: "old-pay", CreatedOn: old},
	})
	s.CreatePaymentSubmission("old-pay", models.PaymentSubmission{
		Resource: models.Resource{ID: "old-sub"},
	})
	s.CreatePaymentAdmission("old-pay", models.PaymentAdmission{
		Resource: models.Resource{ID: "old-adm"},
	})

	// Create recent payment with children
	s.CreatePayment(models.Payment{
		Resource: models.Resource{ID: "new-pay", CreatedOn: recent},
	})
	s.CreatePaymentSubmission("new-pay", models.PaymentSubmission{
		Resource: models.Resource{ID: "new-sub"},
	})

	cutoff := now.Add(-3 * 24 * time.Hour)
	purged := s.PurgeOlderThan(cutoff)

	if purged != 1 {
		t.Errorf("purged = %d, want 1", purged)
	}

	// Old payment and children gone
	if _, err := s.GetPayment("old-pay"); err != ErrNotFound {
		t.Error("expected old payment to be purged")
	}
	if _, err := s.GetPaymentSubmission("old-pay", "old-sub"); err != ErrNotFound {
		t.Error("expected old submission to be purged")
	}
	if _, err := s.GetPaymentAdmission("old-pay", "old-adm"); err != ErrNotFound {
		t.Error("expected old admission to be purged")
	}

	// Recent payment still present
	if _, err := s.GetPayment("new-pay"); err != nil {
		t.Error("expected recent payment to remain")
	}
	if _, err := s.GetPaymentSubmission("new-pay", "new-sub"); err != nil {
		t.Error("expected recent submission to remain")
	}
}

func TestPurgeOlderThan_NothingToPurge(t *testing.T) {
	s := NewMemoryStore()

	s.CreatePayment(models.Payment{
		Resource: models.Resource{ID: "p1", CreatedOn: time.Now().UTC()},
	})

	purged := s.PurgeOlderThan(time.Now().UTC().Add(-24 * time.Hour))
	if purged != 0 {
		t.Errorf("purged = %d, want 0", purged)
	}
}
