package models

// Payment submission statuses - the lifecycle chain.
const (
	StatusAccepted           = "accepted"
	StatusValidationPending  = "validation_pending"
	StatusLimitCheckPending  = "limit_check_pending"
	StatusLimitCheckPassed   = "limit_check_passed"
	StatusReleasedToGateway  = "released_to_gateway"
	StatusQueuedForDelivery  = "queued_for_delivery"
	StatusSubmitted          = "submitted"
	StatusDeliveryConfirmed  = "delivery_confirmed"
	StatusDeliveryFailed     = "delivery_failed"
	StatusFailed             = "failed"

	// Admission statuses.
	StatusPending   = "pending"
	StatusConfirmed = "confirmed"

	// Return/Recall/Reversal submission statuses.
	StatusReturnAccepted          = "accepted"
	StatusReturnDeliveryConfirmed = "delivery_confirmed"
)

// Resource types for JSON:API.
const (
	ResourceTypePayment                  = "payments"
	ResourceTypePaymentSubmission        = "payment_submissions"
	ResourceTypePaymentAdmission         = "payment_admissions"
	ResourceTypeReturnPayment            = "return_payments"
	ResourceTypeReturnSubmission         = "return_submissions"
	ResourceTypeRecall                   = "recalls"
	ResourceTypeRecallSubmission         = "recall_submissions"
	ResourceTypeRecallDecision           = "recall_decisions"
	ResourceTypeRecallDecisionSubmission = "recall_decision_submissions"
	ResourceTypeReversal                 = "reversals"
	ResourceTypeReversalSubmission       = "reversal_submissions"
	ResourceTypeSubscription             = "subscriptions"
)

// Event types for webhook notifications.
const (
	EventCreated  = "created"
	EventUpdated  = "updated"
)
