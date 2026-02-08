package store

import (
	"fmt"
	"sync"

	"github.com/nibble/mock-fps/internal/models"
)

// ErrNotFound is returned when a resource is not found.
var ErrNotFound = fmt.Errorf("not found")

// ErrConflict is returned when a resource already exists.
var ErrConflict = fmt.Errorf("conflict")

// MemoryStore is an in-memory implementation of Store.
type MemoryStore struct {
	mu sync.RWMutex

	payments                  map[string]models.Payment
	paymentSubmissions        map[string]models.PaymentSubmission        // "paymentID:submissionID"
	paymentAdmissions         map[string]models.PaymentAdmission         // "paymentID:admissionID"
	admissionTasks            map[string]models.AdmissionTask            // "paymentID:admissionID:taskID"
	returns                   map[string]models.ReturnPayment            // "paymentID:returnID"
	returnSubmissions         map[string]models.ReturnSubmission         // "paymentID:returnID:submissionID"
	recalls                   map[string]models.Recall                   // "paymentID:recallID"
	recallSubmissions         map[string]models.RecallSubmission         // "paymentID:recallID:submissionID"
	recallDecisions           map[string]models.RecallDecision           // "paymentID:recallID:decisionID"
	recallDecisionSubmissions map[string]models.RecallDecisionSubmission // "paymentID:recallID:decisionID:submissionID"
	reversals                 map[string]models.Reversal                 // "paymentID:reversalID"
	reversalSubmissions       map[string]models.ReversalSubmission       // "paymentID:reversalID:submissionID"
	subscriptions             map[string]models.Subscription
}

// NewMemoryStore creates a new in-memory store.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		payments:                  make(map[string]models.Payment),
		paymentSubmissions:        make(map[string]models.PaymentSubmission),
		paymentAdmissions:         make(map[string]models.PaymentAdmission),
		admissionTasks:            make(map[string]models.AdmissionTask),
		returns:                   make(map[string]models.ReturnPayment),
		returnSubmissions:         make(map[string]models.ReturnSubmission),
		recalls:                   make(map[string]models.Recall),
		recallSubmissions:         make(map[string]models.RecallSubmission),
		recallDecisions:           make(map[string]models.RecallDecision),
		recallDecisionSubmissions: make(map[string]models.RecallDecisionSubmission),
		reversals:                 make(map[string]models.Reversal),
		reversalSubmissions:       make(map[string]models.ReversalSubmission),
		subscriptions:             make(map[string]models.Subscription),
	}
}

func key2(a, b string) string         { return a + ":" + b }
func key3(a, b, c string) string      { return a + ":" + b + ":" + c }
func key4(a, b, c, d string) string   { return a + ":" + b + ":" + c + ":" + d }

// --- Payments ---

func (m *MemoryStore) CreatePayment(p models.Payment) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.payments[p.ID]; ok {
		return ErrConflict
	}
	m.payments[p.ID] = p
	return nil
}

func (m *MemoryStore) GetPayment(id string) (models.Payment, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	p, ok := m.payments[id]
	if !ok {
		return p, ErrNotFound
	}
	return p, nil
}

func (m *MemoryStore) ListPayments() []models.Payment {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]models.Payment, 0, len(m.payments))
	for _, p := range m.payments {
		out = append(out, p)
	}
	return out
}

func (m *MemoryStore) UpdatePayment(p models.Payment) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.payments[p.ID]; !ok {
		return ErrNotFound
	}
	m.payments[p.ID] = p
	return nil
}

// --- Payment Submissions ---

func (m *MemoryStore) CreatePaymentSubmission(paymentID string, s models.PaymentSubmission) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	k := key2(paymentID, s.ID)
	if _, ok := m.paymentSubmissions[k]; ok {
		return ErrConflict
	}
	m.paymentSubmissions[k] = s
	return nil
}

func (m *MemoryStore) GetPaymentSubmission(paymentID, submissionID string) (models.PaymentSubmission, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	s, ok := m.paymentSubmissions[key2(paymentID, submissionID)]
	if !ok {
		return s, ErrNotFound
	}
	return s, nil
}

func (m *MemoryStore) ListPaymentSubmissions(paymentID string) []models.PaymentSubmission {
	m.mu.RLock()
	defer m.mu.RUnlock()
	prefix := paymentID + ":"
	var out []models.PaymentSubmission
	for k, v := range m.paymentSubmissions {
		if len(k) > len(prefix) && k[:len(prefix)] == prefix {
			out = append(out, v)
		}
	}
	return out
}

func (m *MemoryStore) UpdatePaymentSubmission(paymentID string, s models.PaymentSubmission) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	k := key2(paymentID, s.ID)
	if _, ok := m.paymentSubmissions[k]; !ok {
		return ErrNotFound
	}
	m.paymentSubmissions[k] = s
	return nil
}

// --- Payment Admissions ---

func (m *MemoryStore) CreatePaymentAdmission(paymentID string, a models.PaymentAdmission) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	k := key2(paymentID, a.ID)
	if _, ok := m.paymentAdmissions[k]; ok {
		return ErrConflict
	}
	m.paymentAdmissions[k] = a
	return nil
}

func (m *MemoryStore) GetPaymentAdmission(paymentID, admissionID string) (models.PaymentAdmission, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	a, ok := m.paymentAdmissions[key2(paymentID, admissionID)]
	if !ok {
		return a, ErrNotFound
	}
	return a, nil
}

func (m *MemoryStore) ListPaymentAdmissions(paymentID string) []models.PaymentAdmission {
	m.mu.RLock()
	defer m.mu.RUnlock()
	prefix := paymentID + ":"
	var out []models.PaymentAdmission
	for k, v := range m.paymentAdmissions {
		if len(k) > len(prefix) && k[:len(prefix)] == prefix {
			out = append(out, v)
		}
	}
	return out
}

func (m *MemoryStore) UpdatePaymentAdmission(paymentID string, a models.PaymentAdmission) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	k := key2(paymentID, a.ID)
	if _, ok := m.paymentAdmissions[k]; !ok {
		return ErrNotFound
	}
	m.paymentAdmissions[k] = a
	return nil
}

// --- Admission Tasks ---

func (m *MemoryStore) CreateAdmissionTask(paymentID, admissionID string, t models.AdmissionTask) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	k := key3(paymentID, admissionID, t.ID)
	if _, ok := m.admissionTasks[k]; ok {
		return ErrConflict
	}
	m.admissionTasks[k] = t
	return nil
}

func (m *MemoryStore) GetAdmissionTask(paymentID, admissionID, taskID string) (models.AdmissionTask, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	t, ok := m.admissionTasks[key3(paymentID, admissionID, taskID)]
	if !ok {
		return t, ErrNotFound
	}
	return t, nil
}

func (m *MemoryStore) UpdateAdmissionTask(paymentID, admissionID string, t models.AdmissionTask) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	k := key3(paymentID, admissionID, t.ID)
	if _, ok := m.admissionTasks[k]; !ok {
		return ErrNotFound
	}
	m.admissionTasks[k] = t
	return nil
}

// --- Returns ---

func (m *MemoryStore) CreateReturn(paymentID string, r models.ReturnPayment) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	k := key2(paymentID, r.ID)
	if _, ok := m.returns[k]; ok {
		return ErrConflict
	}
	m.returns[k] = r
	return nil
}

func (m *MemoryStore) GetReturn(paymentID, returnID string) (models.ReturnPayment, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	r, ok := m.returns[key2(paymentID, returnID)]
	if !ok {
		return r, ErrNotFound
	}
	return r, nil
}

func (m *MemoryStore) ListReturns(paymentID string) []models.ReturnPayment {
	m.mu.RLock()
	defer m.mu.RUnlock()
	prefix := paymentID + ":"
	var out []models.ReturnPayment
	for k, v := range m.returns {
		if len(k) > len(prefix) && k[:len(prefix)] == prefix {
			out = append(out, v)
		}
	}
	return out
}

// --- Return Submissions ---

func (m *MemoryStore) CreateReturnSubmission(paymentID, returnID string, s models.ReturnSubmission) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	k := key3(paymentID, returnID, s.ID)
	if _, ok := m.returnSubmissions[k]; ok {
		return ErrConflict
	}
	m.returnSubmissions[k] = s
	return nil
}

func (m *MemoryStore) GetReturnSubmission(paymentID, returnID, submissionID string) (models.ReturnSubmission, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	s, ok := m.returnSubmissions[key3(paymentID, returnID, submissionID)]
	if !ok {
		return s, ErrNotFound
	}
	return s, nil
}

func (m *MemoryStore) ListReturnSubmissions(paymentID, returnID string) []models.ReturnSubmission {
	m.mu.RLock()
	defer m.mu.RUnlock()
	prefix := key2(paymentID, returnID) + ":"
	var out []models.ReturnSubmission
	for k, v := range m.returnSubmissions {
		if len(k) > len(prefix) && k[:len(prefix)] == prefix {
			out = append(out, v)
		}
	}
	return out
}

func (m *MemoryStore) UpdateReturnSubmission(paymentID, returnID string, s models.ReturnSubmission) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	k := key3(paymentID, returnID, s.ID)
	if _, ok := m.returnSubmissions[k]; !ok {
		return ErrNotFound
	}
	m.returnSubmissions[k] = s
	return nil
}

// --- Recalls ---

func (m *MemoryStore) CreateRecall(paymentID string, r models.Recall) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	k := key2(paymentID, r.ID)
	if _, ok := m.recalls[k]; ok {
		return ErrConflict
	}
	m.recalls[k] = r
	return nil
}

func (m *MemoryStore) GetRecall(paymentID, recallID string) (models.Recall, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	r, ok := m.recalls[key2(paymentID, recallID)]
	if !ok {
		return r, ErrNotFound
	}
	return r, nil
}

func (m *MemoryStore) ListRecalls(paymentID string) []models.Recall {
	m.mu.RLock()
	defer m.mu.RUnlock()
	prefix := paymentID + ":"
	var out []models.Recall
	for k, v := range m.recalls {
		if len(k) > len(prefix) && k[:len(prefix)] == prefix {
			out = append(out, v)
		}
	}
	return out
}

// --- Recall Submissions ---

func (m *MemoryStore) CreateRecallSubmission(paymentID, recallID string, s models.RecallSubmission) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	k := key3(paymentID, recallID, s.ID)
	if _, ok := m.recallSubmissions[k]; ok {
		return ErrConflict
	}
	m.recallSubmissions[k] = s
	return nil
}

func (m *MemoryStore) GetRecallSubmission(paymentID, recallID, submissionID string) (models.RecallSubmission, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	s, ok := m.recallSubmissions[key3(paymentID, recallID, submissionID)]
	if !ok {
		return s, ErrNotFound
	}
	return s, nil
}

func (m *MemoryStore) ListRecallSubmissions(paymentID, recallID string) []models.RecallSubmission {
	m.mu.RLock()
	defer m.mu.RUnlock()
	prefix := key2(paymentID, recallID) + ":"
	var out []models.RecallSubmission
	for k, v := range m.recallSubmissions {
		if len(k) > len(prefix) && k[:len(prefix)] == prefix {
			out = append(out, v)
		}
	}
	return out
}

func (m *MemoryStore) UpdateRecallSubmission(paymentID, recallID string, s models.RecallSubmission) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	k := key3(paymentID, recallID, s.ID)
	if _, ok := m.recallSubmissions[k]; !ok {
		return ErrNotFound
	}
	m.recallSubmissions[k] = s
	return nil
}

// --- Recall Decisions ---

func (m *MemoryStore) CreateRecallDecision(paymentID, recallID string, d models.RecallDecision) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	k := key3(paymentID, recallID, d.ID)
	if _, ok := m.recallDecisions[k]; ok {
		return ErrConflict
	}
	m.recallDecisions[k] = d
	return nil
}

func (m *MemoryStore) GetRecallDecision(paymentID, recallID, decisionID string) (models.RecallDecision, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	d, ok := m.recallDecisions[key3(paymentID, recallID, decisionID)]
	if !ok {
		return d, ErrNotFound
	}
	return d, nil
}

func (m *MemoryStore) ListRecallDecisions(paymentID, recallID string) []models.RecallDecision {
	m.mu.RLock()
	defer m.mu.RUnlock()
	prefix := key2(paymentID, recallID) + ":"
	var out []models.RecallDecision
	for k, v := range m.recallDecisions {
		if len(k) > len(prefix) && k[:len(prefix)] == prefix {
			out = append(out, v)
		}
	}
	return out
}

// --- Recall Decision Submissions ---

func (m *MemoryStore) CreateRecallDecisionSubmission(paymentID, recallID, decisionID string, s models.RecallDecisionSubmission) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	k := key4(paymentID, recallID, decisionID, s.ID)
	if _, ok := m.recallDecisionSubmissions[k]; ok {
		return ErrConflict
	}
	m.recallDecisionSubmissions[k] = s
	return nil
}

func (m *MemoryStore) GetRecallDecisionSubmission(paymentID, recallID, decisionID, submissionID string) (models.RecallDecisionSubmission, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	s, ok := m.recallDecisionSubmissions[key4(paymentID, recallID, decisionID, submissionID)]
	if !ok {
		return s, ErrNotFound
	}
	return s, nil
}

func (m *MemoryStore) ListRecallDecisionSubmissions(paymentID, recallID, decisionID string) []models.RecallDecisionSubmission {
	m.mu.RLock()
	defer m.mu.RUnlock()
	prefix := key3(paymentID, recallID, decisionID) + ":"
	var out []models.RecallDecisionSubmission
	for k, v := range m.recallDecisionSubmissions {
		if len(k) > len(prefix) && k[:len(prefix)] == prefix {
			out = append(out, v)
		}
	}
	return out
}

func (m *MemoryStore) UpdateRecallDecisionSubmission(paymentID, recallID, decisionID string, s models.RecallDecisionSubmission) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	k := key4(paymentID, recallID, decisionID, s.ID)
	if _, ok := m.recallDecisionSubmissions[k]; !ok {
		return ErrNotFound
	}
	m.recallDecisionSubmissions[k] = s
	return nil
}

// --- Reversals ---

func (m *MemoryStore) CreateReversal(paymentID string, r models.Reversal) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	k := key2(paymentID, r.ID)
	if _, ok := m.reversals[k]; ok {
		return ErrConflict
	}
	m.reversals[k] = r
	return nil
}

func (m *MemoryStore) GetReversal(paymentID, reversalID string) (models.Reversal, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	r, ok := m.reversals[key2(paymentID, reversalID)]
	if !ok {
		return r, ErrNotFound
	}
	return r, nil
}

func (m *MemoryStore) ListReversals(paymentID string) []models.Reversal {
	m.mu.RLock()
	defer m.mu.RUnlock()
	prefix := paymentID + ":"
	var out []models.Reversal
	for k, v := range m.reversals {
		if len(k) > len(prefix) && k[:len(prefix)] == prefix {
			out = append(out, v)
		}
	}
	return out
}

// --- Reversal Submissions ---

func (m *MemoryStore) CreateReversalSubmission(paymentID, reversalID string, s models.ReversalSubmission) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	k := key3(paymentID, reversalID, s.ID)
	if _, ok := m.reversalSubmissions[k]; ok {
		return ErrConflict
	}
	m.reversalSubmissions[k] = s
	return nil
}

func (m *MemoryStore) GetReversalSubmission(paymentID, reversalID, submissionID string) (models.ReversalSubmission, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	s, ok := m.reversalSubmissions[key3(paymentID, reversalID, submissionID)]
	if !ok {
		return s, ErrNotFound
	}
	return s, nil
}

func (m *MemoryStore) ListReversalSubmissions(paymentID, reversalID string) []models.ReversalSubmission {
	m.mu.RLock()
	defer m.mu.RUnlock()
	prefix := key2(paymentID, reversalID) + ":"
	var out []models.ReversalSubmission
	for k, v := range m.reversalSubmissions {
		if len(k) > len(prefix) && k[:len(prefix)] == prefix {
			out = append(out, v)
		}
	}
	return out
}

func (m *MemoryStore) UpdateReversalSubmission(paymentID, reversalID string, s models.ReversalSubmission) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	k := key3(paymentID, reversalID, s.ID)
	if _, ok := m.reversalSubmissions[k]; !ok {
		return ErrNotFound
	}
	m.reversalSubmissions[k] = s
	return nil
}

// --- Subscriptions ---

func (m *MemoryStore) CreateSubscription(s models.Subscription) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.subscriptions[s.ID]; ok {
		return ErrConflict
	}
	m.subscriptions[s.ID] = s
	return nil
}

func (m *MemoryStore) GetSubscription(id string) (models.Subscription, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	s, ok := m.subscriptions[id]
	if !ok {
		return s, ErrNotFound
	}
	return s, nil
}

func (m *MemoryStore) ListSubscriptions() []models.Subscription {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]models.Subscription, 0, len(m.subscriptions))
	for _, s := range m.subscriptions {
		out = append(out, s)
	}
	return out
}

func (m *MemoryStore) UpdateSubscription(s models.Subscription) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.subscriptions[s.ID]; !ok {
		return ErrNotFound
	}
	m.subscriptions[s.ID] = s
	return nil
}

func (m *MemoryStore) DeleteSubscription(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.subscriptions[id]; !ok {
		return ErrNotFound
	}
	delete(m.subscriptions, id)
	return nil
}

func (m *MemoryStore) MatchSubscriptions(recordType, eventType string) []models.Subscription {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var out []models.Subscription
	for _, s := range m.subscriptions {
		if s.Attributes.IsActive && s.Attributes.RecordType == recordType && s.Attributes.EventType == eventType {
			out = append(out, s)
		}
	}
	return out
}
