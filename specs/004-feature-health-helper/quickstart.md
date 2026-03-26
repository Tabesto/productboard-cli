# Quickstart: Feature Health Check Helper

## Prerequisites

- Go 1.22+ installed
- `pboard` configured with a valid ProductBoard API token (`pboard configure`)

## Build & Run

```bash
# Build from repository root
go build -o pboard ./cmd/pboard

# Verify the new commands are available
./pboard features health --help
./pboard features health list --help
./pboard features health get --help
```

## Usage Examples

### List all feature health updates

```bash
# Default: non-archived features with health updates, sorted by most recent
pboard features health list

# JSON output for scripting
pboard features health list -o json
```

### Filter by date

```bash
# Features with health updated in the last 30 days
pboard features health list --updated-since 2026-02-24

# Features with stale health updates (before a cutoff date)
pboard features health list --updated-before 2026-01-01

# Date range
pboard features health list --updated-since 2026-01-01 --updated-before 2026-03-01
```

### Filter by attributes

```bash
# Only at-risk features
pboard features health list --health-status at-risk

# Features owned by a specific person
pboard features health list --owner amel.meradi@deliverect.com

# Combine filters
pboard features health list --status "In Progress" --health-status on-track --updated-since 2026-01-01
```

### Include normally-hidden features

```bash
# Include archived features
pboard features health list --include-archived

# Include features without health updates
pboard features health list --include-no-health
```

### Get detailed health for a specific feature

```bash
# By feature ID
pboard features health get 7043709b-f020-43f0-bebb-8ae7649acbf7

# JSON output
pboard features health get 7043709b-f020-43f0-bebb-8ae7649acbf7 -o json
```

## Files Changed

| File | Change |
|------|--------|
| `internal/cli/features.go` | Register `newFeaturesHealthCmd()` |
| `internal/cli/features_health.go` | **NEW** -- health list + get commands with client-side filtering |
| `internal/mcp/tools.go` | Register `features_health_list` and `features_health_get` tools |
| `internal/mcp/handlers.go` | Add `handleFeaturesHealthList` and `handleFeaturesHealthGet` handlers |
