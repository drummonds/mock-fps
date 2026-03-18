# FPS Resource States

Quick reference for all resource states in the mock-fps server.

## Payment Submission (Outbound)

```
accepted → validation_pending → limit_check_pending → limit_check_passed
  → released_to_gateway → queued_for_delivery → submitted → delivery_confirmed
```

Failure: `delivery_failed` (from any intermediate state)

On reaching `delivery_confirmed`, `settlement_cycle` and `settlement_date` are assigned.

## Payment Admission (Inbound)

```
pending → confirmed
```

Failure: `failed` (with `status_reason`)

On reaching `confirmed`, `settlement_cycle` and `settlement_date` are assigned.

### Admission Status Reasons

When an admission fails, `status_reason` is one of:

| Category | Values |
|---|---|
| Success | `accepted` |
| Account | `account_closed`, `account_closed_beneficiary_deceased`, `account_closed_business_reasons`, `account_closed_currency`, `account_closed_stopped`, `account_closed_transferred`, `blocked_account`, `unknown_accountnumber` |
| Agent | `agent_clearing_process_error`, `agent_clearing_process_timeout`, `agent_reason_unknown`, `agent_suspended`, `agent_unavailable`, `beneficiary_agent_*` variants |
| Validation | `amount_exceeds_settlement_limit`, `amount_invalid_or_missing`, `bankid_not_provisioned`, `invalid_bank_ID`, `invalid_beneficiary_*`, `invalid_debtor_*`, `end_to_end_id_missing_or_invalid` |
| Other | `duplicate_payment`, `regulatory_reason`, `scheme_timeout`, `transaction_forbidden` |

### Acceptance Qualifiers (Stand-in)

| Value | Meaning |
|---|---|
| `none` | Standard accept |
| `same_day` | Credit by end of day |
| `next_calendar_day` | Credit by next calendar day |
| `next_working_day` | Credit by next working day |
| `after_next_working_day` | Credit after next working day |
| `some_other_time` | Credit at some other time |

## Return Submission

```
accepted → validation_pending → limit_check_pending → limit_check_passed
  → released_to_gateway → queued_for_delivery → delivery_confirmed
```

Failure: `delivery_failed`

## Recall Decision Submission

```
accepted → validation_pending → validation_passed → limit_check_pending
  → limit_check_passed → released_to_gateway → queued_for_delivery → delivery_confirmed
```

Failure: `delivery_failed`

### Recall Decision Answers

| Value | Meaning |
|---|---|
| `accepted` | Agree to return funds |
| `rejected` | Refuse to return funds |
| `pending` | Still considering |
| `partially_accepted` | Return partial amount |
| `payment_cancelled` | Payment was already cancelled |

## Recall Decision Admission

```
confirmed
```

Failure: `failed`

## Reversal Admission

```
pending → confirmed
```

## Settlement Cycles

FPS settles 3 times per business day via RTGS:

| Cycle | Window (UTC) | Settlement Date |
|---|---|---|
| 1 | 00:00–07:15 | Today |
| 2 | 07:15–13:00 | Today |
| 3 | 13:00–15:45 | Today |
| 1 (next day) | 15:45–24:00 | Next calendar day |

## Reconciliation Filters

Query `GET /v1/transaction/payments` with these parameters:

| Parameter | Type | Description |
|---|---|---|
| `filter[submission.settlement_cycle]` | int | Settlement cycle on the payment submission |
| `filter[submission.settlement_date]` | date | Settlement date on the payment submission |
| `filter[admission.settlement_cycle]` | int | Settlement cycle on the payment admission |
| `filter[admission.settlement_date]` | date | Settlement date on the payment admission |
| `page[number]` | int | Page number (1-based, default 1) |
| `page[size]` | int | Page size (default 100) |
