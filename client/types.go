package client

import "time"

// DataEnvelope wraps a single resource in JSON:API format.
type DataEnvelope[T any] struct {
	Data T `json:"data"`
}

// ListEnvelope wraps a collection of resources in JSON:API format.
type ListEnvelope[T any] struct {
	Data []T `json:"data"`
}

// Resource is the base for all JSON:API resources.
type Resource struct {
	Type           string    `json:"type"`
	ID             string    `json:"id"`
	OrganisationID string    `json:"organisation_id"`
	Version        int       `json:"version"`
	CreatedOn      time.Time `json:"created_on"`
	ModifiedOn     time.Time `json:"modified_on"`
}

// AccountParty represents a debtor or beneficiary party.
type AccountParty struct {
	AccountName   string `json:"account_name,omitempty"`
	AccountNumber string `json:"account_number,omitempty"`
	SortCode      string `json:"sort_code,omitempty"`
}

// ChargesInformation holds charges data.
type ChargesInformation struct {
	BearerCode string `json:"bearer_code,omitempty"`
}

// FxInfo holds foreign exchange information.
type FxInfo struct {
	ContractReference string `json:"contract_reference,omitempty"`
	ExchangeRate      string `json:"exchange_rate,omitempty"`
	OriginalAmount    string `json:"original_amount,omitempty"`
	OriginalCurrency  string `json:"original_currency,omitempty"`
}

// RelationshipData is a JSON:API relationship link.
type RelationshipData struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

// Relationship wraps a single relationship.
type Relationship struct {
	Data []RelationshipData `json:"data"`
}

// Payment represents a payment resource.
type Payment struct {
	Resource
	Attributes    PaymentAttributes     `json:"attributes"`
	Relationships *PaymentRelationships `json:"relationships,omitempty"`
}

// PaymentAttributes holds the payment data fields.
type PaymentAttributes struct {
	Amount               string              `json:"amount"`
	Currency             string              `json:"currency"`
	EndToEndReference    string              `json:"end_to_end_reference,omitempty"`
	NumericReference     string              `json:"numeric_reference,omitempty"`
	PaymentScheme        string              `json:"payment_scheme,omitempty"`
	PaymentType          string              `json:"payment_type,omitempty"`
	ProcessingDate       string              `json:"processing_date,omitempty"`
	Reference            string              `json:"reference,omitempty"`
	SchemePaymentSubType string              `json:"scheme_payment_sub_type,omitempty"`
	SchemePaymentType    string              `json:"scheme_payment_type,omitempty"`
	BeneficiaryParty     *AccountParty       `json:"beneficiary_party,omitempty"`
	DebtorParty          *AccountParty       `json:"debtor_party,omitempty"`
	ChargesInformation   *ChargesInformation `json:"charges_information,omitempty"`
	Fx                   *FxInfo             `json:"fx,omitempty"`
}

// PaymentRelationships holds relationships to submissions, admissions, etc.
type PaymentRelationships struct {
	PaymentSubmissions *Relationship `json:"payment_submissions,omitempty"`
	PaymentAdmissions  *Relationship `json:"payment_admissions,omitempty"`
	PaymentReturns     *Relationship `json:"payment_returns,omitempty"`
	PaymentRecalls     *Relationship `json:"payment_recalls,omitempty"`
	PaymentReversals   *Relationship `json:"payment_reversals,omitempty"`
}

// PaymentSubmission represents a payment submission resource.
type PaymentSubmission struct {
	Resource
	Attributes PaymentSubmissionAttributes `json:"attributes"`
}

// PaymentSubmissionAttributes holds submission data.
type PaymentSubmissionAttributes struct {
	Status              string `json:"status"`
	SubmissionDate      string `json:"submission_date,omitempty"`
	SchemeStatusCode    string `json:"scheme_status_code,omitempty"`
	SchemeTransactionID string `json:"scheme_transaction_id,omitempty"`
	SettlementCycle     int    `json:"settlement_cycle,omitempty"`
	SettlementDate      string `json:"settlement_date,omitempty"`
}

// PaymentAdmission represents a payment admission resource.
type PaymentAdmission struct {
	Resource
	Attributes PaymentAdmissionAttributes `json:"attributes"`
}

// PaymentAdmissionAttributes holds admission data.
type PaymentAdmissionAttributes struct {
	Status              string `json:"status"`
	AdmissionDate       string `json:"admission_date,omitempty"`
	SchemeTransactionID string `json:"scheme_transaction_id,omitempty"`
	SettlementCycle     int    `json:"settlement_cycle,omitempty"`
	SettlementDate      string `json:"settlement_date,omitempty"`
}

// AdmissionTask represents a task on an admission.
type AdmissionTask struct {
	Resource
	Attributes AdmissionTaskAttributes `json:"attributes"`
}

// AdmissionTaskAttributes holds task data.
type AdmissionTaskAttributes struct {
	Status   string `json:"status"`
	Assignee string `json:"assignee,omitempty"`
	Name     string `json:"name,omitempty"`
}

// ReturnPayment represents a payment return resource.
type ReturnPayment struct {
	Resource
	Attributes ReturnPaymentAttributes `json:"attributes"`
}

// ReturnPaymentAttributes holds return data.
type ReturnPaymentAttributes struct {
	Amount       string `json:"amount"`
	Currency     string `json:"currency"`
	ReturnCode   string `json:"return_code,omitempty"`
	ReturnReason string `json:"return_reason,omitempty"`
	SchemeStatus string `json:"scheme_status,omitempty"`
}

// ReturnSubmission represents a return submission resource.
type ReturnSubmission struct {
	Resource
	Attributes ReturnSubmissionAttributes `json:"attributes"`
}

// ReturnSubmissionAttributes holds return submission data.
type ReturnSubmissionAttributes struct {
	Status         string `json:"status"`
	SubmissionDate string `json:"submission_date,omitempty"`
}

// Recall represents a payment recall resource.
type Recall struct {
	Resource
	Attributes RecallAttributes `json:"attributes"`
}

// RecallAttributes holds recall data.
type RecallAttributes struct {
	Amount       string `json:"amount,omitempty"`
	Currency     string `json:"currency,omitempty"`
	RecallReason string `json:"recall_reason,omitempty"`
	RecallType   string `json:"recall_type,omitempty"`
	Status       string `json:"status,omitempty"`
}

// RecallSubmission represents a recall submission resource.
type RecallSubmission struct {
	Resource
	Attributes RecallSubmissionAttributes `json:"attributes"`
}

// RecallSubmissionAttributes holds recall submission data.
type RecallSubmissionAttributes struct {
	Status         string `json:"status"`
	SubmissionDate string `json:"submission_date,omitempty"`
}

// RecallDecision represents a recall decision resource.
type RecallDecision struct {
	Resource
	Attributes RecallDecisionAttributes `json:"attributes"`
}

// RecallDecisionAttributes holds recall decision data.
type RecallDecisionAttributes struct {
	Answer string `json:"answer,omitempty"`
	Reason string `json:"reason,omitempty"`
	Status string `json:"status,omitempty"`
}

// RecallDecisionSubmission represents a recall decision submission resource.
type RecallDecisionSubmission struct {
	Resource
	Attributes RecallDecisionSubmissionAttributes `json:"attributes"`
}

// RecallDecisionSubmissionAttributes holds recall decision submission data.
type RecallDecisionSubmissionAttributes struct {
	Status         string `json:"status"`
	SubmissionDate string `json:"submission_date,omitempty"`
}

// Reversal represents a payment reversal resource.
type Reversal struct {
	Resource
	Attributes ReversalAttributes `json:"attributes"`
}

// ReversalAttributes holds reversal data.
type ReversalAttributes struct {
	Amount         string `json:"amount,omitempty"`
	Currency       string `json:"currency,omitempty"`
	ReversalReason string `json:"reversal_reason,omitempty"`
	Status         string `json:"status,omitempty"`
}

// ReversalSubmission represents a reversal submission resource.
type ReversalSubmission struct {
	Resource
	Attributes ReversalSubmissionAttributes `json:"attributes"`
}

// ReversalSubmissionAttributes holds reversal submission data.
type ReversalSubmissionAttributes struct {
	Status         string `json:"status"`
	SubmissionDate string `json:"submission_date,omitempty"`
}

// Subscription represents a webhook subscription resource.
type Subscription struct {
	Resource
	Attributes SubscriptionAttributes `json:"attributes"`
}

// SubscriptionAttributes holds subscription data.
type SubscriptionAttributes struct {
	CallbackURI       string `json:"callback_uri"`
	EventType         string `json:"event_type"`
	RecordType        string `json:"record_type"`
	IsActive          bool   `json:"is_active"`
	CallbackTransport string `json:"callback_transport,omitempty"`
	UserID            string `json:"user_id,omitempty"`
}

// Notification is the payload sent to webhook subscribers.
type Notification struct {
	ID             string           `json:"id"`
	OrganisationID string           `json:"organisation_id"`
	Type           string           `json:"type"`
	Version        int              `json:"version"`
	CreatedOn      time.Time        `json:"created_on"`
	Data           NotificationData `json:"data"`
}

// NotificationData holds the event details.
type NotificationData struct {
	RecordType string `json:"record_type"`
	EventType  string `json:"event_type"`
	ResourceID string `json:"resource_id"`
	Payload    any    `json:"payload,omitempty"`
}
