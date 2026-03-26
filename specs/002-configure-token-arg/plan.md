# Implementation Plan: Configure Command Token Argument

**Branch**: `002-configure-token-arg` | **Date**: 2026-03-26 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/002-configure-token-arg/spec.md`

## Summary

Enhance the existing `pboard configure` command to accept an optional positional argument for the API token. When provided, the token is saved directly without interactive prompts. When omitted, the existing interactive prompt behavior is preserved. This is a minimal change to a single file (`internal/cli/configure.go`).

## Technical Context

**Language/Version**: Go 1.22+ (existing project)
**Primary Dependencies**: Cobra (existing), no new dependencies
**Storage**: Existing config file at `~/.config/pboard/config.yaml` (mode 600)
**Testing**: Go standard `testing` package
**Target Platform**: macOS (arm64, amd64), Linux (amd64) — unchanged
**Project Type**: CLI tool (existing)
**Performance Goals**: N/A — instant operation
**Constraints**: Backward-compatible with existing `pboard configure` usage
**Scale/Scope**: Single file change (~20 lines modified)

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

Constitution is an uninitialized template — no project-specific gates defined. **PASS** (no constraints to violate).

**Post-Phase 1 re-check**: Change is a single-file modification to an existing command. **PASS**.

## Project Structure

### Documentation (this feature)

```text
specs/002-configure-token-arg/
├── plan.md              # This file
├── research.md          # Phase 0 output (minimal — no unknowns)
├── contracts/           # Phase 1 output — updated CLI contract
│   └── cli-commands.md
└── tasks.md             # Phase 2 output (/speckit.tasks command)
```

### Source Code (repository root)

```text
internal/
└── cli/
    └── configure.go     # Only file modified
```

**Structure Decision**: No structural changes — this is a behavioral modification to a single existing file.

## Complexity Tracking

> No constitution violations — section not applicable.
