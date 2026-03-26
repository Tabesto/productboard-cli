# Implementation Plan: Feature Health Check Helper

**Branch**: `004-feature-health-helper` | **Date**: 2026-03-26 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/004-feature-health-helper/spec.md`

## Summary

Add a `pboard features health` subcommand group that aggregates feature data with their `lastHealthUpdate` field to provide a consolidated health overview. The implementation adds two CLI commands (`health list` and `health get`) and two matching MCP tools, following the existing codebase patterns exactly. Health data is fetched via a single paginated API call to `/features` (which includes `lastHealthUpdate`), with client-side filtering for date ranges, status, owner, and health status.

## Technical Context

**Language/Version**: Go 1.26.1
**Primary Dependencies**: Cobra (CLI), go-pretty (tables), mcp-go (MCP server), Viper (config) -- all existing, no new dependencies
**Storage**: N/A (reads from ProductBoard API via existing client)
**Testing**: No existing test suite in the project; manual testing via CLI
**Target Platform**: macOS/Linux CLI (single binary)
**Project Type**: CLI tool
**Performance Goals**: Complete health list for up to 200 features within a few seconds (bounded by API pagination latency)
**Constraints**: Read-only API access; client-side filtering (API does not support health-related query parameters); all features fetched then filtered in memory
**Scale/Scope**: Typical boards have ~100-200 features; client-side sort/filter is appropriate at this scale

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| I. Read-Only API Wrapper | PASS | Only HTTP GET to `/features` endpoint |
| II. CLI-First Interface | PASS | CLI commands implemented first; MCP tools reuse `internal/client` |
| III. Single Binary Distribution | PASS | New subcommands within existing `pboard` binary |
| IV. Credential Safety | PASS | No token handling changes; uses existing config/env mechanism |
| Constraints: No new dependencies | PASS | All dependencies already in go.mod |
| Constraints: X-Version: 1 header | PASS | Handled by existing client |
| Constraints: MCP default limit 25 | PASS | MCP tools will use `getLimit()` with default 25 |
| Dev Workflow: APIError exit codes | PASS | Uses existing `handleError()` pattern |

All gates pass. No violations to justify.

## Project Structure

### Documentation (this feature)

```text
specs/004-feature-health-helper/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output (CLI and MCP contracts)
└── tasks.md             # Phase 2 output (/speckit.tasks command)
```

### Source Code (repository root)

```text
internal/
├── cli/
│   ├── features.go              # Existing: add newFeaturesHealthCmd() registration
│   └── features_health.go       # NEW: health list + health get commands + client-side filtering
├── client/
│   └── client.go                # Existing: no changes needed
├── mcp/
│   ├── tools.go                 # Existing: register 2 new MCP tools
│   └── handlers.go              # Existing: add 2 new handler functions
└── output/
    └── formatter.go             # Existing: no changes needed
```

**Structure Decision**: Extend existing `internal/cli/` with a new `features_health.go` file (follows the pattern of keeping related commands grouped, similar to how `features.go` handles list/get/links). MCP tools and handlers are added to the existing `tools.go` and `handlers.go` files since all tools are registered in a single function.

## Complexity Tracking

No violations. This section is intentionally empty.
