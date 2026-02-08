package models

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
