package client

import "fmt"

// CreatePaymentSubmission creates a submission for a payment.
func (c *Client) CreatePaymentSubmission(paymentID string, sub PaymentSubmission) (*PaymentSubmission, error) {
	var env DataEnvelope[PaymentSubmission]
	err := c.doJSON("POST",
		fmt.Sprintf("%s/%s/submissions", paymentsPath, paymentID),
		DataEnvelope[PaymentSubmission]{Data: sub}, &env)
	if err != nil {
		return nil, err
	}
	return &env.Data, nil
}

// GetPaymentSubmission retrieves a specific submission for a payment.
func (c *Client) GetPaymentSubmission(paymentID, submissionID string) (*PaymentSubmission, error) {
	var env DataEnvelope[PaymentSubmission]
	err := c.doJSON("GET",
		fmt.Sprintf("%s/%s/submissions/%s", paymentsPath, paymentID, submissionID),
		nil, &env)
	if err != nil {
		return nil, err
	}
	return &env.Data, nil
}
