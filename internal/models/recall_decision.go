package models

// RecallDecision represents a recall decision resource.
type RecallDecision struct {
	Resource
	Attributes RecallDecisionAttributes `json:"attributes"`
}

// RecallDecisionAttributes holds recall decision data.
type RecallDecisionAttributes struct {
	Answer       string `json:"answer,omitempty"`
	Reason       string `json:"reason,omitempty"`
	Status       string `json:"status,omitempty"`
}
