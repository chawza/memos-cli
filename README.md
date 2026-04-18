# memos-cli

A fast, minimal CLI wrapper for the [Memos](https://www.usememos.com/) API. Manage your self-hosted memos notes directly from the terminal.

> **Status**: Early planning — not yet implemented.

## Why

[Memos](https://www.usememos.com/) is a brilliant self-hosted note-taking app, but there's no convenient CLI to manage memos hands-free. This CLI fills that gap — scriptable, pipe-friendly, and designed for developers who live in the terminal.

Potential use cases:
- Quick note capture from the shell / scripts
- Mirror tasks from other tools (e.g. Linear, cron job outputs) into memos
- Search and retrieve past memos without opening a browser
- Automation workflows (CI/CD pipelines, notifications, reminders)

## Features (Planned)

### Core Commands

- [ ] `memos list` — List memos with optional filters (tag, visibility, date range)
- [ ] `memos create` — Create a new memo from stdin or flag
- [ ] `memos get` — Get a single memo by ID
- [ ] `memos update` — Update a memo's content, visibility, or state
- [ ] `memos delete` — Delete a memo by ID
- [ ] `memos search` — Full-text search across memos

### Output Formats

- [ ] JSON (machine-readable, `--output json`)
- [ ] Plain text (`--output text`)
- [ ] Formatted table (`--output table`)

### Config & Auth

- [ ] Config file (`~/.config/memos-cli/config.toml`) for base URL and token
- [ ] `--token` flag for on-the-fly auth
- [ ] `--base-url` flag to override configured instance

## Architecture

```
memos-cli/
├── cmd/
│   └── main.go           # CLI entry point, cobra root
├── internal/
│   ├── api/
│   │   └── client.go     # Memos API v1 client
│   └── config/
│       └── config.go     # Config loading (TOML)
├── go.mod
└── README.md
```

### API Coverage (from usememos.com/docs/api/latest)

The Memos API is gRPC-inspired REST with a `POST /api/v1/memos` base. Key endpoints:

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/api/v1/memos` | List memos (pagination, AIP-160 filter) |
| `POST` | `/api/v1/memos` | Create memo |
| `GET` | `/api/v1/memos/{id}` | Get memo by ID |
| `PATCH` | `/api/v1/memos/{id}` | Update memo (field mask) |
| `DELETE` | `/api/v1/memos/{id}` | Delete memo |
| `POST` | `/api/v1/memos/{id}/reactions` | Add reaction |
| `DELETE` | `/api/v1/memos/{id}/reactions/{reactionId}` | Remove reaction |

**Request shape** (CreateMemo):
```json
{
  "content": "string (required, markdown)",
  "visibility": "PRIVATE | PROTECTED | PUBLIC",
  "state": "NORMAL | ARCHIVED",
  "pinned": false,
  "createTime": "2024-01-01T00:00:00Z"
}
```

**Auth**: Bearer token — obtained from Memos web UI → Settings → Access Tokens.

### Future Ideas

- [ ] **MCP server mode** — expose memos tools via the Model Context Protocol, enabling an AI agent (like this one!) to manage your memos natively
- [ ] **Watch mode** — long-poll or webhook-driven display of new memos
- [ ] **Import/Export** — bulk import from JSON/Markdown, export all memos
- [ ] **Tags** — first-class tag support if Memos adds it natively
- [ ] **Web UI fallback** — `memos open` to open the web UI

## Quick Start (Planned)

```bash
# Configure once
export MEMOS_BASE_URL="https://memos.example.com"
export MEMOS_TOKEN="your-access-token"

# List recent memos
memos list --limit 10

# Create from stdin
echo "Buy groceries" | memos create --visibility PRIVATE

# Create from flag
memos create --content "Meeting at 3pm" --pinned

# Get a specific memo
memos get 123

# Search
memos search "quarterly report"

# Delete
memos delete 456
```

## Tech Stack

- **Language**: Go 1.21+
- **CLI framework**: [Cobra](https://github.com/spf13/cobra)
- **Config format**: TOML (via [BurntSushi/toml](https://github.com/BurntSushi/toml))
- **HTTP client**: Go stdlib `net/http`
- **JSON**: Go stdlib `encoding/json`

## Install (Planned)

```bash
go install github.com/chawza/memos-cli@latest
```

Or download a prebuilt binary from Releases.

## License

MIT
