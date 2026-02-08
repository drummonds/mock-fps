package models

// ReturnPayment represents a payment return resource.
type ReturnPayment struct {
	Resource
	Attributes ReturnPaymentAttributes `json:"attributes"`
}

// ReturnPaymentAttributes holds return data.
type ReturnPaymentAttributes struct {
	Amount        string `json:"amount"`
	Currency      string `json:"currency"`
	ReturnCode    string `json:"return_code,omitempty"`
	ReturnReason  string `json:"return_reason,omitempty"`
	SchemeStatus  string `json:"scheme_status,omitempty"`
}
