package models

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
