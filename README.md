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

You can also set optional HTTP options:

```bash
memos auth set --base-url https://memos.example.com --token your-access-token --timeout 60 --tls-skip-verify
```

Verify your configuration at any time:

```bash
memos auth check
```

### Config file (manual)

```bash
mkdir -p ~/.config/memos-cli
cat > ~/.config/memos-cli/config.toml << 'EOF'
base_url = "https://memos.example.com"
token = "your-access-token"
timeout = 30
tls_skip_verify = false
EOF
```

Get your access token from **Memos web UI → Settings → Access Tokens**.

## Usage

### Create a memo

```bash
memos create -c "Buy groceries"
memos create -c "Meeting notes" --visibility PRIVATE
memos create -c "Important update" --pinned
```

### List memos

```bash
memos list
memos list --limit 50
memos list --state ARCHIVED
memos list --output json
memos list --output table
```

### Get a memo

```bash
memos get 123
memos get 123 --output json
```

### Update a memo

```bash
memos update 123 -c "Updated content"
memos update 123 --visibility PUBLIC
memos update 123 --pinned
memos update 123 --unpin
memos update 123 --state ARCHIVED
memos update 123 -c "New text" --visibility PUBLIC --pinned
```

### Delete a memo

```bash
memos delete 123
```

## Output formats

The `list` and `get` commands support `--output` / `-o` with three formats:

| Format | Flag | Description |
|---|---|---|
| Text | `--output text` (default) | Human-readable, compact |
| JSON | `--output json` | Machine-readable, indented |
| Table | `--output table` | Aligned columns (list only) |

### Text output examples

**`memos get 123`:**
```
Name:       memos/123
State:      NORMAL
Visibility: PRIVATE
Pinned:     false
Creator:    users/1
Created:    2025-04-18T10:30:00Z

Buy groceries
```

**`memos list`:**
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
| `memos create` | Create a new memo |
| `memos list` | List memos with filters |
| `memos get <id>` | Get a single memo |
| `memos update <id>` | Update a memo |
| `memos delete <id>` | Delete a memo |

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

**`auth check`**
| Flag | Default | Description |
|---|---|---|
| `--base-url` | | Override base URL (optional) |
| `--token` | | Override token (optional) |

**`create`**
| Flag | Default | Description |
|---|---|---|
| `-c, --content` | (required) | Memo content in Markdown |
| `--visibility` | `PRIVATE` | Visibility: `PRIVATE`, `PROTECTED`, `PUBLIC` |
| `--pinned` | `false` | Pin the memo |

**`list`**
| Flag | Default | Description |
|---|---|---|
| `--limit` | `20` | Max memos to return |
| `--filter` | | CEL filter expression |
| `--state` | | Filter by state: `NORMAL` or `ARCHIVED` |
| `-o, --output` | `text` | Output format: `text`, `json`, `table` |

**`get`**
| Flag | Default | Description |
|---|---|---|
| `-o, --output` | `text` | Output format: `text`, `json` |

**`update`**
| Flag | Description |
|---|---|
| `-c, --content` | New content |
| `--visibility` | New visibility: `PRIVATE`, `PROTECTED`, `PUBLIC` |
| `--pinned` | Pin the memo |
| `--unpin` | Unpin the memo |
| `--state` | New state: `NORMAL` or `ARCHIVED` |

## Tech stack

- Go 1.21+
- [Cobra](https://github.com/spf13/cobra) — CLI framework
- [BurntSushi/toml](https://github.com/BurntSushi/toml) — Config parsing
- Go stdlib `net/http` — HTTP client

## License

MIT
