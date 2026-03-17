# FPS Inbound Payment Events — Complete Reference

This document covers every event that can occur for a single inbound FPS payment as seen through Form3's API, from initial receipt through settlement and exception handling.

## Payment Types

FPS supports four payment types:

| Type | Code | Delivery | Notes |
|---|---|---|---|
| Single Immediate Payment | SIP | Synchronous, <15s round-trip | Most common |
| Standing Order Payment | SOP | Asynchronous | Recurring |
| Forward Dated Payment | FDP | Asynchronous | Future-dated |
| Direct Corporate Access | DCA | Batch via Secure-IP | File-based |

This document focuses primarily on SIP flows, which are the most relevant for real-time testing.

## Inbound Payment Lifecycle

### Normal Flow (Happy Path)

```
VocaLink sends ISO 8583 message
  -> Form3 receives, creates Payment + PaymentAdmission (status: pending)
  -> Webhook: payment_admissions.created
  -> Form3 validates, updates admission (status: confirmed)
  -> Webhook: payment_admissions.updated
  -> Participant credits beneficiary account
  -> Settlement at next cycle (07:15, 13:00, or 15:45)
```

### Form3 Resource Model

```
Payment (inbound)
  +-- PaymentAdmission
  |     status: pending -> confirmed | failed
  |     status_reason: accepted | (44 failure reasons — see below)
  |     posting_status: pending -> posted | rejected | blocked | ...
  |     +-- PaymentAdmissionTask (optional, customer action required)
  |           status: pending -> completed | failed | on_hold
  |
  +-- Return (participant-initiated)
  |     +-- ReturnSubmission
  |     |     status: accepted -> ... -> delivery_confirmed | delivery_failed
  |     +-- ReturnAdmission (inbound return on outbound payment)
  |           status: pending -> confirmed | failed
  |
  +-- Recall (inbound recall request from sending bank)
  |     +-- RecallSubmission
  |     |     status: accepted -> delivery_confirmed
  |     +-- RecallDecision (participant's response)
  |     |     answer: accepted | rejected | pending | partially_accepted | payment_cancelled
  |     |     +-- RecallDecisionSubmission
  |     |     |     status: accepted -> ... -> delivery_confirmed | delivery_failed
  |     |     +-- RecallDecisionAdmission
  |     |           status: confirmed | failed
  |     +-- RecallReversal (scheme-generated if recall decision fails)
  |           +-- RecallReversalAdmission
  |
  +-- Reversal
        +-- ReversalSubmission
        |     status: accepted -> delivery_confirmed
        +-- ReversalAdmission
        |     status: pending -> confirmed
        +-- ReversalAdmissionTask (optional)
              status: pending -> completed | failed | on_hold
```

## All Possible Payment Outcomes

### 1. Successful Acceptance

**Unqualified Accept** — Payment accepted, beneficiary credited within 2 hours (typically seconds).

```
PaymentAdmission: pending -> confirmed
```

### 2. Qualified Accept (Stand-in)

Payment accepted, but credit may be delayed beyond the normal SLA. The receiving bank's processor is unavailable; the gateway accepts on their behalf.

```
PaymentAdmission: pending -> confirmed (with acceptance_qualifier)
```

Form3 `acceptance_qualifier` values:
- `none` — standard accept
- `same_day` — credit by end of day
- `next_calendar_day`
- `next_working_day`
- `after_next_working_day`
- `some_other_time`

- Funds stored-and-forwarded to participant when processor recovers
- Beneficiary must still be credited (no rejection)
- The 2-hour SLA for crediting may be exceeded

### 3. Rejection

Payment rejected synchronously by the receiving bank. The admission fails.

```
PaymentAdmission: pending -> failed
```

Form3 `PaymentAdmissionStatusReason` values (from OpenAPI spec):

**Account-related:**
- `account_closed`, `account_closed_beneficiary_deceased`, `account_closed_beneficiary_sensitivities`, `account_closed_business_reasons`, `account_closed_currency`, `account_closed_stopped`, `account_closed_terms_and_conditions`, `account_closed_transferred`
- `blocked_account`, `unknown_accountnumber`

**Agent/infrastructure errors:**
- `agent_clearing_process_error`, `agent_clearing_process_timeout`, `agent_reason_unknown`, `agent_suspended`, `agent_unavailable`
- `beneficiary_agent_clearing_process_error`, `beneficiary_agent_clearing_process_timeout`, `beneficiary_agent_suspended`, `beneficiary_agent_unavailable`

**Validation:**
- `amount_exceeds_settlement_limit`, `amount_invalid_or_missing`, `amount_not_allowed`
- `bankid_not_provisioned`, `invalid_bank_ID`, `invalid_bank_operation_code`
- `invalid_beneficiary_address`, `invalid_beneficiary_agent_BIC`, `invalid_beneficiary_details`
- `invalid_debtor_agent_BIC`, `invalid_debtor_details`
- `end_to_end_id_missing_or_invalid`
- `incorrect_reference_reference_mask`, `incorrect_reference_secondary_identification`, `incorrect_reference_validation_type`
- `beneficiary_name_not_present`

**Other:**
- `accepted` (success)
- `business_reasons`, `customer_check_failed`, `customer_reason_unknown`
- `duplicate_payment`, `original_payment_not_received`
- `regulatory_reason`, `scheme_timeout`
- `transaction_forbidden`, `transaction_type_not_supported`
- `rejected_by_customer`

`PostingStatus` values: `pending`, `posted`, `rejected`, `blocked`, `passed_without_posting`, `received`, `accepted_funds_checked`, `accepted_settlement_in_process`, `accepted_settlement_complete`, `accepted_settlement_complete_creditor_account`

`AccountValidationOutcome`: `passed`, `failed`, `failed_auto_reject_disabled`, `failed_auto_reject_enabled`

### 4. Return

A return is a **new payment** sent back after initial acceptance. The participant creates a Return resource.

```
Return created
  -> ReturnSubmission: accepted -> delivery_confirmed
```

Return reasons (codes use `RET` prefix + 4-8 digit code):
- Beneficiary account cannot accept funds (discovered post-acceptance)
- Failed AML/sanctions checks
- Fraudulent payment identified
- Account validation failures found after acceptance

**Timing**: 1-3 working days after the original payment.

### 5. Recall (Inbound)

The sending bank requests return of funds. FPS does not have automated recall like SEPA — this is a request, not a mandate. Form3 creates the Recall resource; the participant responds with a RecallDecision.

```
Recall created (inbound from sending bank)
  -> Webhook: recalls.created
  -> Participant creates RecallDecision
  -> RecallDecisionSubmission: accepted -> ... -> delivery_confirmed
```

RecallDecision `answer` values:
- `accepted` — agree to return funds
- `rejected` — refuse (with reason)
- `pending` — still considering
- `partially_accepted` — return partial amount
- `payment_cancelled` — payment was already cancelled

- **15-day deadline** for recall processing
- If a positive recall decision cannot be processed by the scheme, a **RecallReversal** is triggered

### 6. Reversal

Undoes a payment that was accepted. Typically occurs in timeout scenarios or when stand-in payments cannot subsequently be applied.

```
Reversal created
  -> ReversalSubmission: accepted -> delivery_confirmed
```

### 7. Timeout

If VocaLink receives no response from the beneficiary bank within ~15-30 seconds:
- The payment is automatically reversed/returned to the originator
- The receiving participant may see a reversal admission

**No partial acceptance exists in FPS for individual payments.** Each payment is all-or-nothing. However, recall decisions support `partially_accepted` (return a partial amount).

## Status Chains (Form3 API)

### Payment Submission (Outbound)

```
accepted
  -> validation_pending
  -> limit_check_pending
  -> limit_check_passed | limit_check_failed
  -> released_to_gateway
  -> queued_for_delivery
  -> submitted
  -> delivery_confirmed | delivery_failed
```

Submissions also carry `limit_breach_start_datetime` and `limit_breach_end_datetime` when held due to MNSC limits.

### Payment Admission (Inbound)

```
pending -> confirmed (status_reason: accepted)
pending -> failed (status_reason: one of 44 reasons above)
```

### Return Submission (Outbound)

```
accepted -> validation_pending -> limit_check_pending -> limit_check_passed
  -> released_to_gateway -> queued_for_delivery -> delivery_confirmed | delivery_failed
```

### Recall Decision Submission

```
accepted -> validation_pending -> validation_passed -> limit_check_pending
  -> limit_check_passed -> released_to_gateway -> queued_for_delivery
  -> delivery_confirmed | delivery_failed
```

### Recall Decision Admission

```
confirmed | failed
```

### Reversal Admission

```
pending -> confirmed
```

## Webhook Event Types

Subscriptions filter by `(record_type, event_type)` pair.

| record_type | event_type | When |
|---|---|---|
| `payments` | `created` | Payment resource created |
| `payment_admissions` | `created` | Inbound payment received |
| `payment_admissions` | `updated` | Admission status changed |
| `payment_admission_tasks` | `created` | Admission task requires customer action |
| `payment_admission_tasks` | `updated` | Admission task status changed |
| `payment_submissions` | `created` | Outbound submission created |
| `payment_submissions` | `updated` | Submission status changed |
| `return_payments` | `created` | Return created |
| `return_submissions` | `updated` | Return submission status changed |
| `recalls` | `created` | Recall request received |
| `recall_submissions` | `updated` | Recall submission status changed |
| `recall_decisions` | `created` | Recall decision created |
| `recall_decision_submissions` | `updated` | Decision submission status changed |
| `reversals` | `created` | Reversal created |
| `reversal_admissions` | `created` | Reversal received |
| `reversal_submissions` | `updated` | Reversal submission status changed |
| `reports` | `created` | Settlement report available |

## FPS Settlement Cycles

FPS uses **deferred net settlement (DNS)**. Payments clear in real-time but interbank settlement is batched.

**3 cycles per business day** via Bank of England RTGS:

| Cycle | Approximate Time |
|---|---|
| 1 | 07:15 |
| 2 | 13:00 |
| 3 | 15:45 |

RTGS operates 06:00-18:00 on business days. Outside these hours (evenings, weekends, bank holidays), payments continue processing but settlement is deferred to the next business day.

### Risk Controls

- **Multilateral Net Sender Cap (MNSC)**: Limits how much a participant can be a net sender between settlement cycles
- **Reserves Collateralisation Account (RCA)**: Pre-funded cash at the BoE (minimum GBP 50m) equal to the MNSC
- If a participant defaults, Pay.UK instructs the BoE to use RCA cash to complete settlement

### Reconciliation

Each settlement cycle produces reconciliation data:
- VocaLink calculates **multilateral net positions** per participant
- Settlement message sent to RTGS at the Bank of England
- Participants reconcile their transaction records against settlement positions

**Form3 reconciliation support:**

| Endpoint | Purpose |
|---|---|
| `GET /v1/transaction/payments` with settlement filters | Query by settlement_cycle, settlement_date |
| `GET /v1/notification/reports` | Scheme-generated settlement reports |
| `GET /v1/notification/reports/{id}/content` | Download report content |
| `GET /v1/organisation/positions` | Net settlement positions |
| `GET /v1/audit/entries/{record_type}` | Change history |

Key filters for reconciliation:
- `filter[submission.settlement_cycle]` — integer
- `filter[submission.settlement_date]` — date
- `filter[admission.settlement_cycle]` / `filter[admission.settlement_date]`
- `filter[processing_date_from]` / `filter[processing_date_to]`

Form3 creates a dummy inbound report every business day in the FPS Direct simulator.

## Stand-in Mode

### What Is Stand-in?

Stand-in is a resilience mechanism where a gateway accepts inbound payments on behalf of a participant whose payment processor is unavailable. Payments are stored and forwarded when the processor recovers.

### When It Happens

- Planned maintenance windows
- Unexpected network disconnection between gateway and processor
- Cloud/platform outages

### Key Rules

- Payments are **never rejected** during stand-in — the gateway accepts them
- The FPS scheme SLA requires beneficiary credit within **2 hours** for unqualified accepts
- Stand-in responses are near-instantaneous (only a timeout triggers the switch)
- When the processor comes back, stored payments are forwarded and applied

### Form3's Handling

Form3 handles stand-in internally:
- Not exposed in the payments API
- Visible only via Prometheus metrics endpoints (queue depths, processing times)
- Queues payments when VocaLink is unavailable, delivers on recovery

### mock-fps Stand-in

The mock server has a custom `/admin/standin` endpoint for testing:
- `GET /admin/standin` — returns `{enabled, queue_length}`
- `PUT /admin/standin` — toggles stand-in on/off
- When enabled, lifecycle transitions are queued, not executed
- When disabled, queued transitions drain and start

### Stand-in Duration Limits

The exact maximum duration before Pay.UK considers a participant in breach is not publicly documented. The key constraint is the 2-hour credit SLA for standard (unqualified) accepts. Extended stand-in periods risk breaching scheme rules, but the specific threshold is in participant-only Pay.UK documentation.

## Complete Event Sequence — Inbound SIP

For a bank receiving a single inbound FPS payment via Form3, the complete set of possible events:

### Phase 1: Receipt (milliseconds)
1. VocaLink routes ISO 8583 message to Form3
2. Form3 creates `Payment` resource
3. Form3 creates `PaymentAdmission` (status: `pending`)
4. Webhook: `payment_admissions.created`

### Phase 2: Validation (milliseconds to minutes)
5. Form3 validates payment
6. **Optional:** `PaymentAdmissionTask` created if customer action needed
   - Webhook: `payment_admission_tasks.created`
   - Participant patches task to `completed` or `failed`
   - Webhook: `payment_admission_tasks.updated`
7. **Either:**
   - Admission status -> `confirmed` (status_reason: `accepted`)
   - Admission status -> `failed` (status_reason: one of 44 reasons) — **END**
8. Webhook: `payment_admissions.updated`

### Phase 3: Participant Processing (seconds to hours)
8. Participant's system receives webhook, credits beneficiary
9. **If unable to credit:**
   - Create `Return` + `ReturnSubmission` (within 1-3 days) — go to Phase 5

### Phase 4: Settlement (next cycle)
10. Payment included in next settlement cycle
11. Net positions calculated, settled via RTGS

### Phase 5: Exception Handling (days)
Any of these may occur independently after acceptance:

**Return flow:**
12. Participant creates `Return` resource
13. Participant creates `ReturnSubmission` (status: `accepted`)
14. Webhook: `return_submissions.updated` through status chain to `delivery_confirmed`

**Recall flow (inbound):**
15. Form3 creates `Recall` resource (from sending bank)
16. Webhook: `recalls.created`
17. Participant creates `RecallDecision` (accepted or rejected)
18. Participant creates `RecallDecisionSubmission`
19. Webhook: `recall_decision_submissions.updated` -> `delivery_confirmed`
20. **If positive decision fails:** `RecallReversal` + `RecallReversalAdmission` created

**Reversal flow:**
21. `Reversal` created (timeout or stand-in failure)
22. `ReversalSubmission` status chain -> `delivery_confirmed`

## Sources

- [Form3 API Docs — FPS Direct](https://www.api-docs.form3.tech/api/schemes/fps-direct/payments/overview)
- [Form3 Reports Tutorial](https://api-docs.form3.tech/tutorial-reports.html)
- [Pay.UK FPS Principles v11 (Jan 2026)](https://www.wearepay.uk/wp-content/uploads/2026/02/Pay.UK-Faster-Payments-System-Principles-v-11-Jan-2026.pdf)
- [Bank of England RTGS Daily Timetable](https://www.bankofengland.co.uk/payment-and-settlement/summary-of-rtgs-daily-timetable)
- [ClearBank — Faster Payments](https://clearbank.github.io/uk/docs/gbp-payments/faster-payments/)
- [Monzo — Processing Payments in Stand-in](https://monzo.com/blog/processing-payments-in-monzo-stand-in)
- [Form3 OpenAPI Spec](https://www.api-docs.form3.tech/assets/swagger/form3-swagger.yaml) (158 paths, authoritative for status enums)
- [go-form3 Client](https://github.com/form3tech-oss/go-form3) (OpenAPI-generated models)
