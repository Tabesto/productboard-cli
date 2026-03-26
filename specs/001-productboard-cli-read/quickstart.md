# Quickstart: pboard CLI

## Installation

```bash
brew tap <org>/homebrew-tap
brew install pboard
```

## Configuration

Set your ProductBoard API token:

```bash
# Option 1: Interactive setup (writes to ~/.config/pboard/config.yaml)
pboard configure

# Option 2: Environment variable
export PRODUCTBOARD_API_TOKEN="pb_your_token_here"
```

## Basic Usage

### List features
```bash
pboard features list
pboard features list --status-name "In Progress"
pboard features list --owner-email user@example.com --output json
pboard features list --limit 10
```

### Get a specific feature
```bash
pboard features get abc123-def456
```

### View linked resources
```bash
pboard features links initiatives abc123-def456
pboard objectives links features abc123-def456
```

### Browse product hierarchy
```bash
pboard products list
pboard components list
pboard feature-statuses list
```

### Search notes
```bash
pboard notes list --term "onboarding" --created-from 2026-01-01
pboard notes get abc123
pboard notes tags abc123
```

### View releases
```bash
pboard releases list
pboard release-groups list
pboard feature-release-assignments list --release-id abc123
```

### View OKRs
```bash
pboard objectives list
pboard key-results list
pboard initiatives list
```

### JSON output for scripting
```bash
# Pipe to jq for advanced filtering
pboard features list -o json | jq '.[] | select(.status.name == "Done")'

# Export to file
pboard notes list -o json > notes-export.json
```

## Help

```bash
pboard --help
pboard features --help
pboard features list --help
```
