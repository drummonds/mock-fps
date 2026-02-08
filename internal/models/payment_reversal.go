package models

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
