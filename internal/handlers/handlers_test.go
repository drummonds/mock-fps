package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/nibble/mock-fps/internal/handlers"
	"github.com/nibble/mock-fps/internal/jsonapi"
	"github.com/nibble/mock-fps/internal/lifecycle"
	"github.com/nibble/mock-fps/internal/models"
	"github.com/nibble/mock-fps/internal/store"
)

func setupServer() *httptest.Server {
	s := store.NewMemoryStore()
	engine := lifecycle.NewEngine(10, nil) // 10ms steps for fast tests
	mux := http.NewServeMux()
	handlers.RegisterRoutes(mux, s, engine)
	return httptest.NewServer(mux)
}

func TestHealthCheck(t *testing.T) {
	srv := setupServer()
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/health")
	if err != nil {
		t.Fatalf("GET /health: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestCreateAndGetPayment(t *testing.T) {
	srv := setupServer()
	defer srv.Close()

	payment := models.Payment{
		Resource: models.Resource{
			ID:             "test-payment-1",
			OrganisationID: "org-1",
		},
		Attributes: models.PaymentAttributes{
			Amount:   "100.50",
			Currency: "GBP",
		},
	}
	body, _ := json.Marshal(jsonapi.DataEnvelope[models.Payment]{Data: payment})

	// Create
	resp, err := http.Post(srv.URL+"/v1/transaction/payments", jsonapi.ContentType, bytes.NewReader(body))
	if err != nil {
		t.Fatalf("POST payment: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}

	var created jsonapi.DataEnvelope[models.Payment]
	json.NewDecoder(resp.Body).Decode(&created)
	if created.Data.ID != "test-payment-1" {
		t.Errorf("expected id test-payment-1, got %s", created.Data.ID)
	}

	// Get
	resp2, err := http.Get(srv.URL + "/v1/transaction/payments/test-payment-1")
	if err != nil {
		t.Fatalf("GET payment: %v", err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp2.StatusCode)
	}
}

func TestListPayments(t *testing.T) {
	srv := setupServer()
	defer srv.Close()

	// Create two payments
	for _, id := range []string{"p1", "p2"} {
		p := models.Payment{
			Resource:   models.Resource{ID: id},
			Attributes: models.PaymentAttributes{Amount: "50.00", Currency: "GBP"},
		}
		body, _ := json.Marshal(jsonapi.DataEnvelope[models.Payment]{Data: p})
		http.Post(srv.URL+"/v1/transaction/payments", jsonapi.ContentType, bytes.NewReader(body))
	}

	resp, err := http.Get(srv.URL + "/v1/transaction/payments")
	if err != nil {
		t.Fatalf("GET payments: %v", err)
	}
	defer resp.Body.Close()

	var list jsonapi.ListEnvelope[models.Payment]
	json.NewDecoder(resp.Body).Decode(&list)
	if len(list.Data) != 2 {
		t.Errorf("expected 2 payments, got %d", len(list.Data))
	}
}

func TestPaymentNotFound(t *testing.T) {
	srv := setupServer()
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/v1/transaction/payments/nonexistent")
	if err != nil {
		t.Fatalf("GET nonexistent: %v", err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404, got %d", resp.StatusCode)
	}
}

func TestPaymentSubmissionLifecycle(t *testing.T) {
	srv := setupServer()
	defer srv.Close()

	// Create payment
	payment := models.Payment{
		Resource:   models.Resource{ID: "p1"},
		Attributes: models.PaymentAttributes{Amount: "100.00", Currency: "GBP"},
	}
	body, _ := json.Marshal(jsonapi.DataEnvelope[models.Payment]{Data: payment})
	http.Post(srv.URL+"/v1/transaction/payments", jsonapi.ContentType, bytes.NewReader(body))

	// Create submission
	sub := models.PaymentSubmission{Resource: models.Resource{ID: "s1"}}
	body, _ = json.Marshal(jsonapi.DataEnvelope[models.PaymentSubmission]{Data: sub})
	resp, err := http.Post(srv.URL+"/v1/transaction/payments/p1/submissions", jsonapi.ContentType, bytes.NewReader(body))
	if err != nil {
		t.Fatalf("POST submission: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}

	var created jsonapi.DataEnvelope[models.PaymentSubmission]
	json.NewDecoder(resp.Body).Decode(&created)
	if created.Data.Attributes.Status != "accepted" {
		t.Errorf("expected initial status accepted, got %s", created.Data.Attributes.Status)
	}

	// Wait for lifecycle to complete (10ms * 7 steps + buffer)
	time.Sleep(200 * time.Millisecond)

	// Check final status
	resp2, err := http.Get(srv.URL + "/v1/transaction/payments/p1/submissions/s1")
	if err != nil {
		t.Fatalf("GET submission: %v", err)
	}
	defer resp2.Body.Close()

	var final jsonapi.DataEnvelope[models.PaymentSubmission]
	json.NewDecoder(resp2.Body).Decode(&final)
	if final.Data.Attributes.Status != "delivery_confirmed" {
		t.Errorf("expected final status delivery_confirmed, got %s", final.Data.Attributes.Status)
	}
}

func TestSubmissionRequiresPayment(t *testing.T) {
	srv := setupServer()
	defer srv.Close()

	sub := models.PaymentSubmission{Resource: models.Resource{ID: "s1"}}
	body, _ := json.Marshal(jsonapi.DataEnvelope[models.PaymentSubmission]{Data: sub})
	resp, err := http.Post(srv.URL+"/v1/transaction/payments/nonexistent/submissions", jsonapi.ContentType, bytes.NewReader(body))
	if err != nil {
		t.Fatalf("POST submission: %v", err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404, got %d", resp.StatusCode)
	}
}

func TestReturnFlow(t *testing.T) {
	srv := setupServer()
	defer srv.Close()

	// Create payment
	payment := models.Payment{
		Resource:   models.Resource{ID: "p1"},
		Attributes: models.PaymentAttributes{Amount: "100.00", Currency: "GBP"},
	}
	body, _ := json.Marshal(jsonapi.DataEnvelope[models.Payment]{Data: payment})
	http.Post(srv.URL+"/v1/transaction/payments", jsonapi.ContentType, bytes.NewReader(body))

	// Create return
	ret := models.ReturnPayment{
		Resource:   models.Resource{ID: "r1"},
		Attributes: models.ReturnPaymentAttributes{Amount: "50.00", Currency: "GBP", ReturnCode: "1100"},
	}
	body, _ = json.Marshal(jsonapi.DataEnvelope[models.ReturnPayment]{Data: ret})
	resp, err := http.Post(srv.URL+"/v1/transaction/payments/p1/returns", jsonapi.ContentType, bytes.NewReader(body))
	if err != nil {
		t.Fatalf("POST return: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}

	// Get return
	resp2, err := http.Get(srv.URL + "/v1/transaction/payments/p1/returns/r1")
	if err != nil {
		t.Fatalf("GET return: %v", err)
	}
	defer resp2.Body.Close()
	if resp2.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp2.StatusCode)
	}

	// List returns
	resp3, err := http.Get(srv.URL + "/v1/transaction/payments/p1/returns")
	if err != nil {
		t.Fatalf("GET returns: %v", err)
	}
	defer resp3.Body.Close()
	var list jsonapi.ListEnvelope[models.ReturnPayment]
	json.NewDecoder(resp3.Body).Decode(&list)
	if len(list.Data) != 1 {
		t.Errorf("expected 1 return, got %d", len(list.Data))
	}
}

func TestSubscriptionCRUD(t *testing.T) {
	srv := setupServer()
	defer srv.Close()

	sub := models.Subscription{
		Resource: models.Resource{ID: "sub1"},
		Attributes: models.SubscriptionAttributes{
			CallbackURI: "http://example.com/webhook",
			EventType:   "updated",
			RecordType:  "payment_submissions",
		},
	}
	body, _ := json.Marshal(jsonapi.DataEnvelope[models.Subscription]{Data: sub})

	// Create
	resp, err := http.Post(srv.URL+"/v1/notification/subscriptions", jsonapi.ContentType, bytes.NewReader(body))
	if err != nil {
		t.Fatalf("POST subscription: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}

	// Get
	resp2, err := http.Get(srv.URL + "/v1/notification/subscriptions/sub1")
	if err != nil {
		t.Fatalf("GET subscription: %v", err)
	}
	defer resp2.Body.Close()
	if resp2.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp2.StatusCode)
	}

	// List
	resp3, err := http.Get(srv.URL + "/v1/notification/subscriptions")
	if err != nil {
		t.Fatalf("GET subscriptions: %v", err)
	}
	defer resp3.Body.Close()
	var list jsonapi.ListEnvelope[models.Subscription]
	json.NewDecoder(resp3.Body).Decode(&list)
	if len(list.Data) != 1 {
		t.Errorf("expected 1 subscription, got %d", len(list.Data))
	}

	// Patch
	patch := models.Subscription{
		Attributes: models.SubscriptionAttributes{
			CallbackURI: "http://example.com/webhook2",
			IsActive:    false,
		},
	}
	patchBody, _ := json.Marshal(jsonapi.DataEnvelope[models.Subscription]{Data: patch})
	req, _ := http.NewRequest(http.MethodPatch, srv.URL+"/v1/notification/subscriptions/sub1", bytes.NewReader(patchBody))
	req.Header.Set("Content-Type", jsonapi.ContentType)
	resp4, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("PATCH subscription: %v", err)
	}
	defer resp4.Body.Close()
	if resp4.StatusCode != http.StatusOK {
		t.Errorf("expected 200 for PATCH, got %d", resp4.StatusCode)
	}

	// Delete
	delReq, _ := http.NewRequest(http.MethodDelete, srv.URL+"/v1/notification/subscriptions/sub1", nil)
	resp5, err := http.DefaultClient.Do(delReq)
	if err != nil {
		t.Fatalf("DELETE subscription: %v", err)
	}
	if resp5.StatusCode != http.StatusNoContent {
		t.Errorf("expected 204 for DELETE, got %d", resp5.StatusCode)
	}
}

func TestRecallDecisionFlow(t *testing.T) {
	srv := setupServer()
	defer srv.Close()

	// Setup: payment -> recall -> decision -> decision submission
	payment := models.Payment{Resource: models.Resource{ID: "p1"}, Attributes: models.PaymentAttributes{Amount: "100.00", Currency: "GBP"}}
	body, _ := json.Marshal(jsonapi.DataEnvelope[models.Payment]{Data: payment})
	http.Post(srv.URL+"/v1/transaction/payments", jsonapi.ContentType, bytes.NewReader(body))

	recall := models.Recall{Resource: models.Resource{ID: "rec1"}}
	body, _ = json.Marshal(jsonapi.DataEnvelope[models.Recall]{Data: recall})
	resp, err := http.Post(srv.URL+"/v1/transaction/payments/p1/recalls", jsonapi.ContentType, bytes.NewReader(body))
	if err != nil {
		t.Fatalf("POST recall: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("create recall: expected 201, got %d", resp.StatusCode)
	}

	decision := models.RecallDecision{
		Resource:   models.Resource{ID: "dec1"},
		Attributes: models.RecallDecisionAttributes{Answer: "accepted"},
	}
	body, _ = json.Marshal(jsonapi.DataEnvelope[models.RecallDecision]{Data: decision})
	resp, err = http.Post(srv.URL+"/v1/transaction/payments/p1/recalls/rec1/decisions", jsonapi.ContentType, bytes.NewReader(body))
	if err != nil {
		t.Fatalf("POST decision: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("create decision: expected 201, got %d", resp.StatusCode)
	}

	decSub := models.RecallDecisionSubmission{Resource: models.Resource{ID: "ds1"}}
	body, _ = json.Marshal(jsonapi.DataEnvelope[models.RecallDecisionSubmission]{Data: decSub})
	resp, err = http.Post(srv.URL+"/v1/transaction/payments/p1/recalls/rec1/decisions/dec1/submissions", jsonapi.ContentType, bytes.NewReader(body))
	if err != nil {
		t.Fatalf("POST decision submission: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("create decision submission: expected 201, got %d", resp.StatusCode)
	}

	// Wait for lifecycle
	time.Sleep(100 * time.Millisecond)

	resp2, err := http.Get(srv.URL + "/v1/transaction/payments/p1/recalls/rec1/decisions/dec1/submissions/ds1")
	if err != nil {
		t.Fatalf("GET decision submission: %v", err)
	}
	defer resp2.Body.Close()
	var got jsonapi.DataEnvelope[models.RecallDecisionSubmission]
	json.NewDecoder(resp2.Body).Decode(&got)
	if got.Data.Attributes.Status != "delivery_confirmed" {
		t.Errorf("expected delivery_confirmed, got %s", got.Data.Attributes.Status)
	}
}

func TestReversalFlow(t *testing.T) {
	srv := setupServer()
	defer srv.Close()

	// Payment -> reversal -> reversal submission
	payment := models.Payment{Resource: models.Resource{ID: "p1"}, Attributes: models.PaymentAttributes{Amount: "100.00", Currency: "GBP"}}
	body, _ := json.Marshal(jsonapi.DataEnvelope[models.Payment]{Data: payment})
	http.Post(srv.URL+"/v1/transaction/payments", jsonapi.ContentType, bytes.NewReader(body))

	rev := models.Reversal{Resource: models.Resource{ID: "rev1"}, Attributes: models.ReversalAttributes{Amount: "100.00"}}
	body, _ = json.Marshal(jsonapi.DataEnvelope[models.Reversal]{Data: rev})
	resp, err := http.Post(srv.URL+"/v1/transaction/payments/p1/reversals", jsonapi.ContentType, bytes.NewReader(body))
	if err != nil {
		t.Fatalf("POST reversal: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("create reversal: expected 201, got %d", resp.StatusCode)
	}

	revSub := models.ReversalSubmission{Resource: models.Resource{ID: "rs1"}}
	body, _ = json.Marshal(jsonapi.DataEnvelope[models.ReversalSubmission]{Data: revSub})
	resp, err = http.Post(srv.URL+"/v1/transaction/payments/p1/reversals/rev1/submissions", jsonapi.ContentType, bytes.NewReader(body))
	if err != nil {
		t.Fatalf("POST reversal submission: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("create reversal submission: expected 201, got %d", resp.StatusCode)
	}

	time.Sleep(100 * time.Millisecond)

	resp2, err := http.Get(srv.URL + "/v1/transaction/payments/p1/reversals/rev1/submissions/rs1")
	if err != nil {
		t.Fatalf("GET reversal submission: %v", err)
	}
	defer resp2.Body.Close()
	var got jsonapi.DataEnvelope[models.ReversalSubmission]
	json.NewDecoder(resp2.Body).Decode(&got)
	if got.Data.Attributes.Status != "delivery_confirmed" {
		t.Errorf("expected delivery_confirmed, got %s", got.Data.Attributes.Status)
	}
}

func TestPaymentRelationships(t *testing.T) {
	srv := setupServer()
	defer srv.Close()

	// Create payment with submission
	payment := models.Payment{Resource: models.Resource{ID: "p1"}, Attributes: models.PaymentAttributes{Amount: "100.00", Currency: "GBP"}}
	body, _ := json.Marshal(jsonapi.DataEnvelope[models.Payment]{Data: payment})
	http.Post(srv.URL+"/v1/transaction/payments", jsonapi.ContentType, bytes.NewReader(body))

	sub := models.PaymentSubmission{Resource: models.Resource{ID: "s1"}}
	body, _ = json.Marshal(jsonapi.DataEnvelope[models.PaymentSubmission]{Data: sub})
	http.Post(srv.URL+"/v1/transaction/payments/p1/submissions", jsonapi.ContentType, bytes.NewReader(body))

	// Get payment - should have relationships
	resp, err := http.Get(srv.URL + "/v1/transaction/payments/p1")
	if err != nil {
		t.Fatalf("GET payment: %v", err)
	}
	defer resp.Body.Close()
	var got jsonapi.DataEnvelope[models.Payment]
	json.NewDecoder(resp.Body).Decode(&got)

	if got.Data.Relationships == nil {
		t.Fatal("expected relationships to be set")
	}
	if got.Data.Relationships.PaymentSubmissions == nil {
		t.Fatal("expected payment_submissions relationship")
	}
	if len(got.Data.Relationships.PaymentSubmissions.Data) != 1 {
		t.Errorf("expected 1 submission relationship, got %d", len(got.Data.Relationships.PaymentSubmissions.Data))
	}
}

func TestPaymentFPSDefaults(t *testing.T) {
	srv := setupServer()
	defer srv.Close()

	// Create payment with no FPS fields
	payment := models.Payment{
		Resource:   models.Resource{ID: "fps-test"},
		Attributes: models.PaymentAttributes{Amount: "25.00", Currency: "GBP"},
	}
	body, _ := json.Marshal(jsonapi.DataEnvelope[models.Payment]{Data: payment})
	resp, err := http.Post(srv.URL+"/v1/transaction/payments", jsonapi.ContentType, bytes.NewReader(body))
	if err != nil {
		t.Fatalf("POST payment: %v", err)
	}
	defer resp.Body.Close()

	var created jsonapi.DataEnvelope[models.Payment]
	json.NewDecoder(resp.Body).Decode(&created)
	attrs := created.Data.Attributes

	if attrs.PaymentScheme != "FPS" {
		t.Errorf("payment_scheme = %q, want FPS", attrs.PaymentScheme)
	}
	if attrs.SchemePaymentType != "ImmediatePayment" {
		t.Errorf("scheme_payment_type = %q, want ImmediatePayment", attrs.SchemePaymentType)
	}
	if matched, _ := regexp.MatchString(`^\d{6}$`, attrs.NumericReference); !matched {
		t.Errorf("numeric_reference = %q, want 6 digits", attrs.NumericReference)
	}
	if matched, _ := regexp.MatchString(`^FPS\d{14}$`, attrs.EndToEndReference); !matched {
		t.Errorf("end_to_end_reference = %q, want FPS+date+6digits", attrs.EndToEndReference)
	}
	if matched, _ := regexp.MatchString(`^\d{4}-\d{2}-\d{2}$`, attrs.ProcessingDate); !matched {
		t.Errorf("processing_date = %q, want YYYY-MM-DD", attrs.ProcessingDate)
	}
}

func TestPaymentFPSClientValuesPreserved(t *testing.T) {
	srv := setupServer()
	defer srv.Close()

	payment := models.Payment{
		Resource: models.Resource{ID: "fps-custom"},
		Attributes: models.PaymentAttributes{
			Amount:            "10.00",
			Currency:          "GBP",
			PaymentScheme:     "BACS",
			SchemePaymentType: "Priority",
			NumericReference:  "999999",
			EndToEndReference: "CUSTOM-REF-123",
			ProcessingDate:    "2025-01-15",
		},
	}
	body, _ := json.Marshal(jsonapi.DataEnvelope[models.Payment]{Data: payment})
	resp, err := http.Post(srv.URL+"/v1/transaction/payments", jsonapi.ContentType, bytes.NewReader(body))
	if err != nil {
		t.Fatalf("POST payment: %v", err)
	}
	defer resp.Body.Close()

	var created jsonapi.DataEnvelope[models.Payment]
	json.NewDecoder(resp.Body).Decode(&created)
	attrs := created.Data.Attributes

	if attrs.PaymentScheme != "BACS" {
		t.Errorf("payment_scheme = %q, want BACS", attrs.PaymentScheme)
	}
	if attrs.SchemePaymentType != "Priority" {
		t.Errorf("scheme_payment_type = %q, want Priority", attrs.SchemePaymentType)
	}
	if attrs.NumericReference != "999999" {
		t.Errorf("numeric_reference = %q, want 999999", attrs.NumericReference)
	}
	if attrs.EndToEndReference != "CUSTOM-REF-123" {
		t.Errorf("end_to_end_reference = %q, want CUSTOM-REF-123", attrs.EndToEndReference)
	}
	if attrs.ProcessingDate != "2025-01-15" {
		t.Errorf("processing_date = %q, want 2025-01-15", attrs.ProcessingDate)
	}
}

func TestStandInGetDefault(t *testing.T) {
	srv := setupServer()
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/admin/standin")
	if err != nil {
		t.Fatalf("GET standin: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var got struct {
		Enabled     bool `json:"enabled"`
		QueueLength int  `json:"queue_length"`
	}
	json.NewDecoder(resp.Body).Decode(&got)
	if got.Enabled {
		t.Error("expected standin disabled by default")
	}
	if got.QueueLength != 0 {
		t.Errorf("expected queue_length 0, got %d", got.QueueLength)
	}
}

func TestStandInToggle(t *testing.T) {
	srv := setupServer()
	defer srv.Close()

	// Enable standin
	body, _ := json.Marshal(map[string]bool{"enabled": true})
	req, _ := http.NewRequest(http.MethodPut, srv.URL+"/admin/standin", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("PUT standin: %v", err)
	}
	defer resp.Body.Close()

	var got struct {
		Enabled     bool `json:"enabled"`
		QueueLength int  `json:"queue_length"`
	}
	json.NewDecoder(resp.Body).Decode(&got)
	if !got.Enabled {
		t.Error("expected standin enabled after PUT")
	}
}

func TestStandInQueuesSubmissions(t *testing.T) {
	srv := setupServer()
	defer srv.Close()

	// Enable standin
	body, _ := json.Marshal(map[string]bool{"enabled": true})
	req, _ := http.NewRequest(http.MethodPut, srv.URL+"/admin/standin", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	http.DefaultClient.Do(req)

	// Create payment + submission
	payment := models.Payment{Resource: models.Resource{ID: "p1"}, Attributes: models.PaymentAttributes{Amount: "100.00", Currency: "GBP"}}
	body, _ = json.Marshal(jsonapi.DataEnvelope[models.Payment]{Data: payment})
	http.Post(srv.URL+"/v1/transaction/payments", jsonapi.ContentType, bytes.NewReader(body))

	sub := models.PaymentSubmission{Resource: models.Resource{ID: "s1"}}
	body, _ = json.Marshal(jsonapi.DataEnvelope[models.PaymentSubmission]{Data: sub})
	http.Post(srv.URL+"/v1/transaction/payments/p1/submissions", jsonapi.ContentType, bytes.NewReader(body))

	// Wait a bit — lifecycle should NOT progress
	time.Sleep(200 * time.Millisecond)

	// Check submission is still at initial status
	resp, err := http.Get(srv.URL + "/v1/transaction/payments/p1/submissions/s1")
	if err != nil {
		t.Fatalf("GET submission: %v", err)
	}
	defer resp.Body.Close()
	var got jsonapi.DataEnvelope[models.PaymentSubmission]
	json.NewDecoder(resp.Body).Decode(&got)
	if got.Data.Attributes.Status != "accepted" {
		t.Errorf("expected status accepted (queued), got %s", got.Data.Attributes.Status)
	}

	// Check queue length is 1
	resp2, err := http.Get(srv.URL + "/admin/standin")
	if err != nil {
		t.Fatalf("GET standin: %v", err)
	}
	defer resp2.Body.Close()
	var state struct {
		QueueLength int `json:"queue_length"`
	}
	json.NewDecoder(resp2.Body).Decode(&state)
	if state.QueueLength != 1 {
		t.Errorf("expected queue_length 1, got %d", state.QueueLength)
	}

	// Disable standin — transitions should drain
	body, _ = json.Marshal(map[string]bool{"enabled": false})
	req, _ = http.NewRequest(http.MethodPut, srv.URL+"/admin/standin", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	http.DefaultClient.Do(req)

	// Wait for lifecycle to complete
	time.Sleep(200 * time.Millisecond)

	resp3, err := http.Get(srv.URL + "/v1/transaction/payments/p1/submissions/s1")
	if err != nil {
		t.Fatalf("GET submission after drain: %v", err)
	}
	defer resp3.Body.Close()
	var final jsonapi.DataEnvelope[models.PaymentSubmission]
	json.NewDecoder(resp3.Body).Decode(&final)
	if final.Data.Attributes.Status != "delivery_confirmed" {
		t.Errorf("expected delivery_confirmed after drain, got %s", final.Data.Attributes.Status)
	}
}

func TestSubmissionSchemeTransactionID(t *testing.T) {
	srv := setupServer()
	defer srv.Close()

	// Create payment
	payment := models.Payment{
		Resource:   models.Resource{ID: "p-stid"},
		Attributes: models.PaymentAttributes{Amount: "50.00", Currency: "GBP"},
	}
	body, _ := json.Marshal(jsonapi.DataEnvelope[models.Payment]{Data: payment})
	http.Post(srv.URL+"/v1/transaction/payments", jsonapi.ContentType, bytes.NewReader(body))

	// Create submission — should get auto scheme_transaction_id
	sub := models.PaymentSubmission{Resource: models.Resource{ID: "s-stid"}}
	body, _ = json.Marshal(jsonapi.DataEnvelope[models.PaymentSubmission]{Data: sub})
	resp, err := http.Post(srv.URL+"/v1/transaction/payments/p-stid/submissions", jsonapi.ContentType, bytes.NewReader(body))
	if err != nil {
		t.Fatalf("POST submission: %v", err)
	}
	defer resp.Body.Close()

	var created jsonapi.DataEnvelope[models.PaymentSubmission]
	json.NewDecoder(resp.Body).Decode(&created)

	if matched, _ := regexp.MatchString(`^\d{26}$`, created.Data.Attributes.SchemeTransactionID); !matched {
		t.Errorf("scheme_transaction_id = %q, want 26 digits", created.Data.Attributes.SchemeTransactionID)
	}

	// Create admission — should also get auto scheme_transaction_id
	adm := models.PaymentAdmission{Resource: models.Resource{ID: "a-stid"}}
	body, _ = json.Marshal(jsonapi.DataEnvelope[models.PaymentAdmission]{Data: adm})
	resp2, err := http.Post(srv.URL+"/v1/transaction/payments/p-stid/admissions", jsonapi.ContentType, bytes.NewReader(body))
	if err != nil {
		t.Fatalf("POST admission: %v", err)
	}
	defer resp2.Body.Close()

	var admCreated jsonapi.DataEnvelope[models.PaymentAdmission]
	json.NewDecoder(resp2.Body).Decode(&admCreated)

	if matched, _ := regexp.MatchString(`^\d{26}$`, admCreated.Data.Attributes.SchemeTransactionID); !matched {
		t.Errorf("admission scheme_transaction_id = %q, want 26 digits", admCreated.Data.Attributes.SchemeTransactionID)
	}
}
