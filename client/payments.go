package client

import (
	"fmt"
	"net/url"
	"strconv"
)

const paymentsPath = "/v1/transaction/payments"

// ListPaymentsOptions holds optional filter and pagination parameters.
type ListPaymentsOptions struct {
	// Filter by submission settlement cycle.
	SubmissionSettlementCycle *int
	// Filter by submission settlement date (YYYY-MM-DD).
	SubmissionSettlementDate string
	// Filter by admission settlement cycle.
	AdmissionSettlementCycle *int
	// Filter by admission settlement date (YYYY-MM-DD).
	AdmissionSettlementDate string
	// Page size (default 100).
	PageSize int
	// Page number (default 1).
	PageNumber int
}

// CreatePayment creates a new payment.
func (c *Client) CreatePayment(payment Payment) (*Payment, error) {
	var env DataEnvelope[Payment]
	err := c.doJSON("POST", paymentsPath, DataEnvelope[Payment]{Data: payment}, &env)
	if err != nil {
		return nil, err
	}
	return &env.Data, nil
}

// GetPayment retrieves a payment by ID.
func (c *Client) GetPayment(id string) (*Payment, error) {
	var env DataEnvelope[Payment]
	err := c.doJSON("GET", fmt.Sprintf("%s/%s", paymentsPath, id), nil, &env)
	if err != nil {
		return nil, err
	}
	return &env.Data, nil
}

// ListPayments lists payments with optional filters and pagination.
func (c *Client) ListPayments(opts *ListPaymentsOptions) ([]Payment, error) {
	path := paymentsPath
	if opts != nil {
		q := url.Values{}
		if opts.SubmissionSettlementCycle != nil {
			q.Set("filter[submission.settlement_cycle]", strconv.Itoa(*opts.SubmissionSettlementCycle))
		}
		if opts.SubmissionSettlementDate != "" {
			q.Set("filter[submission.settlement_date]", opts.SubmissionSettlementDate)
		}
		if opts.AdmissionSettlementCycle != nil {
			q.Set("filter[admission.settlement_cycle]", strconv.Itoa(*opts.AdmissionSettlementCycle))
		}
		if opts.AdmissionSettlementDate != "" {
			q.Set("filter[admission.settlement_date]", opts.AdmissionSettlementDate)
		}
		if opts.PageSize > 0 {
			q.Set("page[size]", strconv.Itoa(opts.PageSize))
		}
		if opts.PageNumber > 0 {
			q.Set("page[number]", strconv.Itoa(opts.PageNumber))
		}
		if encoded := q.Encode(); encoded != "" {
			path += "?" + encoded
		}
	}
	var env ListEnvelope[Payment]
	err := c.doJSON("GET", path, nil, &env)
	if err != nil {
		return nil, err
	}
	return env.Data, nil
}
