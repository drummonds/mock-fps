package models

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
