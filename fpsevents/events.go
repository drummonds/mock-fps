package fpsevents

// PaymentTypes lists all FPS payment types.
var PaymentTypes = []PaymentType{
	{"Single Immediate Payment", "SIP", "Synchronous, <15s round-trip"},
	{"Standing Order Payment", "SOP", "Asynchronous"},
	{"Forward Dated Payment", "FDP", "Asynchronous"},
	{"Direct Corporate Access", "DCA", "Batch via Secure-IP"},
}

// StatusChains maps resource type to its status transition chain.
var StatusChains = map[string]StatusChain{
	"payment_submission": {
		Resource: "payment_submission",
		Statuses: []string{
			"accepted",
			"validation_pending",
			"limit_check_pending",
			"limit_check_passed",
			"released_to_gateway",
			"queued_for_delivery",
			"submitted",
			"delivery_confirmed",
		},
		FailStatus: "delivery_failed",
	},
	"payment_admission": {
		Resource:   "payment_admission",
		Statuses:   []string{"pending", "confirmed"},
		FailStatus: "failed",
	},
	"return_submission": {
		Resource: "return_submission",
		Statuses: []string{
			"accepted",
			"validation_pending",
			"limit_check_pending",
			"limit_check_passed",
			"released_to_gateway",
			"queued_for_delivery",
			"delivery_confirmed",
		},
		FailStatus: "delivery_failed",
	},
	"recall_decision_submission": {
		Resource: "recall_decision_submission",
		Statuses: []string{
			"accepted",
			"validation_pending",
			"validation_passed",
			"limit_check_pending",
			"limit_check_passed",
			"released_to_gateway",
			"queued_for_delivery",
			"delivery_confirmed",
		},
		FailStatus: "delivery_failed",
	},
	"recall_decision_admission": {
		Resource:   "recall_decision_admission",
		Statuses:   []string{"confirmed"},
		FailStatus: "failed",
	},
	"reversal_admission": {
		Resource: "reversal_admission",
		Statuses: []string{"pending", "confirmed"},
	},
}

// WebhookEvents lists all Form3 FPS webhook event types.
var WebhookEvents = []WebhookEvent{
	{"payments", "created", "Payment resource created"},
	{"payment_admissions", "created", "Inbound payment received"},
	{"payment_admissions", "updated", "Admission status changed"},
	{"payment_admission_tasks", "created", "Admission task requires customer action"},
	{"payment_admission_tasks", "updated", "Admission task status changed"},
	{"payment_submissions", "created", "Outbound submission created"},
	{"payment_submissions", "updated", "Submission status changed"},
	{"return_payments", "created", "Return created"},
	{"return_submissions", "updated", "Return submission status changed"},
	{"recalls", "created", "Recall request received"},
	{"recall_submissions", "updated", "Recall submission status changed"},
	{"recall_decisions", "created", "Recall decision created"},
	{"recall_decision_submissions", "updated", "Decision submission status changed"},
	{"reversals", "created", "Reversal created"},
	{"reversal_admissions", "created", "Reversal received"},
	{"reversal_submissions", "updated", "Reversal submission status changed"},
	{"reports", "created", "Settlement report available"},
}

// AdmissionStatusReasons lists all possible PaymentAdmission status_reason values
// from the Form3 OpenAPI spec.
var AdmissionStatusReasons = []string{
	// Success
	"accepted",

	// Account-related
	"account_closed",
	"account_closed_beneficiary_deceased",
	"account_closed_beneficiary_sensitivities",
	"account_closed_business_reasons",
	"account_closed_currency",
	"account_closed_stopped",
	"account_closed_terms_and_conditions",
	"account_closed_transferred",
	"blocked_account",
	"unknown_accountnumber",

	// Agent/infrastructure errors
	"agent_clearing_process_error",
	"agent_clearing_process_timeout",
	"agent_reason_unknown",
	"agent_suspended",
	"agent_unavailable",
	"beneficiary_agent_clearing_process_error",
	"beneficiary_agent_clearing_process_timeout",
	"beneficiary_agent_suspended",
	"beneficiary_agent_unavailable",

	// Validation
	"amount_exceeds_settlement_limit",
	"amount_invalid_or_missing",
	"amount_not_allowed",
	"bankid_not_provisioned",
	"invalid_bank_ID",
	"invalid_bank_operation_code",
	"invalid_beneficiary_address",
	"invalid_beneficiary_agent_BIC",
	"invalid_beneficiary_details",
	"invalid_debtor_agent_BIC",
	"invalid_debtor_details",
	"end_to_end_id_missing_or_invalid",
	"incorrect_reference_reference_mask",
	"incorrect_reference_secondary_identification",
	"incorrect_reference_validation_type",
	"beneficiary_name_not_present",

	// Other
	"business_reasons",
	"customer_check_failed",
	"customer_reason_unknown",
	"duplicate_payment",
	"original_payment_not_received",
	"regulatory_reason",
	"scheme_timeout",
	"transaction_forbidden",
	"transaction_type_not_supported",
	"rejected_by_customer",
}

// AcceptanceQualifiers lists Form3 acceptance_qualifier values for stand-in payments.
var AcceptanceQualifiers = []AcceptanceQualifier{
	{"none", "Standard accept"},
	{"same_day", "Credit by end of day"},
	{"next_calendar_day", "Credit by next calendar day"},
	{"next_working_day", "Credit by next working day"},
	{"after_next_working_day", "Credit after next working day"},
	{"some_other_time", "Credit at some other time"},
}

// PostingStatuses lists all possible posting_status values on admissions.
var PostingStatuses = []string{
	"pending",
	"posted",
	"rejected",
	"blocked",
	"passed_without_posting",
	"received",
	"accepted_funds_checked",
	"accepted_settlement_in_process",
	"accepted_settlement_complete",
	"accepted_settlement_complete_creditor_account",
}

// AccountValidationOutcomes lists all account_validation_outcome values.
var AccountValidationOutcomes = []string{
	"passed",
	"failed",
	"failed_auto_reject_disabled",
	"failed_auto_reject_enabled",
}

// RecallAnswers lists all possible recall decision answer values.
var RecallAnswers = []RecallAnswer{
	{"accepted", "Agree to return funds"},
	{"rejected", "Refuse to return funds"},
	{"pending", "Still considering"},
	{"partially_accepted", "Return partial amount"},
	{"payment_cancelled", "Payment was already cancelled"},
}

// SettlementCycles lists the 3 daily FPS settlement cycles.
var SettlementCycles = []SettlementCycle{
	{1, "07:15"},
	{2, "13:00"},
	{3, "15:45"},
}

// ReconciliationFilters lists the query parameters supported by
// GET /v1/transaction/payments for reconciliation filtering.
var ReconciliationFilters = []ReconciliationFilter{
	{"filter[submission.settlement_cycle]", "int", "Settlement cycle on the payment submission"},
	{"filter[submission.settlement_date]", "string (YYYY-MM-DD)", "Settlement date on the payment submission"},
	{"filter[admission.settlement_cycle]", "int", "Settlement cycle on the payment admission"},
	{"filter[admission.settlement_date]", "string (YYYY-MM-DD)", "Settlement date on the payment admission"},
	{"page[number]", "int", "Page number (1-based, default 1)"},
	{"page[size]", "int", "Page size (default 100)"},
}
