package client

import "fmt"

// CreateRecall creates a recall on a payment.
func (c *Client) CreateRecall(paymentID string, recall Recall) (*Recall, error) {
	var env DataEnvelope[Recall]
	err := c.doJSON("POST",
		fmt.Sprintf("%s/%s/recalls", paymentsPath, paymentID),
		DataEnvelope[Recall]{Data: recall}, &env)
	if err != nil {
		return nil, err
	}
	return &env.Data, nil
}

// GetRecall retrieves a specific recall on a payment.
func (c *Client) GetRecall(paymentID, recallID string) (*Recall, error) {
	var env DataEnvelope[Recall]
	err := c.doJSON("GET",
		fmt.Sprintf("%s/%s/recalls/%s", paymentsPath, paymentID, recallID),
		nil, &env)
	if err != nil {
		return nil, err
	}
	return &env.Data, nil
}

// ListRecalls lists all recalls on a payment.
func (c *Client) ListRecalls(paymentID string) ([]Recall, error) {
	var env ListEnvelope[Recall]
	err := c.doJSON("GET",
		fmt.Sprintf("%s/%s/recalls", paymentsPath, paymentID),
		nil, &env)
	if err != nil {
		return nil, err
	}
	return env.Data, nil
}

// CreateRecallSubmission creates a submission for a recall.
func (c *Client) CreateRecallSubmission(paymentID, recallID string, sub RecallSubmission) (*RecallSubmission, error) {
	var env DataEnvelope[RecallSubmission]
	err := c.doJSON("POST",
		fmt.Sprintf("%s/%s/recalls/%s/submissions", paymentsPath, paymentID, recallID),
		DataEnvelope[RecallSubmission]{Data: sub}, &env)
	if err != nil {
		return nil, err
	}
	return &env.Data, nil
}

// GetRecallSubmission retrieves a specific recall submission.
func (c *Client) GetRecallSubmission(paymentID, recallID, submissionID string) (*RecallSubmission, error) {
	var env DataEnvelope[RecallSubmission]
	err := c.doJSON("GET",
		fmt.Sprintf("%s/%s/recalls/%s/submissions/%s", paymentsPath, paymentID, recallID, submissionID),
		nil, &env)
	if err != nil {
		return nil, err
	}
	return &env.Data, nil
}

// CreateRecallDecision creates a decision on a recall.
func (c *Client) CreateRecallDecision(paymentID, recallID string, dec RecallDecision) (*RecallDecision, error) {
	var env DataEnvelope[RecallDecision]
	err := c.doJSON("POST",
		fmt.Sprintf("%s/%s/recalls/%s/decisions", paymentsPath, paymentID, recallID),
		DataEnvelope[RecallDecision]{Data: dec}, &env)
	if err != nil {
		return nil, err
	}
	return &env.Data, nil
}

// GetRecallDecision retrieves a specific recall decision.
func (c *Client) GetRecallDecision(paymentID, recallID, decisionID string) (*RecallDecision, error) {
	var env DataEnvelope[RecallDecision]
	err := c.doJSON("GET",
		fmt.Sprintf("%s/%s/recalls/%s/decisions/%s", paymentsPath, paymentID, recallID, decisionID),
		nil, &env)
	if err != nil {
		return nil, err
	}
	return &env.Data, nil
}

// ListRecallDecisions lists all decisions on a recall.
func (c *Client) ListRecallDecisions(paymentID, recallID string) ([]RecallDecision, error) {
	var env ListEnvelope[RecallDecision]
	err := c.doJSON("GET",
		fmt.Sprintf("%s/%s/recalls/%s/decisions", paymentsPath, paymentID, recallID),
		nil, &env)
	if err != nil {
		return nil, err
	}
	return env.Data, nil
}

// CreateRecallDecisionSubmission creates a submission for a recall decision.
func (c *Client) CreateRecallDecisionSubmission(paymentID, recallID, decisionID string, sub RecallDecisionSubmission) (*RecallDecisionSubmission, error) {
	var env DataEnvelope[RecallDecisionSubmission]
	err := c.doJSON("POST",
		fmt.Sprintf("%s/%s/recalls/%s/decisions/%s/submissions", paymentsPath, paymentID, recallID, decisionID),
		DataEnvelope[RecallDecisionSubmission]{Data: sub}, &env)
	if err != nil {
		return nil, err
	}
	return &env.Data, nil
}

// GetRecallDecisionSubmission retrieves a specific recall decision submission.
func (c *Client) GetRecallDecisionSubmission(paymentID, recallID, decisionID, submissionID string) (*RecallDecisionSubmission, error) {
	var env DataEnvelope[RecallDecisionSubmission]
	err := c.doJSON("GET",
		fmt.Sprintf("%s/%s/recalls/%s/decisions/%s/submissions/%s", paymentsPath, paymentID, recallID, decisionID, submissionID),
		nil, &env)
	if err != nil {
		return nil, err
	}
	return &env.Data, nil
}
