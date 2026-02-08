package models

// Subscription represents a webhook subscription resource.
type Subscription struct {
	Resource
	Attributes SubscriptionAttributes `json:"attributes"`
}

// SubscriptionAttributes holds subscription data.
type SubscriptionAttributes struct {
	CallbackURI      string `json:"callback_uri"`
	EventType        string `json:"event_type"`
	RecordType       string `json:"record_type"`
	IsActive         bool   `json:"is_active"`
	CallbackTransport string `json:"callback_transport,omitempty"`
	UserID           string `json:"user_id,omitempty"`
}
