# memos-cli — Project Progress

## Overview

A CLI wrapper for the [Memos](https://www.usememos.com/) API v1 (targeting **v0.26.2**, latest stable).

**Verified against**: `usememos/memos` tag `v0.26.2` proto definitions.

---

## Phase 6: Support Comments, Reactions, Attachments `[x]`

### CLI Commands

```bash
# Memo namespace
memos memo list [--limit N] [--state STATE] [--output text|json|table]
memos memo get <id> [--include-comments] [--include-reactions] [--include-attachments] [-a]
memos memo create --content <text> [--visibility PRIVATE|PROTECTED|PUBLIC] [--pinned]
memos memo update <id> [--content <text>] [--visibility] [--pinned] [--unpin] [--state]
memos memo delete <id>

# Comments namespace
memos comments list <memo-id>
memos comments create <memo-id> --content <text>
memos comments delete <comment-id>

# Reactions namespace
memos reactions list <memo-id>
memos reactions create <memo-id> --type <emoji>
memos reactions delete <memo-id> <reaction-id>

# Attachments namespace
memos attachments list <memo-id>
memos attachments set <memo-id> --file <path> [--file <path>...]
```

### Implementation

1. [x] Add types (Reaction, Attachment, UpsertReaction, response types)
2. [x] Add read methods (ListMemoComments, ListMemoReactions, ListMemoAttachments)
3. [x] Add write methods (CreateMemoComment, UpsertMemoReaction, DeleteMemoReaction, SetMemoAttachments)
4. [x] Update `memos memo get` with `--include-comments`, `--include-reactions`, `--include-attachments`, `-a` flags
5. [x] Add `memos comments list/create/delete` commands
6. [x] Add `memos reactions list/create/delete` commands
7. [x] Add `memos attachments list/set` commands
8. [x] Add `memo` namespace to all memo commands (list, get, create, update, delete)

---

## Previous Phases

### Phase 1: Foundation `[x]`
- [x] Clean project structure and `go.mod`
- [x] `internal/api/types.go` — Memo, Visibility, State, request/response types
- [x] `internal/config/config.go` — TOML config load/save, path resolution

### Phase 2: API Client `[x]`
- [x] `internal/api/client.go` — HTTP client, generic `do()` method, auth headers
- [x] `internal/api/memo.go` — CreateMemo, ListMemos, GetMemo, UpdateMemo, DeleteMemo

### Phase 3: CLI Commands `[x]`
- [x] `cmd/root.go` — Root command, persistent flags
- [x] `cmd/list.go` — `memos memo list [--limit] [--state]`
- [x] `cmd/get.go` — `memos memo get <id>`
- [x] `cmd/create.go` — `memos memo create --content <text>`
- [x] `cmd/update.go` — `memos memo update <id>`
- [x] `cmd/delete.go` — `memos memo delete <id>`
- [x] `cmd/comments.go` — `memos comments list/create/delete`
- [x] `cmd/reactions.go` — `memos reactions list/create/delete`
- [x] `cmd/attachments.go` — `memos attachments list/set`

### Phase 4: Output Formatting `[x]`
- [x] `cmd/output.go` — Text, JSON, table output

### Phase 5: Testing `[x]`
- [x] `internal/api/memo_test.go` — 7 tests
- [x] `internal/config/config_test.go` — 4 tests

---

## API Reference (Memos v1)

| Method | Endpoint | Notes |
|---|---|---|
| `POST` | `/api/v1/memos` | CreateMemo |
| `GET` | `/api/v1/memos` | ListMemos |
| `GET` | `/api/v1/memos/{id}` | GetMemo |
| `PATCH` | `/api/v1/memos/{id}` | UpdateMemo (updateMask as query param) |
| `DELETE` | `/api/v1/memos/{id}` | DeleteMemo |
| `GET` | `/api/v1/memos/{id}/comments` | ListMemoComments |
| `POST` | `/api/v1/memos/{id}/comments` | CreateMemoComment |
| `GET` | `/api/v1/memos/{id}/reactions` | ListMemoReactions |
| `POST` | `/api/v1/memos/{id}/reactions` | UpsertMemoReaction |
| `DELETE` | `/api/v1/memos/{id}/reactions/{reaction}` | DeleteMemoReaction |
| `GET` | `/api/v1/memos/{id}/attachments` | ListMemoAttachments |
| `PATCH` | `/api/v1/memos/{id}/attachments` | SetMemoAttachments |

---

## Changelog

- **Phase 6 complete** — Added comments, reactions, attachments support with proper CLI namespace structure