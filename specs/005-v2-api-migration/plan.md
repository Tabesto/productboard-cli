# Implementation Plan: Migrate CLI to ProductBoard V2 API

**Branch**: `005-v2-api-migration` | **Date**: 2026-04-01 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/005-v2-api-migration/spec.md`

## Summary

Migrate the pboard CLI from ProductBoard V1 API (using `X-Version: 1` header with per-type endpoints) to V2 API (using `/v2` URL prefix with unified `/entities` endpoint). Dual V1/V2 support with V2 as default. Add new `members` and `teams` commands (V2 only). CLI remains read-only.

## Technical Context

**Language/Version**: Go 1.26.1 (go.mod)
**Primary Dependencies**: Cobra (CLI), Viper (config), go-pretty/v6 (tables), mcp-go (MCP server) — no new dependencies
**Storage**: Config file at `~/.config/pboard/config.yaml` (mode 600)
**Testing**: Manual verification against live API (no test framework in project)
**Target Platform**: macOS/Linux (single binary via Homebrew/go install)
**Project Type**: CLI tool
**Performance Goals**: N/A (interactive CLI, network-bound)
**Constraints**: Read-only API access, single binary, no new dependencies
**Scale/Scope**: ~20 CLI commands, ~50 MCP tool handlers

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| I. Read-Only API Wrapper | PASS | Migration preserves read-only constraint. No write operations added. |
| II. CLI-First Interface | PASS | New members/teams commands follow CLI-first pattern. MCP reuses same client. |
| III. Single Binary Distribution | PASS | No new binaries or processes. |
| IV. Credential Safety | PASS | Token handling unchanged. |
| **Constraint: API versioning** | **VIOLATION** | Constitution says "All ProductBoard API requests MUST include the `X-Version: 1` header". V2 migration requires amending this. |
| Constraint: Default result limit | PASS | MCP tools retain 25-item default. |
| Development Workflow | PASS | Following numbered branch pattern. |

### Post-Design Re-check

All principles pass. The API versioning constraint violation is justified and requires a constitution amendment (see Complexity Tracking).

## Project Structure

### Documentation (this feature)

```text
specs/005-v2-api-migration/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/
│   └── cli-commands.md  # Phase 1 output
└── tasks.md             # Phase 2 output (via /speckit.tasks)
```

### Source Code (repository root)

```text
cmd/pboard/main.go                    # Entry point (unchanged)
internal/
├── cli/
│   ├── root.go                       # MODIFY: add --api-version flag, register new commands
│   ├── features.go                   # MODIFY: V2 path/param mapping (handled by client)
│   ├── notes.go                      # MODIFY: V2 path/param mapping (handled by client)
│   ├── companies.go                  # MODIFY: V2 path/param mapping (handled by client)
│   ├── users.go                      # MODIFY: V2 path/param mapping (handled by client)
│   ├── products.go                   # MODIFY: V2 path/param mapping (handled by client)
│   ├── components.go                 # MODIFY: V2 path/param mapping (handled by client)
│   ├── initiatives.go               # MODIFY: V2 path/param mapping (handled by client)
│   ├── objectives.go                # MODIFY: V2 path/param mapping (handled by client)
│   ├── key_results.go               # MODIFY: V2 path/param mapping (handled by client)
│   ├── releases.go                  # MODIFY: V2 path/param mapping (handled by client)
│   ├── release_groups.go            # MODIFY: V2 path/param mapping (handled by client)
│   ├── assignments.go               # MODIFY: V2 path/param mapping (handled by client)
│   ├── custom_fields.go             # MODIFY: V2 path/param mapping (handled by client)
│   ├── feature_statuses.go          # MODIFY: V2 path/param mapping (handled by client)
│   ├── jira_integrations.go         # MODIFY: V2 path/param mapping (handled by client)
│   ├── plugin_integrations.go       # MODIFY: V2 path/param mapping (handled by client)
│   ├── webhooks.go                  # MODIFY: V2 path/param mapping (handled by client)
│   ├── feedback_forms.go            # MODIFY: add V2 unavailability guard
│   ├── members.go                   # NEW: V2-only members commands
│   ├── teams.go                     # NEW: V2-only teams commands
│   ├── configure.go                 # MODIFY: save api_version to config
│   └── mcp.go                       # MODIFY: pass api_version through
├── client/
│   ├── client.go                    # MODIFY: core V2 support (path mapping, response flattening, pagination)
│   ├── errors.go                    # MODIFY: V2 error parsing
│   └── v2.go                        # NEW: V2 path/param translation logic
├── config/
│   └── config.go                    # MODIFY: add APIVersion field, env var, config key
├── mcp/
│   ├── server.go                    # MODIFY: pass api_version to client
│   ├── handlers.go                  # MODIFY: adapt to flattened V2 responses if needed
│   └── tools.go                     # MODIFY: add members/teams tools
├── models/
│   └── models.go                    # MODIFY: add Member, Team structs (documentation only)
└── output/
    └── formatter.go                 # No changes (SafeStr/SafeNested work with flattened data)
```

**Structure Decision**: No new packages or directories. V2 translation logic lives in a new `internal/client/v2.go` file to keep `client.go` manageable. All other changes are modifications to existing files.

## Complexity Tracking

> **Constitution violation justification**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| API versioning constraint (`X-Version: 1` header) | V2 API uses URL prefix instead of header. V1 API will be deprecated. | Staying on V1 only would block access to new V2 features (members, teams, unified entities) and eventually break when V1 is sunset. |

**Required amendment**: Update constitution constraint from "All ProductBoard API requests MUST include the `X-Version: 1` header" to "API version is configurable via `--api-version` flag or config (default V2). V1 uses `X-Version: 1` header; V2 uses `/v2` URL prefix." Bump constitution version 1.0.0 → 1.1.0.
