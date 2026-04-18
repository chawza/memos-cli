# AGENTS.md

## Project

memos-cli — Go CLI for the [Memos](https://www.usememos.com/) API v1 (targeting v0.26.2+).

## Build & Run

```bash
go build -o memos .
go vet ./...
go test ./...
```

## Structure

- `main.go` — entry point
- `cmd/` — Cobra commands (one file per command + root.go + output.go)
- `internal/api/` — API client (client.go, memo.go, types.go)
- `internal/config/` — TOML config load/save

## Conventions

- Follow existing code style — no comments unless asked
- Use `internal/` for packages not meant to be imported externally
- Memo IDs are passed as plain strings by users, converted to `memos/{id}` internally via `memoName()` in `internal/api/memo.go`
- API types use direct request bodies (no wrapper): memo fields sent directly for create/update
- `updateMask` is passed as a query parameter on PATCH requests, not in the body
- Config file at `~/.config/memos-cli/config.toml` with 0600 permissions
- Auth priority: flags > env vars > config file

## Key Dependencies

- `github.com/spf13/cobra` — CLI framework
- `github.com/BurntSushi/toml` — TOML config (note: uses `NewEncoder().Encode()`, not `Marshal()`)

## API Reference

Targeting Memos v1 proto API. Verified against `usememos/memos` tag `v0.26.2`.

| Method | Endpoint | Notes |
|---|---|---|
| POST | `/api/v1/memos` | Body: memo fields directly |
| GET | `/api/v1/memos` | Query: pageSize, pageToken, filter, state, orderBy |
| GET | `/api/v1/memos/{id}` | |
| PATCH | `/api/v1/memos/{id}` | Body: memo fields directly; `updateMask` as query param |
| DELETE | `/api/v1/memos/{id}` | |

## Testing

Tests use `httptest.NewServer` for mocking the API. No external test dependencies.

```bash
go test ./... -v
```
