package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
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

	resp, _ := http.Get(srv.URL + "/v1/transaction/payments/nonexistent")
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
	resp2, _ := http.Get(srv.URL + "/v1/transaction/payments/p1/submissions/s1")
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
	resp, _ := http.Post(srv.URL+"/v1/transaction/payments/nonexistent/submissions", jsonapi.ContentType, bytes.NewReader(body))
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
	resp, _ := http.Post(srv.URL+"/v1/transaction/payments/p1/returns", jsonapi.ContentType, bytes.NewReader(body))
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}

	// Get return
	resp2, _ := http.Get(srv.URL + "/v1/transaction/payments/p1/returns/r1")
	defer resp2.Body.Close()
	if resp2.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp2.StatusCode)
	}

	// List returns
	resp3, _ := http.Get(srv.URL + "/v1/transaction/payments/p1/returns")
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
	resp, _ := http.Post(srv.URL+"/v1/notification/subscriptions", jsonapi.ContentType, bytes.NewReader(body))
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}

	// Get
	resp2, _ := http.Get(srv.URL + "/v1/notification/subscriptions/sub1")
	defer resp2.Body.Close()
	if resp2.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp2.StatusCode)
	}

	// List
	resp3, _ := http.Get(srv.URL + "/v1/notification/subscriptions")
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
	resp4, _ := http.DefaultClient.Do(req)
	defer resp4.Body.Close()
	if resp4.StatusCode != http.StatusOK {
		t.Errorf("expected 200 for PATCH, got %d", resp4.StatusCode)
	}

	// Delete
	delReq, _ := http.NewRequest(http.MethodDelete, srv.URL+"/v1/notification/subscriptions/sub1", nil)
	resp5, _ := http.DefaultClient.Do(delReq)
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
	resp, _ := http.Post(srv.URL+"/v1/transaction/payments/p1/recalls", jsonapi.ContentType, bytes.NewReader(body))
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("create recall: expected 201, got %d", resp.StatusCode)
	}

	decision := models.RecallDecision{
		Resource:   models.Resource{ID: "dec1"},
		Attributes: models.RecallDecisionAttributes{Answer: "accepted"},
	}
	body, _ = json.Marshal(jsonapi.DataEnvelope[models.RecallDecision]{Data: decision})
	resp, _ = http.Post(srv.URL+"/v1/transaction/payments/p1/recalls/rec1/decisions", jsonapi.ContentType, bytes.NewReader(body))
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("create decision: expected 201, got %d", resp.StatusCode)
	}

	decSub := models.RecallDecisionSubmission{Resource: models.Resource{ID: "ds1"}}
	body, _ = json.Marshal(jsonapi.DataEnvelope[models.RecallDecisionSubmission]{Data: decSub})
	resp, _ = http.Post(srv.URL+"/v1/transaction/payments/p1/recalls/rec1/decisions/dec1/submissions", jsonapi.ContentType, bytes.NewReader(body))
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("create decision submission: expected 201, got %d", resp.StatusCode)
	}

	// Wait for lifecycle
	time.Sleep(100 * time.Millisecond)

	resp2, _ := http.Get(srv.URL + "/v1/transaction/payments/p1/recalls/rec1/decisions/dec1/submissions/ds1")
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
	resp, _ := http.Post(srv.URL+"/v1/transaction/payments/p1/reversals", jsonapi.ContentType, bytes.NewReader(body))
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("create reversal: expected 201, got %d", resp.StatusCode)
	}

	revSub := models.ReversalSubmission{Resource: models.Resource{ID: "rs1"}}
	body, _ = json.Marshal(jsonapi.DataEnvelope[models.ReversalSubmission]{Data: revSub})
	resp, _ = http.Post(srv.URL+"/v1/transaction/payments/p1/reversals/rev1/submissions", jsonapi.ContentType, bytes.NewReader(body))
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("create reversal submission: expected 201, got %d", resp.StatusCode)
	}

	time.Sleep(100 * time.Millisecond)

	resp2, _ := http.Get(srv.URL + "/v1/transaction/payments/p1/reversals/rev1/submissions/rs1")
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
	resp, _ := http.Get(srv.URL + "/v1/transaction/payments/p1")
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
