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
- `cmd/` — Cobra commands (one file per command + root.go + output.go + memo.go)
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

## Issue Fix Workflow

When fixing issues from GitHub:
1. Create a branch for the fix (do not commit directly to main)
2. Make changes and verify with build/vet/test
3. Create PR on GitHub and wait for merge
4. Do NOT push commits until PR is created and ready for review

## Key Dependencies

- `github.com/spf13/cobra` — CLI framework
- `github.com/BurntSushi/toml` — TOML config (note: uses `NewEncoder().Encode()`, not `Marshal()`)

## API Reference

Targeting Memos v1 proto API. Verified against `usememos/memos` tag `v0.26.2`.

### Memo Endpoints

| Method | Endpoint | Notes |
|---|---|---|
| POST | `/api/v1/memos` | Body: memo fields directly |
| GET | `/api/v1/memos` | Query: pageSize, pageToken, state |
| GET | `/api/v1/memos/{id}` | |
| PATCH | `/api/v1/memos/{id}` | Body: memo fields directly; `updateMask` as query param |
| DELETE | `/api/v1/memos/{id}` | |

### Comments Endpoints

| Method | Endpoint | Notes |
|---|---|---|
| GET | `/api/v1/memos/{id}/comments` | List comments |
| POST | `/api/v1/memos/{id}/comments` | Create comment |

### Reactions Endpoints

| Method | Endpoint | Notes |
|---|---|---|
| GET | `/api/v1/memos/{id}/reactions` | List reactions |
| POST | `/api/v1/memos/{id}/reactions` | Upsert reaction |
| DELETE | `/api/v1/memos/{id}/reactions/{reaction}` | Delete reaction |

### Attachments Endpoints

| Method | Endpoint | Notes |
|---|---|---|
| GET | `/api/v1/memos/{id}/attachments` | List attachments |
| PATCH | `/api/v1/memos/{id}/attachments` | Set attachments (replaces all) |

## CLI Command Structure

All commands use namespace pattern: `memos <resource> <verb>`

```
memos memo list|get|create|update|delete
memos comments list|create|delete
memos reactions list|create|delete
memos attachments list|set
```

## Testing

Tests use `httptest.NewServer` for mocking the API. No external test dependencies.

```bash
go test ./... -v
```

## Releases

In every release make sure CHANGELOG.md is updated
