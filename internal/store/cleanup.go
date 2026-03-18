package store

import (
	"context"
	"log"
	"strings"
	"time"
)

// PurgeOlderThan removes payments (and all child resources) with CreatedOn before cutoff.
// Returns the number of purged payments.
func (m *MemoryStore) PurgeOlderThan(cutoff time.Time) int {
	m.mu.Lock()
	defer m.mu.Unlock()

	var purged int
	for id, p := range m.payments {
		if p.CreatedOn.Before(cutoff) {
			delete(m.payments, id)
			m.deleteChildrenLocked(id)
			purged++
		}
	}
	return purged
}

// deleteChildrenLocked removes all child resources keyed by paymentID prefix.
// Must be called with mu held.
func (m *MemoryStore) deleteChildrenLocked(paymentID string) {
	prefix := paymentID + ":"
	deleteByPrefix(m.paymentSubmissions, prefix)
	deleteByPrefix(m.paymentAdmissions, prefix)
	deleteByPrefix(m.admissionTasks, prefix)
	deleteByPrefix(m.returns, prefix)
	deleteByPrefix(m.returnSubmissions, prefix)
	deleteByPrefix(m.recalls, prefix)
	deleteByPrefix(m.recallSubmissions, prefix)
	deleteByPrefix(m.recallDecisions, prefix)
	deleteByPrefix(m.recallDecisionSubmissions, prefix)
	deleteByPrefix(m.reversals, prefix)
	deleteByPrefix(m.reversalSubmissions, prefix)
}

func deleteByPrefix[V any](m map[string]V, prefix string) {
	for k := range m {
		if strings.HasPrefix(k, prefix) {
			delete(m, k)
		}
	}
}

// RunCleanup runs periodic cleanup of old payments. Blocks until ctx is cancelled.
func RunCleanup(ctx context.Context, s *MemoryStore, retentionDays int) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			cutoff := time.Now().UTC().AddDate(0, 0, -retentionDays)
			if n := s.PurgeOlderThan(cutoff); n > 0 {
				log.Printf("cleanup: purged %d payments older than %s", n, cutoff.Format(time.DateOnly))
			}
		}
	}
}
