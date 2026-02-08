package models

// PaymentAdmission represents a payment admission resource.
type PaymentAdmission struct {
	Resource
	Attributes PaymentAdmissionAttributes `json:"attributes"`
}

// PaymentAdmissionAttributes holds admission data.
type PaymentAdmissionAttributes struct {
	Status        string `json:"status"`
	AdmissionDate string `json:"admission_date,omitempty"`
}

// AdmissionTask represents a task on an admission.
type AdmissionTask struct {
	Resource
	Attributes AdmissionTaskAttributes `json:"attributes"`
}

// AdmissionTaskAttributes holds task data.
type AdmissionTaskAttributes struct {
	Status   string `json:"status"`
	Assignee string `json:"assignee,omitempty"`
	Name     string `json:"name,omitempty"`
}
