# memos-cli

A fast, minimal CLI wrapper for the [Memos](https://www.usememos.com/) API v1. Manage your self-hosted memos notes directly from the terminal.

## Install

```bash
go install github.com/chawza/memos-cli@latest
```

Or build from source:

```bash
git clone https://github.com/chawza/memos-cli.git
cd memos-cli
go build -o memos .
```

## Configuration

There are three ways to configure the CLI, checked in priority order:

1. **Flags** — `--base-url` and `--token` on any command
2. **Environment variables** — `MEMOS_BASE_URL` and `MEMOS_TOKEN`
3. **Config file** — `~/.config/memos-cli/config.toml`

### Auth command (recommended)

```bash
memos auth set --base-url https://memos.example.com --token your-access-token
```

This saves credentials to `~/.config/memos-cli/config.toml` and verifies connectivity.

Verify your configuration at any time:

```bash
memos auth check
```

Get your access token from **Memos web UI → Settings → Access Tokens**.

## Usage

All commands use a namespace pattern: `memos <resource> <verb>`

### Memo commands

```bash
memos memo create -c "Buy groceries #shopping"
memos memo create -c "Meeting notes" --visibility PRIVATE
memos memo create -c "Important update" --pinned
memos memo list
memos memo list --limit 50
memos memo list --state ARCHIVED
memos memo list --output json
memos memo list --output table
memos memo get 123
memos memo get 123 --output json
memos memo get 123 --include-comments --include-reactions -a
memos memo update 123 -c "Updated content"
memos memo update 123 --visibility PUBLIC
memos memo update 123 --pinned
memos memo update 123 --unpin
memos memo update 123 --state ARCHIVED
memos memo delete 123
```

### Comments commands

```bash
memos comments list 123
memos comments create 123 -c "Great point!"
memos comments delete 456
```

### Reactions commands

```bash
memos reactions list 123
memos reactions create 123 --type "👍"
memos reactions delete 123 <reaction-id>
```

### Attachments commands

```bash
memos attachments list 123
memos attachments set 123 --file /path/to/file.pdf
```

## Output formats

The `list` and `get` commands support `--output` / `-o` with three formats:

| Format | Flag | Description |
|---|---|---|
| Text | `--output text` (default) | Human-readable, compact |
| JSON | `--output json` | Machine-readable, indented |
| Table | `--output table` | Aligned columns (list only) |

### Text output examples

**`memos memo get 123`:**
```
Name:       memos/123
State:      NORMAL
Visibility: PRIVATE
Pinned:     false
Creator:    users/1
Created:    2025-04-18T10:30:00Z

Buy groceries

Reactions:
  👍 by users/1

Attachments:
  receipt.pdf (application/pdf)
```

**`memos memo list`:**
```
  [PRIVATE] Buy groceries
* [PUBLIC]  Important announcement
  [PRIVATE] Meeting notes
```

## Commands reference

| Command | Description |
|---|---|
| `memos auth set` | Save credentials to config file |
| `memos auth check` | Verify saved configuration |
| `memos memo list` | List memos with filters |
| `memos memo get <id>` | Get a memo (with optional related data) |
| `memos memo create` | Create a new memo |
| `memos memo update <id>` | Update a memo |
| `memos memo delete <id>` | Delete a memo |
| `memos comments list <memo-id>` | List comments |
| `memos comments create <memo-id>` | Create a comment |
| `memos comments delete <comment-id>` | Delete a comment |
| `memos reactions list <memo-id>` | List reactions |
| `memos reactions create <memo-id>` | Add a reaction |
| `memos reactions delete <memo-id> <reaction-id>` | Remove a reaction |
| `memos attachments list <memo-id>` | List attachments |
| `memos attachments set <memo-id>` | Set attachments (replaces all) |

### Global flags

| Flag | Env variable | Description |
|---|---|---|
| `--base-url` | `MEMOS_BASE_URL` | Memos instance URL |
| `--token` | `MEMOS_TOKEN` | Access token |

### Command flags

**`auth set`**
| Flag | Default | Description |
|---|---|---|
| `--base-url` | (required) | Memos instance URL |
| `--token` | (required) | Access token |
| `--timeout` | `30` | HTTP timeout in seconds |
| `--tls-skip-verify` | `false` | Skip TLS certificate verification |

**`memo list`**
| Flag | Default | Description |
|---|---|---|
| `--limit` | `20` | Max memos to return |
| `--state` | | Filter by state: `NORMAL` or `ARCHIVED` |
| `-o, --output` | `text` | Output format: `text`, `json`, `table` |

**`memo get`**
| Flag | Default | Description |
|---|---|---|
| `-o, --output` | `text` | Output format: `text`, `json` |
| `--include-comments` | `false` | Include comments |
| `--include-reactions` | `false` | Include reactions |
| `--include-attachments` | `false` | Include attachments |
| `-a` | `false` | Include all (comments, reactions, attachments) |

**`memo create`**
| Flag | Default | Description |
|---|---|---|
| `-c, --content` | (required) | Memo content in Markdown |
| `--visibility` | `PRIVATE` | Visibility: `PRIVATE`, `PROTECTED`, `PUBLIC` |
| `--pinned` | `false` | Pin the memo |

**`memo update`**
| Flag | Description |
|---|---|
| `-c, --content` | New content |
| `--visibility` | New visibility: `PRIVATE`, `PROTECTED`, `PUBLIC` |
| `--pinned` | Pin the memo |
| `--unpin` | Unpin the memo |
| `--state` | New state: `NORMAL` or `ARCHIVED` |

**`comments create`**
| Flag | Default | Description |
|---|---|---|
| `-c, --content` | (required) | Comment content in Markdown |

**`reactions create`**
| Flag | Default | Description |
|---|---|---|
| `--type` | (required) | Reaction emoji |

**`attachments set`**
| Flag | Default | Description |
|---|---|---|
| `--file` | (required) | File path(s) to attach |

## Tech stack

- Go 1.21+
- [Cobra](https://github.com/spf13/cobra) — CLI framework
- [BurntSushi/toml](https://github.com/BurntSushi/toml) — Config parsing
- Go stdlib `net/http` — HTTP client

## License

MIT
