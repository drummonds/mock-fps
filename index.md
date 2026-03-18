# mock-fps Documentation

A Go mock server simulating the Form3 Faster Payments Service (FPS) API for testing.

Returns realistic JSON:API responses with async lifecycle status transitions, stand-in mode, and webhook notifications.

## FPS Events Reference

Complete reference for the FPS payment lifecycle — every event that can occur for an inbound payment.

- [FPS Events Reference](fps-events.html) — Payment types, status chains, webhook events, settlement cycles, stand-in mode, reconciliation
- [Resource States](states.html) — Quick reference for all resource states, settlement cycles, and reconciliation filters

## Stand-in Mode

Stand-in is a resilience mechanism where a gateway accepts inbound payments on behalf of a participant whose processor is unavailable. Payments are stored and forwarded when the processor recovers.

The mock server exposes `/admin/standin` for testing:

- `GET /admin/standin` — returns current status and queue length
- `PUT /admin/standin` — toggle stand-in on/off

When enabled, lifecycle transitions are queued rather than executed. When disabled, queued transitions drain and resume. See the [FPS Events Reference](fps-events.html#stand-in-mode) for full details on Form3's handling and acceptance qualifiers.

## Reconciliation

FPS uses deferred net settlement (DNS) with 3 settlement cycles per business day (07:15, 13:00, 15:45 via RTGS). Each cycle produces reconciliation data.

Form3 provides these endpoints for reconciliation (not yet implemented in mock-fps):

| Endpoint | Purpose |
|---|---|
| `GET /v1/transaction/payments` with settlement filters | Query by settlement_cycle, settlement_date |
| `GET /v1/notification/reports` | Scheme-generated settlement reports |
| `GET /v1/notification/reports/{id}/content` | Download report content |
| `GET /v1/organisation/positions` | Net settlement positions |

See the [FPS Events Reference](fps-events.html#fps-settlement-cycles) for cycle timing, risk controls, and filter parameters.

## Programmatic Access — fpsevents Package

The `fpsevents` package exports all FPS event definitions as Go types for use by other projects (e.g. gobank-products):

```go
import "github.com/nibble/mock-fps/fpsevents"

// All payment types
for _, pt := range fpsevents.PaymentTypes { ... }

// Status chains for each resource
chain := fpsevents.StatusChains["payment_submission"]

// All webhook event types
for _, evt := range fpsevents.WebhookEvents { ... }

// All 44 admission status reasons
for _, r := range fpsevents.AdmissionStatusReasons { ... }
```

## API Reference

Interactive API documentation with request/response schemas for all endpoints.

- [API Reference](api.html) — Full OpenAPI specification with Swagger UI

## Links

| | |
|---|---|
| Documentation | https://h3-mock-fps.statichost.eu |
| Source (Codeberg) | https://codeberg.org/hum3/mock-fps |
| Mirror (GitHub) | https://github.com/drummonds/mock-fps |
