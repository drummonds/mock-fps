package models

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
