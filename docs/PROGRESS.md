# memos-cli — Project Progress

## Overview

A CLI wrapper for the [Memos](https://www.usememos.com/) API v1 (targeting **v0.26.2**, latest stable).
Full rewrite from scratch — the initial skeleton had compilation errors and incorrect API usage.

**Verified against**: `usememos/memos` tag `v0.26.2` proto definitions (`proto/api/v1/memo_service.proto`, `proto/api/v1/common.proto`).

**Scope**: Basic CRUD (create, get, list, update, delete) + config management + tests.
**Out of scope**: Search command, stdin support, reactions, attachments, relations, comments.

---

## Phases

### Phase 1: Foundation `[x]`
> Project structure, Go module init, types, and config.

- [x] Clean project structure and `go.mod`
- [x] `internal/api/types.go` — Memo, Visibility, State, request/response types
- [x] `internal/config/config.go` — TOML config load/save, path resolution

### Phase 2: API Client `[x]`
> HTTP client with proper auth, error handling, and CRUD methods.

- [x] `internal/api/client.go` — HTTP client, generic `do()` method, auth headers
- [x] `internal/api/memo.go` — CreateMemo, ListMemos, GetMemo, UpdateMemo, DeleteMemo
- [x] Proper error handling (API errors, network errors, decode errors)

### Phase 3: CLI Commands `[x]`
> Cobra commands for each CRUD operation.

- [x] `cmd/root.go` — Root command, persistent flags (`--base-url`, `--token`), env var fallbacks, client resolution
- [x] `cmd/create.go` — `memos create --content <text> [--visibility] [--pinned]`
- [x] `cmd/get.go` — `memos get <id>`
- [x] `cmd/list.go` — `memos list [--limit] [--filter] [--state] [--output]`
- [x] `cmd/update.go` — `memos update <id> [--content] [--visibility] [--pinned] [--unpin] [--state]`
- [x] `cmd/delete.go` — `memos delete <id>`
- [x] `main.go` — Entry point wiring

### Phase 4: Output Formatting `[x]`
> Multi-format output for list/get commands.

- [x] `cmd/output.go` — Shared output helpers (text, json, table)
- [x] Text output for `list` (compact one-line per memo)
- [x] Text output for `get` (detailed view)
- [x] JSON output (`--output json`)
- [x] Table output (`--output table`)

### Phase 5: Testing `[x]`
> Unit tests for API client and config.

- [x] `internal/api/memo_test.go` — Test CRUD methods with mocked HTTP server (7 tests)
- [x] `internal/config/config_test.go` — Test config load/save (4 tests)

### Phase 6: Polish & Ship `[x]`
> Final cleanup, documentation, build verification.

- [x] Verify `go build` and binary works
- [x] `go vet` passes clean
- [x] All 11 tests pass
- [x] Consistent error messages and exit codes

---

## Key Design Decisions

| Decision | Choice | Rationale |
|---|---|---|
| Memo identifier | `name` field (format: `memos/{id}`) | Matches Memos v1 proto API |
| Config format | TOML | Simple, human-readable, already in README |
| Config location | `~/.config/memos-cli/config.toml` | XDG convention |
| Auth priority | Flags > Env vars > Config file | Most explicit wins |
| HTTP client | Go stdlib `net/http` | No external dependency needed |
| CLI framework | Cobra | Standard Go CLI framework |
| Search command | Skipped | No server-side full-text search in API |

---

## API Reference (Memos v1 — verified against v0.26.2)

Proto source: `github.com/usememos/memos` tag `v0.26.2`

### Endpoints used:
| Method | Endpoint | Request | Notes |
|---|---|---|---|
| `POST` | `/api/v1/memos` | `{"memo": {...}, "memoId": "optional"}` | CreateMemo |
| `GET` | `/api/v1/memos` | Query: `pageSize`, `pageToken`, `filter`, `state`, `orderBy` | ListMemos |
| `GET` | `/api/v1/{name=memos/*}` | — | GetMemo, name = `memos/{id}` |
| `PATCH` | `/api/v1/{memo.name=memos/*}` | `{"memo": {...}, "updateMask": "content,visibility"}` | UpdateMemo |
| `DELETE` | `/api/v1/{name=memos/*}` | — | DeleteMemo |

### Memo resource fields (v0.26.2):
- `name` (string, format `memos/{id}`) — resource identifier
- `state` (enum: `NORMAL`, `ARCHIVED`)
- `creator` (string, output only, format `users/{user}`)
- `create_time`, `update_time`, `display_time` (timestamps)
- `content` (string, markdown, required)
- `visibility` (enum: `PRIVATE`, `PROTECTED`, `PUBLIC`)
- `tags` (repeated string, output only)
- `pinned` (bool)
- `snippet` (string, output only)

### CreateMemoRequest body mapping:
```json
{
  "memo": {
    "content": "Hello world",
    "visibility": "PRIVATE",
    "pinned": false
  }
}
```

### UpdateMemoRequest body mapping:
```json
{
  "memo": {
    "name": "memos/123",
    "content": "Updated content"
  },
  "updateMask": "content,visibility"
}
```

---

## Changelog

- **Phase 1–6 complete** — Full rewrite implemented, all tests passing, builds cleanly
