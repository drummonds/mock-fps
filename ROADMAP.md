# Roadmap

## Current (v0.1.x)

- [x] Core payment lifecycle (create, submit, admit)
- [x] Returns, recalls, reversals with submissions
- [x] Webhook notifications via subscriptions
- [x] Stand-in mode (admin toggle, queue/drain)
- [x] Realistic FPS IDs (scheme_transaction_id, end_to_end_reference)
- [x] FPS events documentation

## Next (v0.2.x)

- [ ] Settlement cycle fields on submissions/admissions (settlement_cycle, settlement_date)
- [ ] Query parameter filtering on GET /v1/transaction/payments (settlement filters, date ranges, status, amount)
- [ ] Payment admission failure scenarios (configurable rejection)
- [ ] Return codes and rejection codes on responses

## Future

- [ ] Reports API (GET /v1/notification/reports) with canned settlement reports
- [ ] Positions endpoint (GET /v1/organisation/positions)
- [ ] Audit trail endpoint (GET /v1/audit/entries/{record_type})
- [ ] Qualified accept (stand-in) indicator on admissions
- [ ] Recall decision timeout enforcement (15-day deadline)
- [ ] Prometheus metrics endpoints
