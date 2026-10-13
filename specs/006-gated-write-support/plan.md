# Implementation Plan: Gated Write Support

**Branch**: `006-gated-write-support` | **Date**: 2026-10-13 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/006-gated-write-support/spec.md`

## Summary

Add create/update of entities and notes through Productboard API V2, off by default. A single boolean "write gate" (flag `--allow-write` OR config `write_enabled` OR env `PRODUCTBOARD_WRITE_ENABLED`) is enforced inside `internal/client` for every non-GET request, so CLI and MCP share one enforcement point. MCP write tools are registered only when the gate is on at server start. `configure` is fixed to preserve existing settings and gains `--write-enabled`. Constitution principle I is amended (read-only by default, gated writes).

## Technical Context

**Language/Version**: Go 1.26.1 (go.mod)
**Primary Dependencies**: Cobra, Viper, go-pretty/v6, mcp-go v0.46.0 — all existing, no new dependencies
**Storage**: `~/.config/pboard/config.yaml` (mode 600), new key `write_enabled`
**Testing**: `go test` with `net/http/httptest` (no tests exist today; first test files added here)
**Target Platform**: macOS/Linux CLI + stdio MCP server
**Project Type**: single Go module, CLI + embedded MCP server
**Performance Goals**: N/A (interactive writes; Productboard limit 50 req/s per token)
**Constraints**: writes V2 only; no retry; refuse before any network call when gate is off
**Scale/Scope**: 4 write operations (entity create/update, note create/update), 4 MCP tools

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| I. Read-Only API Wrapper | **VIOLATION (justified)** | Feature deliberately adds writes. Resolved by constitution amendment 1.1.0 → 2.0.0 ("Read-Only by Default, Gated Writes"), delivered in this feature. See Complexity Tracking. |
| II. CLI-First Interface | PASS | Writes are Cobra subcommands with table/JSON output; MCP reuses `internal/client`, no second HTTP client. |
| III. Single Binary | PASS | No new binaries; `--allow-write` is a flag on existing `mcp serve`. |
| IV. Credential Safety | PASS | Token handling unchanged; config rewrite keeps mode 600 and never prints the token. |
| Constraints: dependencies | PASS | None added (Viper already writes YAML). |
| Workflow: error scheme | AMEND | Adds exit code 6 (`write disabled`) to the documented scheme. |

**Post-design re-check**: same result. Gate is enforced centrally, so the amended principle ("no write path may bypass the gate") is verifiable by test.

## Project Structure

### Documentation (this feature)

```text
specs/006-gated-write-support/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── interfaces.md    # CLI commands + MCP tools + HTTP requests
├── checklists/requirements.md
└── tasks.md             # /speckit.tasks
```

### Source Code (repository root)

```text
internal/
├── config/
│   ├── config.go        # + WriteEnabled, env parsing, Update() merge-writer (replaces WriteConfig)
│   └── config_test.go   # new
├── client/
│   ├── client.go        # Get -> do(); + Post/Patch/RequireWrite; http timeout
│   ├── errors.go        # + ExitWriteDisabled = 6, NewWriteDisabledError
│   ├── write.go         # new: WriteRef parsing, body helpers
│   └── client_test.go   # new (httptest)
├── cli/
│   ├── root.go          # + --allow-write persistent flag; getClient applies it; help text
│   ├── configure.go     # + --write-enabled; preserve settings; token optional when flag set
│   ├── write.go         # new: --field parsing, shared write helpers, print result
│   ├── entities.go      # new: `entities create|update`
│   ├── notes.go         # + `notes create|update`
│   ├── mcp.go           # + --allow-write on serve and install
│   ├── skill_content.go # docs for write + gate
│   └── write_test.go    # new
└── mcp/
    ├── server.go        # Serve(allowWrite bool)
    ├── tools.go         # registerTools(s, allowWrite) -> registerWriteTools(s)
    ├── handlers.go      # + 4 write handlers (or handlers_write.go)
    └── tools_test.go    # new: tool list with/without gate
README.md                # document gate, write commands, API V2
.specify/memory/constitution.md  # amendment 2.0.0
```

**Structure Decision**: extend existing packages; only new files are `client/write.go`, `cli/write.go`, `cli/entities.go` and test files. No new packages.

## Complexity Tracking

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| Principle I (read-only) amended to gated writes | User requirement: create/update entities and notes from CLI and MCP | A separate write-only binary violates Principle III (single binary); leaving writes out defeats the feature. Gate keeps default behavior identical to today. |
