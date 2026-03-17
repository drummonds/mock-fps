# mock-fps

A Go mock server simulating the Form3 Faster Payments Service (FPS) API for testing.

Returns realistic JSON:API responses with async lifecycle status transitions, stand-in mode, and webhook notifications.

## Documentation

- [FPS Events Reference](docs/fps-events.md) — Complete set of events for inbound payments, settlement cycles, stand-in, and reconciliation
- [API Reference](README.md#api-overview) — Endpoints and configuration

## Quick Start

```bash
task build
task run    # server on :8080
```

## Links

| | |
|---|---|
| Source (Codeberg) | https://codeberg.org/hum3/mock-fps |
| Mirror (GitHub) | https://github.com/drummonds/mock-fps |
