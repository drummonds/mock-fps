package jsonapi

// DataEnvelope wraps a single resource in JSON:API format.
type DataEnvelope[T any] struct {
	Data T `json:"data"`
}

// ListEnvelope wraps a collection of resources in JSON:API format.
type ListEnvelope[T any] struct {
	Data []T `json:"data"`
}

// ErrorResponse is a JSON:API error response.
type ErrorResponse struct {
	Errors []Error `json:"errors"`
}

// Error is a single JSON:API error object.
type Error struct {
	Status string `json:"status"`
	Title  string `json:"title"`
	Detail string `json:"detail,omitempty"`
}
