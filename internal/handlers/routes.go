package handlers

import (
	"net/http"

	"github.com/nibble/mock-fps/internal/lifecycle"
	"github.com/nibble/mock-fps/internal/store"
)

const basePath = "/v1/transaction/payments"
const subsPath = "/v1/notification/subscriptions"

// RegisterRoutes registers all API routes on the given mux.
func RegisterRoutes(mux *http.ServeMux, s store.Store, engine *lifecycle.Engine) {
	payments := NewPaymentHandler(s)
	submissions := NewPaymentSubmissionHandler(s, engine)
	admissions := NewPaymentAdmissionHandler(s, engine)
	returns := NewPaymentReturnHandler(s)
	returnSubs := NewReturnSubmissionHandler(s, engine)
	recalls := NewPaymentRecallHandler(s)
	recallSubs := NewRecallSubmissionHandler(s, engine)
	decisions := NewRecallDecisionHandler(s)
	decisionSubs := NewRecallDecisionSubmissionHandler(s, engine)
	reversals := NewPaymentReversalHandler(s)
	reversalSubs := NewReversalSubmissionHandler(s, engine)
	subscriptions := NewSubscriptionHandler(s)

	// Payments
	mux.HandleFunc("POST "+basePath, payments.Create)
	mux.HandleFunc("GET "+basePath, payments.List)
	mux.HandleFunc("GET "+basePath+"/{paymentID}", payments.Get)

	// Payment Submissions
	mux.HandleFunc("POST "+basePath+"/{paymentID}/submissions", submissions.Create)
	mux.HandleFunc("GET "+basePath+"/{paymentID}/submissions/{submissionID}", submissions.Get)

	// Payment Admissions
	mux.HandleFunc("POST "+basePath+"/{paymentID}/admissions", admissions.Create)
	mux.HandleFunc("GET "+basePath+"/{paymentID}/admissions/{admissionID}", admissions.Get)
	mux.HandleFunc("PATCH "+basePath+"/{paymentID}/admissions/{admissionID}/tasks/{taskID}", admissions.PatchTask)

	// Returns
	mux.HandleFunc("POST "+basePath+"/{paymentID}/returns", returns.Create)
	mux.HandleFunc("GET "+basePath+"/{paymentID}/returns", returns.List)
	mux.HandleFunc("GET "+basePath+"/{paymentID}/returns/{returnID}", returns.Get)

	// Return Submissions
	mux.HandleFunc("POST "+basePath+"/{paymentID}/returns/{returnID}/submissions", returnSubs.Create)
	mux.HandleFunc("GET "+basePath+"/{paymentID}/returns/{returnID}/submissions/{submissionID}", returnSubs.Get)

	// Recalls
	mux.HandleFunc("POST "+basePath+"/{paymentID}/recalls", recalls.Create)
	mux.HandleFunc("GET "+basePath+"/{paymentID}/recalls", recalls.List)
	mux.HandleFunc("GET "+basePath+"/{paymentID}/recalls/{recallID}", recalls.Get)

	// Recall Submissions
	mux.HandleFunc("POST "+basePath+"/{paymentID}/recalls/{recallID}/submissions", recallSubs.Create)
	mux.HandleFunc("GET "+basePath+"/{paymentID}/recalls/{recallID}/submissions/{submissionID}", recallSubs.Get)

	// Recall Decisions
	mux.HandleFunc("POST "+basePath+"/{paymentID}/recalls/{recallID}/decisions", decisions.Create)
	mux.HandleFunc("GET "+basePath+"/{paymentID}/recalls/{recallID}/decisions", decisions.List)
	mux.HandleFunc("GET "+basePath+"/{paymentID}/recalls/{recallID}/decisions/{decisionID}", decisions.Get)

	// Recall Decision Submissions
	mux.HandleFunc("POST "+basePath+"/{paymentID}/recalls/{recallID}/decisions/{decisionID}/submissions", decisionSubs.Create)
	mux.HandleFunc("GET "+basePath+"/{paymentID}/recalls/{recallID}/decisions/{decisionID}/submissions/{submissionID}", decisionSubs.Get)

	// Reversals
	mux.HandleFunc("POST "+basePath+"/{paymentID}/reversals", reversals.Create)
	mux.HandleFunc("GET "+basePath+"/{paymentID}/reversals", reversals.List)
	mux.HandleFunc("GET "+basePath+"/{paymentID}/reversals/{reversalID}", reversals.Get)

	// Reversal Submissions
	mux.HandleFunc("POST "+basePath+"/{paymentID}/reversals/{reversalID}/submissions", reversalSubs.Create)
	mux.HandleFunc("GET "+basePath+"/{paymentID}/reversals/{reversalID}/submissions/{submissionID}", reversalSubs.Get)

	// Subscriptions
	mux.HandleFunc("POST "+subsPath, subscriptions.Create)
	mux.HandleFunc("GET "+subsPath, subscriptions.List)
	mux.HandleFunc("GET "+subsPath+"/{subscriptionID}", subscriptions.Get)
	mux.HandleFunc("PATCH "+subsPath+"/{subscriptionID}", subscriptions.Patch)
	mux.HandleFunc("DELETE "+subsPath+"/{subscriptionID}", subscriptions.Delete)

	// Health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

}
