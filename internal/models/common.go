package models

import "time"

// Resource is the base for all JSON:API resources.
type Resource struct {
	Type           string    `json:"type"`
	ID             string    `json:"id"`
	OrganisationID string    `json:"organisation_id"`
	Version        int       `json:"version"`
	CreatedOn      time.Time `json:"created_on"`
	ModifiedOn     time.Time `json:"modified_on"`
}

// AccountParty represents a debtor or beneficiary party.
type AccountParty struct {
	AccountName   string `json:"account_name,omitempty"`
	AccountNumber string `json:"account_number,omitempty"`
	SortCode      string `json:"sort_code,omitempty"`
}

// ChargesInformation holds charges data.
type ChargesInformation struct {
	BearerCode string `json:"bearer_code,omitempty"`
}

// FxInfo holds foreign exchange information.
type FxInfo struct {
	ContractReference string  `json:"contract_reference,omitempty"`
	ExchangeRate      string  `json:"exchange_rate,omitempty"`
	OriginalAmount    string  `json:"original_amount,omitempty"`
	OriginalCurrency  string  `json:"original_currency,omitempty"`
}

// RelationshipData is a JSON:API relationship link.
type RelationshipData struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

// Relationship wraps a single relationship.
type Relationship struct {
	Data []RelationshipData `json:"data"`
}
