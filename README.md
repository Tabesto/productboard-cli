# pboard - ProductBoard CLI

A command-line tool for browsing ProductBoard data. **Read-only by default**: supports all GET endpoints, plus optional, gated create/update of entities and notes (API V2 only, see [Write access](#write-access)).

Uses ProductBoard API V2 by default (`--api-version 1`, `PRODUCTBOARD_API_VERSION` or `api_version` in the config select V1, which ProductBoard retired on 2026-07-08).

## Installation

### Homebrew (macOS / Linux)

```bash
brew tap tabesto/tap
brew install productboard-cli
```

### From source

```bash
go install github.com/tabesto/productboard-cli/cmd/pboard@latest
```

## Configuration

Set your ProductBoard API token using one of these methods:

### One-command setup

```bash
pboard configure pb_your_token_here
```

### Interactive setup

```bash
pboard configure
```

Both methods save the token to `~/.config/pboard/config.yaml` (mode 600).

### Environment variable

```bash
export PRODUCTBOARD_API_TOKEN="pb_your_token_here"
```

The environment variable takes precedence over the config file.

## Usage

### Features

```bash
pboard features list
pboard features list --status-name "In Progress" --limit 10
pboard features get <id>
pboard features links initiatives <id>
pboard features links objectives <id>
```

### Product hierarchy

```bash
pboard products list
pboard products get <id>
pboard components list
pboard feature-statuses list
```

### Notes & feedback

```bash
pboard notes list --term "onboarding" --created-from 2024-01-01
pboard notes get <id>
pboard notes tags <id>
pboard notes links <id>
pboard feedback-forms list
```

### Releases

```bash
pboard releases list
pboard release-groups list
pboard feature-release-assignments list --release-id <id>
```

### OKRs & initiatives

```bash
pboard objectives list
pboard key-results list
pboard initiatives list
pboard objectives links features <id>
pboard initiatives links objectives <id>
```

### Companies & users

```bash
pboard companies list --term "Acme"
pboard users list
```

### Custom fields

```bash
pboard custom-fields list --type text
pboard custom-fields values list --custom-field-id <id>
```

### Integrations & webhooks

```bash
pboard plugin-integrations list
pboard jira-integrations list
pboard webhooks list
```

## Output formats

```bash
# Default: human-readable table
pboard features list

# JSON for scripting
pboard features list -o json

# Pipe to jq
pboard features list -o json | jq '.[] | select(.status.name == "Done")'

# Limit results
pboard features list --limit 10
```

## Global flags

| Flag | Short | Description |
|------|-------|-------------|
| `--output` | `-o` | Output format: `table` (default) or `json` |
| `--limit` | `-l` | Maximum number of results (0 = all) |
| `--api-version` | | API version: `1` or `2` (default `2`) |
| `--allow-write` | | Enable write operations for this run (see [Write access](#write-access)) |
| `--help` | `-h` | Show help |
| `--version` | `-v` | Show version |

## Write access

Writes are **disabled by default**: any create/update is refused (exit code 6) before a request is sent. Enable them with any one of:

| Mechanism | Scope |
|-----------|-------|
| `--allow-write` flag | this run only |
| `pboard configure --write-enabled` (`=false` to disable) | persisted as `write_enabled` in `~/.config/pboard/config.yaml` |
| `PRODUCTBOARD_WRITE_ENABLED=1` | environment (`0`/`false`/other values force writes off, overriding the config) |

Supported operations (API V2 only; field names depend on your workspace, so inspect existing records first):

```bash
pboard --allow-write entities create --type feature --field name="My feature"
pboard --allow-write entities update <id> --field name="Renamed"
pboard --allow-write entities update <id> --patch '[{"op":"clear","path":"owner"}]'
pboard --allow-write notes create --name "Customer asked for X" --content "Details"
pboard --allow-write notes update <id> --field content="Edited"
```

`--field key=value` is repeatable (the value is read as JSON when valid, otherwise as text). `--json '<data object>'` sends a full request `data` object. An update takes exactly one of `--field`, `--patch` or `--json`. Create/update print only the new record's id, type and link; use the usual `get` commands for the full record. Deleting and relationship management are not supported.

### Audit log and reverting

Every write (CLI and MCP) is recorded in `~/.config/pboard/audit.log` (mode 600, one JSON object per line, append-only, never rotated). Each write has a `request` line written **before** the call and a `response` line (same `op_id`) with the returned id/link or the error. For updates, the `request` line also holds `previous`, the record as read just before the change, and `revert`, a request body that puts back the previous value of each field that was replaced.

The log is fail-closed: if the log cannot be written, or the current value cannot be read before an update, the write is not sent. Note contents sent on create/update end up in the log, so treat the file as sensitive.

```bash
# list recent updates (op_id, target, revert body)
jq -c 'select(.event=="request" and .method=="PATCH") | {op_id, path, revert}' ~/.config/pboard/audit.log | tail -5
# undo one (check the body first: it is built from the read format and has not been verified against the real API)
pboard --allow-write entities update <id> --json "$(jq -c 'select(.op_id=="<op_id>" and .event=="request") | .revert.data' ~/.config/pboard/audit.log)"
```

Limits: a field that did not exist before an update is not in `revert` (clear it by hand); updates made with `--patch` keep the full `previous` snapshot but no ready-made `revert`; a create cannot be undone by pboard (delete the record in Productboard using the id in the `response` line).

## MCP Server for Claude Desktop

Use ProductBoard data directly from Claude Desktop as MCP tools. One command to install:

```bash
pboard mcp install
```

Then restart Claude Desktop. All ProductBoard tools (features, notes, companies, releases, etc.) will be available in your conversations.

```bash
# Preview what would be configured
pboard mcp install --dry-run

# Overwrite existing installation
pboard mcp install --force

# Remove the MCP server
pboard mcp uninstall

# Also expose write tools (create_entity, update_entity, create_note, update_note)
pboard --allow-write mcp install
```

Write tools are registered only when writes are enabled when the server starts (`--allow-write` on `mcp serve`/`mcp install`, `write_enabled` in the config, or `PRODUCTBOARD_WRITE_ENABLED=1`); otherwise the assistant cannot see them.

The MCP server exposes ~50 tools covering every pboard CLI command. Results default to 25 items per query to keep Claude's context manageable.

## Claude Code integration

Install a skill that teaches Claude Code agents how to use `pboard`:

```bash
pboard skill install
```

This writes a usage guide to `~/.claude/commands/pboard.md`, making it available across all Claude Code sessions. Agents will automatically know every command, flag, and workflow.

```bash
# Overwrite an existing skill
pboard skill install --force

# Preview without writing
pboard skill install --dry-run

# Remove the skill
pboard skill uninstall
```

## Shell completion

```bash
# Bash
pboard completion bash > /etc/bash_completion.d/pboard

# Zsh
pboard completion zsh > "${fpath[1]}/_pboard"

# Fish
pboard completion fish > ~/.config/fish/completions/pboard.fish
```

## Exit codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General error |
| 2 | Authentication error |
| 3 | Resource not found |
| 4 | Rate limited |
| 5 | Invalid input |
| 6 | Write operations disabled |
