package lifecycle

// StatusChain defines the sequence of statuses a resource transitions through.
type StatusChain []string

// PaymentSubmissionChain is the status lifecycle for payment submissions.
var PaymentSubmissionChain = StatusChain{
	"accepted",
	"validation_pending",
	"limit_check_pending",
	"limit_check_passed",
	"released_to_gateway",
	"queued_for_delivery",
	"submitted",
	"delivery_confirmed",
}

// AdmissionChain is the status lifecycle for payment admissions.
var AdmissionChain = StatusChain{
	"pending",
	"confirmed",
}

// SimpleSubmissionChain is the status lifecycle for return/recall/reversal submissions.
var SimpleSubmissionChain = StatusChain{
	"accepted",
	"delivery_confirmed",
}
