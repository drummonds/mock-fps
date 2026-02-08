package models

// PaymentSubmission represents a payment submission resource.
type PaymentSubmission struct {
	Resource
	Attributes PaymentSubmissionAttributes `json:"attributes"`
}

// PaymentSubmissionAttributes holds submission data.
type PaymentSubmissionAttributes struct {
	Status         string `json:"status"`
	SubmissionDate string `json:"submission_date,omitempty"`
	SchemeStatusCode string `json:"scheme_status_code,omitempty"`
}
