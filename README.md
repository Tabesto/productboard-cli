# pboard - ProductBoard CLI

A read-only command-line tool for browsing ProductBoard data. Supports all GET endpoints from the ProductBoard API.

## Installation

### Homebrew (macOS / Linux)

```bash
brew tap tabesto/tap
brew install pboard
```

### From source

```bash
go install github.com/tabesto/productboard-cli/cmd/pboard@latest
```

## Configuration

Set your ProductBoard API token using one of these methods:

### Interactive setup

```bash
pboard configure
```

This saves the token to `~/.config/pboard/config.yaml` (mode 600).

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
| `--help` | `-h` | Show help |
| `--version` | `-v` | Show version |

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
