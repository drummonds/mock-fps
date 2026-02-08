package models

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
