# HW2

## Run

Start server:

```bash
cd hw2
go run ./cmd/server
```

Run client (server should already be running):

```bash
cd hw2
go run ./cmd/client
```

## Server API

- `GET /version` → `{"version":"v1.0.0"}`
- `POST /decode` (JSON `{"inputString":"<base64>"}`) → `{"outputString":"<decoded>"}`
- `GET /hard-op` → sleeps randomly 10..20 seconds and returns random HTTP `200` or `5xx`

Graceful shutdown: `Ctrl+C`. Port can be controlled via `PORT` environment variable.
