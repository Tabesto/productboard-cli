# Implementation Plan: MCP Server for pboard CLI

**Branch**: `003-pboard-mcp-server` | **Date**: 2026-03-26 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/003-pboard-mcp-server/spec.md`

## Summary

Build an MCP (Model Context Protocol) server embedded in the existing pboard binary that exposes all ~37 pboard CLI read commands as individual MCP tools accessible from Claude Desktop. The server communicates via stdio transport and is installed with a single `pboard mcp install` command that auto-configures Claude Desktop. Uses `mark3labs/mcp-go` library and reuses the existing `internal/client` for all ProductBoard API calls.

## Technical Context

**Language/Version**: Go 1.22+ (existing project, go.mod shows 1.26.1)
**Primary Dependencies**: Cobra (CLI), mcp-go (MCP server), viper (config), go-pretty (output)
**Storage**: N/A (stateless server; reads config from `~/.config/pboard/config.yaml`)
**Testing**: `go test` with table-driven tests
**Target Platform**: macOS (Claude Desktop config path; binary is cross-platform)
**Project Type**: CLI tool with embedded MCP server
**Performance Goals**: Same response time as direct CLI commands
**Constraints**: Stdio transport (no networking), default 25 results per list query, read-only
**Scale/Scope**: ~37 MCP tools, single-user local server

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

Constitution is an unfilled template -- no gates defined. **PASS** (no violations possible).

**Post-Phase 1 re-check**: Still PASS. Design adds one new dependency (`mcp-go`) and one new package (`internal/mcp`), consistent with existing project patterns.

## Project Structure

### Documentation (this feature)

```text
specs/003-pboard-mcp-server/
├── plan.md              # This file
├── spec.md              # Feature specification
├── research.md          # Phase 0: Technology decisions
├── data-model.md        # Phase 1: Tool registry & config model
├── quickstart.md        # Phase 1: Dev setup guide
├── contracts/
│   └── mcp-tools.md     # Phase 1: MCP tool & error contracts
├── checklists/
│   └── requirements.md  # Spec quality checklist
└── tasks.md             # Phase 2 output (/speckit.tasks command)
```

### Source Code (repository root)

```text
cmd/pboard/main.go              # Entry point (unchanged)
internal/
├── cli/
│   ├── root.go                  # Modify: register mcp command
│   ├── mcp.go                   # Create: pboard mcp (install/uninstall/serve)
│   └── [existing commands]      # Unchanged
├── mcp/
│   ├── server.go                # Create: MCP server setup & stdio entry
│   ├── tools.go                 # Create: Tool definitions (~37 tools)
│   └── handlers.go              # Create: Tool handlers calling client
├── client/                      # Unchanged (reused by handlers)
├── config/                      # Unchanged (reused for token loading)
└── output/                      # Unchanged
go.mod                           # Modify: add mcp-go dependency
```

**Structure Decision**: Follows existing project layout. New `internal/mcp/` package encapsulates all MCP server logic. New `internal/cli/mcp.go` adds the `pboard mcp` command group following the same pattern as `skill.go`. No new top-level directories needed.

## Complexity Tracking

No constitution violations to justify.
