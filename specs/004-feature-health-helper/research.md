# Research: Feature Health Check Helper

**Date**: 2026-03-26 | **Branch**: `004-feature-health-helper`

## R1: lastHealthUpdate Field Structure

**Decision**: The `lastHealthUpdate` field is available in both the features list (`/features`) and single feature GET (`/features/{id}`) API responses. It contains:
- `status`: string enum (e.g., "on-track", "at-risk", "off-track")
- `updatedAt`: ISO 8601 timestamp
- `message`: single free-text field (authors may follow progress/problems/plans convention)

**Rationale**: Confirmed during spec clarification phase. The field is included in list results, so a single paginated call to `/features` with `limit=0` retrieves all features with their health data. No N+1 API call problem.

**Alternatives considered**:
- Individual GET per feature: Rejected -- would require N+1 calls, slow for boards with 100+ features
- Separate health endpoint: Does not exist in ProductBoard API

## R2: Client-Side Filtering Strategy

**Decision**: Fetch all features via `GetList("/features", params, 0)` then filter in-memory by:
1. Health update date range (`--updated-since`, `--updated-before`)
2. Health status (`--health-status`)
3. Feature-level filters delegated to API where possible (`ownerEmail`, `statusName`, `archived`)

**Rationale**: The ProductBoard API supports filtering by `ownerEmail`, `statusName`, and `archived` as query parameters (server-side). However, it does not support filtering by `lastHealthUpdate` fields. The hybrid approach minimizes data transfer by using server-side filters where available, then applying health-specific filters client-side.

**Alternatives considered**:
- Pure client-side filtering (fetch everything, filter locally): Simpler but wastes bandwidth when owner/status filters are used
- Pure server-side filtering: Not possible -- API lacks health-related filter parameters

## R3: Command Placement

**Decision**: Add `health` as a subcommand group under `features`, creating `pboard features health list` and `pboard features health get <id>`.

**Rationale**: Follows the existing nesting pattern (e.g., `features links initiatives`). Health data is a property of features, so it belongs under the features command tree. This also establishes the "helper" pattern as nested subcommands rather than top-level commands.

**Alternatives considered**:
- Top-level `pboard health` command: Would break the convention where every command maps to a ProductBoard entity
- Adding health columns to existing `features list`: Would make the default table too wide; better as a focused view

## R4: New File vs Extending Existing

**Decision**: Create a new file `internal/cli/features_health.go` for the health commands rather than extending `features.go`.

**Rationale**: `features.go` is already 220 lines with list/get/links commands. Adding ~150 lines of health logic (list with filtering + get with detail view) would make it harder to navigate. A separate file keeps concerns isolated while the registration still happens in `newFeaturesCmd()` via `cmd.AddCommand(newFeaturesHealthCmd())`.

**Alternatives considered**:
- Inline in `features.go`: Would work but results in a 370+ line file
- Separate package: Overkill for two commands that share the same client/output patterns

## R5: No New Dependencies Required

**Decision**: No new Go dependencies needed.

**Rationale**: Date parsing uses `time.Parse` from stdlib. Sorting uses `sort.Slice` from stdlib. All other needs (HTTP client, CLI framework, table output, JSON output, MCP server) are covered by existing dependencies.

**Alternatives considered**: None -- stdlib covers all new requirements.
