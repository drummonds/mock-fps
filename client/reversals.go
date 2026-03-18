package client

import "fmt"

// CreateReversal creates a reversal on a payment.
func (c *Client) CreateReversal(paymentID string, rev Reversal) (*Reversal, error) {
	var env DataEnvelope[Reversal]
	err := c.doJSON("POST",
		fmt.Sprintf("%s/%s/reversals", paymentsPath, paymentID),
		DataEnvelope[Reversal]{Data: rev}, &env)
	if err != nil {
		return nil, err
	}
	return &env.Data, nil
}

// GetReversal retrieves a specific reversal on a payment.
func (c *Client) GetReversal(paymentID, reversalID string) (*Reversal, error) {
	var env DataEnvelope[Reversal]
	err := c.doJSON("GET",
		fmt.Sprintf("%s/%s/reversals/%s", paymentsPath, paymentID, reversalID),
		nil, &env)
	if err != nil {
		return nil, err
	}
	return &env.Data, nil
}

// ListReversals lists all reversals on a payment.
func (c *Client) ListReversals(paymentID string) ([]Reversal, error) {
	var env ListEnvelope[Reversal]
	err := c.doJSON("GET",
		fmt.Sprintf("%s/%s/reversals", paymentsPath, paymentID),
		nil, &env)
	if err != nil {
		return nil, err
	}
	return env.Data, nil
}

// CreateReversalSubmission creates a submission for a reversal.
func (c *Client) CreateReversalSubmission(paymentID, reversalID string, sub ReversalSubmission) (*ReversalSubmission, error) {
	var env DataEnvelope[ReversalSubmission]
	err := c.doJSON("POST",
		fmt.Sprintf("%s/%s/reversals/%s/submissions", paymentsPath, paymentID, reversalID),
		DataEnvelope[ReversalSubmission]{Data: sub}, &env)
	if err != nil {
		return nil, err
	}
	return &env.Data, nil
}

// GetReversalSubmission retrieves a specific reversal submission.
func (c *Client) GetReversalSubmission(paymentID, reversalID, submissionID string) (*ReversalSubmission, error) {
	var env DataEnvelope[ReversalSubmission]
	err := c.doJSON("GET",
		fmt.Sprintf("%s/%s/reversals/%s/submissions/%s", paymentsPath, paymentID, reversalID, submissionID),
		nil, &env)
	if err != nil {
		return nil, err
	}
	return &env.Data, nil
}
