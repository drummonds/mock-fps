# Mock FPS — Documentation & Feature Plan

## Context

The mock-fps project simulates the Form3 Faster Payments Service API for testing.
It currently has no README.md or docs/ folder. This plan covers creating those,
plus research findings on two Form3 API capabilities we may want to mock.

---

## Task 1: Create README.md

Status: **TODO**

Should cover:
- What the project is (mock Form3 FPS API for testing)
- Quick start (`task run`, env vars)
- Supported endpoints (30 routes under `/v1`)
- Architecture overview (handlers, store, lifecycle engine, webhooks)
- Status transition chains
- Configuration (env vars from `internal/config/config.go`)

### Current endpoints
- Payments: CRUD + list
- Payment submissions, admissions (with task patching)
- Returns + return submissions
- Recalls + recall submissions
- Recall decisions + decision submissions
- Reversals + reversal submissions
- Subscriptions: CRUD + list + patch + delete
- Health: GET /health

---

## Task 2: Create docs/ folder

Status: **TODO**

Suggested structure:
```
docs/
  architecture.md    — internal design, middleware chain, store keys
  api-reference.md   — endpoint table with methods/paths/descriptions
  form3-comparison.md — what we mock vs what real Form3 API provides
```

---

## Task 3: Research Findings — Standin Mode

Status: **RESEARCH COMPLETE**

### Finding: Form3 does NOT expose standin in its public API

- Searched the full Form3 OpenAPI spec (158 paths, 701KB) for `standin`, `stand-in`,
  `stand_in`, `contingency`, `fallback`, `outage`, `degraded` — zero matches.
- Form3 handles stand-in processing internally (queues payments when FPS central
  infrastructure / VocaLink is unavailable) and delivers when it recovers.
- The only visibility is via **Prometheus metrics** endpoints:
  - `GET /v1/metrics/prometheus/api/v1/query` (point-in-time)
  - `GET /v1/metrics/prometheus/api/v1/query_range` (range)
  - `GET /v1/metrics/prometheus/federate`
  - These expose stand-in queue depths, processing times, volumes.
- At the scheme level (Pay.UK / VocaLink), stand-in is an infrastructure contingency
  mechanism, not a formally documented API concept.

### Implication for mock-fps
- No standin endpoints to mock from the payments API.
- Could optionally mock the Prometheus metrics endpoints if we want to simulate
  monitoring, but this is low priority.
- If we want to simulate degraded mode for testing, we could add a custom
  `/admin/standin` toggle that queues submissions instead of processing them
  immediately — but this would be our own invention, not a Form3 API feature.

---

## Task 4: Research Findings — Transaction Listing & Reconciliation

Status: **RESEARCH COMPLETE**

### Finding: Form3 has extensive reconciliation support

#### A. `GET /v1/transaction/payments` — primary reconciliation endpoint

Already partially mocked (our `payments.List` handler). The real Form3 API supports
rich filtering that we don't yet implement:

**Settlement-cycle filters (key for reconciliation):**
- `filter[submission.settlement_cycle]` — integer (e.g. `1`)
- `filter[submission.settlement_date]` — date
- `filter[admission.settlement_cycle]` — integer
- `filter[admission.settlement_date]` — date
- `filter[return_submission.settlement_cycle]`
- `filter[return_admission.settlement_cycle]`

**Date-range filters:**
- `filter[processing_date_from]` / `filter[processing_date_to]`
- `filter[created_date_from]` / `filter[created_date_to]`
- `filter[submission.submission_date_from]` / `filter[submission.submission_date_to]`
- `filter[admission.admission_date_from]` / `filter[admission.admission_date_to]`

**Other useful filters:**
- `filter[organisation_id]` — UUID array
- `filter[payment_scheme]` — e.g. "FPS"
- `filter[amount]`, `filter[amount][from]`, `filter[amount][to]`
- `filter[submission.status]`, `filter[admission.status]`
- `filter[relationships]` / `filter[not_relationships]`

**Pagination:**
- `page[number]`, `page[size]`, `page[before]`, `page[after]` (cursor-based)
- Limit: `page[number] * page[size]` must be < 10,000

#### B. Reports API — `GET /v1/notification/reports`

Not currently mocked. Form3 provides scheme-generated settlement reports as resources:

| Endpoint | Method | Description |
|---|---|---|
| `/v1/notification/reports` | GET | List reports with filters |
| `/v1/notification/reports/{id}` | GET | Fetch a specific report |
| `/v1/notification/reports/{id}/admissions/{admissionId}` | GET | Report admission |
| `/v1/notification/reports/{id}/content` | GET | Download report content (binary) |

Report filters: `report_type`, `report_source`, `organisation_id`,
`created_on_after/before`, `processing_date_from/to`.

In the FPS Direct simulator, Form3 creates a dummy inbound report every business day.

#### C. Audit Trail — `GET /v1/audit/entries/{record_type}`

Not currently mocked. Provides change history per resource type.

#### D. Positions — `GET /v1/organisation/positions`

Not currently mocked. Shows net settlement positions.

### Implication for mock-fps

**High value additions:**
1. Add query parameter filtering to `GET /v1/transaction/payments` — settlement_cycle,
   settlement_date, date ranges, status, amount. This directly supports reconciliation
   testing against a downstream service.
2. Add `settlement_cycle` and `settlement_date` fields to submission/admission models.

**Medium value additions:**
3. Mock the Reports API (`/v1/notification/reports`) with canned settlement reports.
4. Add basic audit trail endpoint.

**Low value additions:**
5. Positions endpoint.
6. Prometheus metrics endpoints.

---

## Execution Order

1. Create `README.md`
2. Create `docs/` folder with architecture, API reference, and Form3 comparison docs
3. (Future) Add settlement_cycle/settlement_date to models
4. (Future) Add query filtering to payments list endpoint
5. (Future) Mock reports API

---

## Sources

- [Form3 API Docs](https://www.api-docs.form3.tech/)
- [Form3 FPS Direct Payments](https://www.api-docs.form3.tech/api/schemes/fps-direct/payments/overview)
- [Form3 Reports Tutorial](https://api-docs.form3.tech/tutorial-reports.html)
- [Form3 Subscriptions Tutorial](https://api-docs.form3.tech/tutorial-subscription-event-notification.html)
- [Form3 Metrics Blog](https://www.form3.tech/news/payment-insights/how-form3s-metrics-service-opens-up-granular-insights)
- [Form3 OpenAPI Spec](https://www.api-docs.form3.tech/assets/swagger/form3-swagger.yaml) (158 paths)
- [go-form3 Client](https://github.com/form3tech-oss/go-form3)
