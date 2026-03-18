// Package fpsevents defines the complete set of FPS payment events,
// status chains, and webhook types for programmatic consumption.
//
// Import this package to enumerate all possible FPS scenarios for testing:
//
//	import "github.com/nibble/mock-fps/fpsevents"
//
//	for _, evt := range fpsevents.WebhookEvents { ... }
package fpsevents

// PaymentType represents an FPS payment type.
type PaymentType struct {
	Name     string // e.g. "Single Immediate Payment"
	Code     string // e.g. "SIP"
	Delivery string // e.g. "Synchronous, <15s round-trip"
}

// StatusChain describes the ordered status transitions for a resource type.
type StatusChain struct {
	Resource   string   // e.g. "payment_submission"
	Statuses   []string // ordered happy-path statuses
	FailStatus string   // terminal failure status, if any
}

// WebhookEvent describes a Form3 webhook event type.
type WebhookEvent struct {
	RecordType string // e.g. "payment_admissions"
	EventType  string // e.g. "created"
	When       string // human description of when this fires
}

// AcceptanceQualifier describes a stand-in acceptance qualifier value.
type AcceptanceQualifier struct {
	Value       string
	Description string
}

// RecallAnswer describes a possible recall decision answer.
type RecallAnswer struct {
	Value       string
	Description string
}

// SettlementCycle describes an FPS settlement cycle.
type SettlementCycle struct {
	Cycle int
	Time  string // approximate time, e.g. "07:15"
}

// ReconciliationFilter describes a query parameter for payment reconciliation filtering.
type ReconciliationFilter struct {
	Param       string // e.g. "filter[submission.settlement_cycle]"
	Type        string // e.g. "int"
	Description string
}
