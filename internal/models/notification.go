package models

import "time"

// Notification is the payload sent to webhook subscribers.
type Notification struct {
	ID             string    `json:"id"`
	OrganisationID string    `json:"organisation_id"`
	Type           string    `json:"type"`
	Version        int       `json:"version"`
	CreatedOn      time.Time `json:"created_on"`
	Data           NotificationData `json:"data"`
}

// NotificationData holds the event details.
type NotificationData struct {
	RecordType string      `json:"record_type"`
	EventType  string      `json:"event_type"`
	ResourceID string      `json:"resource_id"`
	Payload    interface{} `json:"payload,omitempty"`
}
