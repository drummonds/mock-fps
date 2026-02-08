package store

import "github.com/nibble/mock-fps/internal/models"

// Store defines the storage interface for all resources.
type Store interface {
	// Payments
	CreatePayment(p models.Payment) error
	GetPayment(id string) (models.Payment, error)
	ListPayments() []models.Payment
	UpdatePayment(p models.Payment) error

	// Payment Submissions
	CreatePaymentSubmission(paymentID string, s models.PaymentSubmission) error
	GetPaymentSubmission(paymentID, submissionID string) (models.PaymentSubmission, error)
	ListPaymentSubmissions(paymentID string) []models.PaymentSubmission
	UpdatePaymentSubmission(paymentID string, s models.PaymentSubmission) error

	// Payment Admissions
	CreatePaymentAdmission(paymentID string, a models.PaymentAdmission) error
	GetPaymentAdmission(paymentID, admissionID string) (models.PaymentAdmission, error)
	ListPaymentAdmissions(paymentID string) []models.PaymentAdmission
	UpdatePaymentAdmission(paymentID string, a models.PaymentAdmission) error

	// Admission Tasks
	CreateAdmissionTask(paymentID, admissionID string, t models.AdmissionTask) error
	GetAdmissionTask(paymentID, admissionID, taskID string) (models.AdmissionTask, error)
	UpdateAdmissionTask(paymentID, admissionID string, t models.AdmissionTask) error

	// Returns
	CreateReturn(paymentID string, r models.ReturnPayment) error
	GetReturn(paymentID, returnID string) (models.ReturnPayment, error)
	ListReturns(paymentID string) []models.ReturnPayment

	// Return Submissions
	CreateReturnSubmission(paymentID, returnID string, s models.ReturnSubmission) error
	GetReturnSubmission(paymentID, returnID, submissionID string) (models.ReturnSubmission, error)
	ListReturnSubmissions(paymentID, returnID string) []models.ReturnSubmission
	UpdateReturnSubmission(paymentID, returnID string, s models.ReturnSubmission) error

	// Recalls
	CreateRecall(paymentID string, r models.Recall) error
	GetRecall(paymentID, recallID string) (models.Recall, error)
	ListRecalls(paymentID string) []models.Recall

	// Recall Submissions
	CreateRecallSubmission(paymentID, recallID string, s models.RecallSubmission) error
	GetRecallSubmission(paymentID, recallID, submissionID string) (models.RecallSubmission, error)
	ListRecallSubmissions(paymentID, recallID string) []models.RecallSubmission
	UpdateRecallSubmission(paymentID, recallID string, s models.RecallSubmission) error

	// Recall Decisions
	CreateRecallDecision(paymentID, recallID string, d models.RecallDecision) error
	GetRecallDecision(paymentID, recallID, decisionID string) (models.RecallDecision, error)
	ListRecallDecisions(paymentID, recallID string) []models.RecallDecision

	// Recall Decision Submissions
	CreateRecallDecisionSubmission(paymentID, recallID, decisionID string, s models.RecallDecisionSubmission) error
	GetRecallDecisionSubmission(paymentID, recallID, decisionID, submissionID string) (models.RecallDecisionSubmission, error)
	ListRecallDecisionSubmissions(paymentID, recallID, decisionID string) []models.RecallDecisionSubmission
	UpdateRecallDecisionSubmission(paymentID, recallID, decisionID string, s models.RecallDecisionSubmission) error

	// Reversals
	CreateReversal(paymentID string, r models.Reversal) error
	GetReversal(paymentID, reversalID string) (models.Reversal, error)
	ListReversals(paymentID string) []models.Reversal

	// Reversal Submissions
	CreateReversalSubmission(paymentID, reversalID string, s models.ReversalSubmission) error
	GetReversalSubmission(paymentID, reversalID, submissionID string) (models.ReversalSubmission, error)
	ListReversalSubmissions(paymentID, reversalID string) []models.ReversalSubmission
	UpdateReversalSubmission(paymentID, reversalID string, s models.ReversalSubmission) error

	// Subscriptions
	CreateSubscription(s models.Subscription) error
	GetSubscription(id string) (models.Subscription, error)
	ListSubscriptions() []models.Subscription
	UpdateSubscription(s models.Subscription) error
	DeleteSubscription(id string) error
	MatchSubscriptions(recordType, eventType string) []models.Subscription
}
