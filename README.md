# mock-fps

A Go mock server simulating the Form3 Faster Payments Service (FPS) API. Returns realistic JSON:API responses with async lifecycle status transitions, stand-in mode, and webhook notifications.

Full documentation: [mock-fps-docs](https://github.com/drummonds/mock-fps-docs) (coming soon)

## Quick start

```bash
task build
task run
```

The server starts on `:8080` by default.

## Configuration

| Variable | Default | Description |
|---|---|---|
| `PORT` | `8080` | HTTP listen port |
| `LIFECYCLE_STEP_DELAY_MS` | `500` | Delay between status transitions |
| `WEBHOOK_WORKERS` | `4` | Concurrent webhook delivery workers |
| `WEBHOOK_BUFFER_SIZE` | `1000` | Webhook notification queue size |

## API overview

All endpoints are under `/v1` and follow JSON:API conventions (`application/vnd.api+json`).

- **Payments** — `POST/GET /v1/transaction/payments`
- **Submissions** — `POST/GET .../submissions`
- **Admissions** — `POST/GET/PATCH .../admissions`
- **Returns** — `POST/GET .../returns` and `.../returns/{id}/submissions`
- **Recalls** — `POST/GET .../recalls`, `.../submissions`, `.../decisions`
- **Reversals** — `POST/GET .../reversals` and `.../reversals/{id}/submissions`
- **Subscriptions** — `CRUD /v1/notification/subscriptions`
- **Health** — `GET /health`
- **Admin** — `GET/PUT /admin/standin` (stand-in mode control)

## Tasks

```
task build    # Build binary
task run      # Build and run
task test     # Run tests with -race
task check    # Test + clean
task clean    # Remove build artifacts
```

## License

See [LICENSE](LICENSE).
