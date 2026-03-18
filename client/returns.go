package client

import "fmt"

// CreateReturn creates a return on a payment.
func (c *Client) CreateReturn(paymentID string, ret ReturnPayment) (*ReturnPayment, error) {
	var env DataEnvelope[ReturnPayment]
	err := c.doJSON("POST",
		fmt.Sprintf("%s/%s/returns", paymentsPath, paymentID),
		DataEnvelope[ReturnPayment]{Data: ret}, &env)
	if err != nil {
		return nil, err
	}
	return &env.Data, nil
}

// GetReturn retrieves a specific return on a payment.
func (c *Client) GetReturn(paymentID, returnID string) (*ReturnPayment, error) {
	var env DataEnvelope[ReturnPayment]
	err := c.doJSON("GET",
		fmt.Sprintf("%s/%s/returns/%s", paymentsPath, paymentID, returnID),
		nil, &env)
	if err != nil {
		return nil, err
	}
	return &env.Data, nil
}

// ListReturns lists all returns on a payment.
func (c *Client) ListReturns(paymentID string) ([]ReturnPayment, error) {
	var env ListEnvelope[ReturnPayment]
	err := c.doJSON("GET",
		fmt.Sprintf("%s/%s/returns", paymentsPath, paymentID),
		nil, &env)
	if err != nil {
		return nil, err
	}
	return env.Data, nil
}

// CreateReturnSubmission creates a submission for a return.
func (c *Client) CreateReturnSubmission(paymentID, returnID string, sub ReturnSubmission) (*ReturnSubmission, error) {
	var env DataEnvelope[ReturnSubmission]
	err := c.doJSON("POST",
		fmt.Sprintf("%s/%s/returns/%s/submissions", paymentsPath, paymentID, returnID),
		DataEnvelope[ReturnSubmission]{Data: sub}, &env)
	if err != nil {
		return nil, err
	}
	return &env.Data, nil
}

// GetReturnSubmission retrieves a specific return submission.
func (c *Client) GetReturnSubmission(paymentID, returnID, submissionID string) (*ReturnSubmission, error) {
	var env DataEnvelope[ReturnSubmission]
	err := c.doJSON("GET",
		fmt.Sprintf("%s/%s/returns/%s/submissions/%s", paymentsPath, paymentID, returnID, submissionID),
		nil, &env)
	if err != nil {
		return nil, err
	}
	return &env.Data, nil
}
