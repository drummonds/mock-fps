package client

import "fmt"

// CreatePaymentAdmission creates an admission for a payment.
func (c *Client) CreatePaymentAdmission(paymentID string, adm PaymentAdmission) (*PaymentAdmission, error) {
	var env DataEnvelope[PaymentAdmission]
	err := c.doJSON("POST",
		fmt.Sprintf("%s/%s/admissions", paymentsPath, paymentID),
		DataEnvelope[PaymentAdmission]{Data: adm}, &env)
	if err != nil {
		return nil, err
	}
	return &env.Data, nil
}

// GetPaymentAdmission retrieves a specific admission for a payment.
func (c *Client) GetPaymentAdmission(paymentID, admissionID string) (*PaymentAdmission, error) {
	var env DataEnvelope[PaymentAdmission]
	err := c.doJSON("GET",
		fmt.Sprintf("%s/%s/admissions/%s", paymentsPath, paymentID, admissionID),
		nil, &env)
	if err != nil {
		return nil, err
	}
	return &env.Data, nil
}

// UpdateAdmissionTask patches a task on an admission.
func (c *Client) UpdateAdmissionTask(paymentID, admissionID, taskID string, task AdmissionTask) (*AdmissionTask, error) {
	var env DataEnvelope[AdmissionTask]
	err := c.doJSON("PATCH",
		fmt.Sprintf("%s/%s/admissions/%s/tasks/%s", paymentsPath, paymentID, admissionID, taskID),
		DataEnvelope[AdmissionTask]{Data: task}, &env)
	if err != nil {
		return nil, err
	}
	return &env.Data, nil
}
