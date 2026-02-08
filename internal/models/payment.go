package models

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
