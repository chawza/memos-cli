# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v0.3.0] - 2026-04-18

### Added
- JSON fixtures for realistic API request/response testing (`internal/api/fixtures/`)
- `Tags`, `Property`, and `Location` fields to `Memo` struct
- `MemoProperty` type (hasLink, hasTaskList, hasCode, hasIncompleteTasks, title)
- `MemoLocation` type (placeholder, latitude, longitude)
- Structured error parsing in `APIError` (code, message, details)
- Tests: empty list, filter/state query params, pagination, structured errors, HTTP error codes (400/401/403/500), Ping/PingFailure

### Changed
- All tests now use fixtures matching the official Memos API v1 spec
- `APIError.Error()` format includes HTTP status code and API error code

## [v0.2.0] - 2026-04-18

### Fixed
- Memo content not being transmitted during creation — removed wrapper types, send memo fields directly in request body
- `updateMask` passed as query parameter on PATCH requests instead of in body

### Changed
- Removed `CreateMemoRequest` and `UpdateMemoRequest` wrapper types
- Updated `UpdateMemo` client method signature to accept `updateMask` as separate parameter
- Updated `AGENTS.md` with correct API request format documentation

## [v0.1.0] - 2026-04-18

### Added
- Initial Go CLI scaffold with Memos API v1 client
- Full CRUD operations: create, list, get, update, delete memos
- `auth` command with set/check subcommands
- HTTP config options (timeout, TLS skip verify)
- TOML config support at `~/.config/memos-cli/config.toml`
- README and AGENTS documentation
- Release configuration

[v0.3.0]: https://github.com/chawza/memos-cli/releases/tag/v0.3.0
[v0.2.0]: https://github.com/chawza/memos-cli/releases/tag/v0.2.0
[v0.1.0]: https://github.com/chawza/memos-cli/releases/tag/v0.1.0
