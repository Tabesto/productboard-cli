# Tasks: Feature Health Check Helper

**Input**: Design documents from `/specs/004-feature-health-helper/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: No tests requested. This project has no existing test suite.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Register health subcommand group and create the new file

- [x] T001 Create `internal/cli/features_health.go` with `newFeaturesHealthCmd()` that returns a Cobra command group (`Use: "health"`, `Short: "Feature health updates"`) with no Run function (displays help by default)
- [x] T002 Register health subcommand in `internal/cli/features.go` by adding `cmd.AddCommand(newFeaturesHealthCmd())` inside `newFeaturesCmd()`

**Checkpoint**: `pboard features health --help` displays the help message with no subcommands yet

---

## Phase 2: User Story 1 - List Health Updates Across All Features (Priority: P1) MVP

**Goal**: Fetch all features and display a table of features with their health updates, excluding archived and no-health features by default, sorted by health update date (most recent first)

**Independent Test**: Run `pboard features health list` and verify it shows a table with columns: Feature Name, Status, Owner, Health, Health Updated, Message (truncated to 50 chars). Verify archived features and features without health updates are excluded.

### Implementation for User Story 1

- [x] T003 [US1] Implement `newFeaturesHealthListCmd()` in `internal/cli/features_health.go`: create Cobra command (`Use: "list"`, `Short: "List health updates across features"`), define local flags `--include-archived` (bool) and `--include-no-health` (bool), implement Run function that calls `getClient()` then `c.GetList("/features", params, 0)` (always fetch ALL features for client-side filtering) with `params["archived"] = "false"` when `--include-archived` is not set
- [x] T004 [US1] Add client-side filtering logic in the `health list` Run function in `internal/cli/features_health.go`: iterate over results, skip features where `lastHealthUpdate` is nil (unless `--include-no-health` is set), extract health update fields using `output.SafeNested()`
- [x] T005 [US1] Add sorting logic in the `health list` Run function in `internal/cli/features_health.go`: use `sort.Slice` to sort filtered results by `lastHealthUpdate.updatedAt` descending (most recent first), features without health updates sort to the end
- [x] T006 [US1] Apply `limitFlag` to the filtered+sorted slice (if `limitFlag > 0`, truncate to that length), then render output using `output.Print(outputFormat, finalResults, headers, toRows)` with headers `["Feature Name", "Status", "Owner", "Health", "Health Updated", "Message"]`, truncate message to 50 chars using `output.Truncate()` in `internal/cli/features_health.go`
- [x] T007 [US1] Extract the client-side filter+sort logic into a shared function `filterAndSortHealthFeatures(features []map[string]interface{}, opts HealthFilterOpts) []map[string]interface{}` in `internal/cli/features_health.go`, where `HealthFilterOpts` is a struct holding all filter parameters (includeNoHealth, includeArchived, updatedSince, updatedBefore, healthStatus, statusName, ownerEmail). Refactor T004/T005 logic to call this function. This enables reuse by MCP handlers (D1 remediation).
- [x] T007b [US1] Register `newFeaturesHealthListCmd()` in `newFeaturesHealthCmd()` via `cmd.AddCommand()` in `internal/cli/features_health.go`

**Checkpoint**: `pboard features health list` shows a table of features with health updates. `pboard features health list --include-no-health` also shows features without health updates. `-o json` outputs valid JSON array. MVP is functional.

---

## Phase 3: User Story 2 - Filter Health Updates by Date (Priority: P1)

**Goal**: Add `--updated-since` and `--updated-before` date range filtering to the health list command

**Independent Test**: Run `pboard features health list --updated-since 2025-11-01 --updated-before 2026-01-01` and verify only features with health updates in that range appear. Run with invalid date like `--updated-since bad-date` and verify error message with exit code 5.

### Implementation for User Story 2

- [x] T008 [US2] Add `--updated-since` and `--updated-before` string flags to `newFeaturesHealthListCmd()` in `internal/cli/features_health.go`
- [x] T009 [US2] Implement date parsing and validation in `internal/cli/features_health.go`: parse flags using `time.Parse("2006-01-02", value)`, on parse error call `handleError(&client.APIError{StatusCode: 0, Message: "invalid date format, expected YYYY-MM-DD", ExitCode: client.ExitInvalidInput})` (constitution requires all errors go through APIError type)
- [x] T010 [US2] Add date range filtering to the client-side filter loop in `internal/cli/features_health.go`: parse each feature's `lastHealthUpdate.updatedAt` (ISO 8601), skip features outside the `--updated-since` / `--updated-before` range using `time.Before()` / `time.After()` comparisons

**Checkpoint**: Date filtering works with both flags individually and combined. Invalid dates produce exit code 5 with a clear error.

---

## Phase 4: User Story 3 - Filter Health Updates by Feature Attributes (Priority: P2)

**Goal**: Add `--status`, `--owner`, and `--health-status` filters to the health list command

**Independent Test**: Run `pboard features health list --health-status at-risk` and verify only at-risk features appear. Combine with `--owner` and verify AND logic.

### Implementation for User Story 3

- [x] T011 [US3] Add `--status`, `--owner`, and `--health-status` string flags to `newFeaturesHealthListCmd()` in `internal/cli/features_health.go`
- [x] T012 [US3] Implement server-side delegation for `--status` and `--owner` flags in `internal/cli/features_health.go`: when `--status` is set, add `params["statusName"] = statusFlag`; when `--owner` is set, add `params["ownerEmail"] = ownerFlag` (these are handled by the ProductBoard API)
- [x] T013 [US3] Implement client-side `--health-status` filter in the filter loop in `internal/cli/features_health.go`: compare `lastHealthUpdate.status` case-insensitively against the flag value using `strings.EqualFold()`

**Checkpoint**: All three attribute filters work individually and in combination with date filters. Filters combine with AND logic.

---

## Phase 5: User Story 4 - View Detailed Health Update for a Single Feature (Priority: P2)

**Goal**: Add `pboard features health get <feature-id>` command that shows full health details for one feature

**Independent Test**: Run `pboard features health get <valid-id>` and verify it shows a key-value table with ID, Name, Status, Owner, Health Status, Health Updated, and full Message. Run with invalid ID and verify exit code 3.

### Implementation for User Story 4

- [x] T014 [US4] Implement `newFeaturesHealthGetCmd()` in `internal/cli/features_health.go`: create Cobra command (`Use: "get <feature-id>"`, `Args: cobra.ExactArgs(1)`), implement Run function that calls `c.GetSingle(fmt.Sprintf("/features/%s", args[0]))`
- [x] T015 [US4] Add output rendering in `newFeaturesHealthGetCmd()` Run function in `internal/cli/features_health.go`: for JSON use `output.PrintJSON(feature)`, for table use `output.PrintSingleTable()` with rows for ID, Name, Status, Owner, Health Status, Health Updated, Message; when `lastHealthUpdate` is nil display Health Status as "(none)" and omit Health Updated and Message rows
- [x] T016 [US4] Register `newFeaturesHealthGetCmd()` in `newFeaturesHealthCmd()` via `cmd.AddCommand()` in `internal/cli/features_health.go`

**Checkpoint**: `pboard features health get <id>` shows full health details. Features without health updates show "(none)". Invalid IDs produce exit code 3.

---

## Phase 6: User Story 5 + MCP Tools + Polish (Priority: P3)

**Goal**: Expose health commands as MCP tools and ensure JSON output works correctly across all commands

**Independent Test**: Run `pboard features health list -o json | python3 -c "import json,sys; json.load(sys.stdin)"` to verify valid JSON. Test MCP tools via `pboard mcp serve` with a test client.

### Implementation

- [x] T017 [P] [US5] Verify JSON output for `health list` and `health get` commands by running manual tests with `-o json` flag -- no code changes expected since `output.Print()` and `output.PrintJSON()` handle this automatically
- [x] T018 [P] Register `features_health_list` MCP tool in `internal/mcp/tools.go`: add `s.AddTool(mcp.NewTool("features_health_list", ...))` with string params `updated_since`, `updated_before`, `status`, `owner`, `health_status`, boolean params `include_archived`, `include_no_health`, and number param `limit` (default 25)
- [x] T019 [P] Register `features_health_get` MCP tool in `internal/mcp/tools.go`: add `s.AddTool(mcp.NewTool("features_health_get", ...))` with required string param `id`
- [x] T020 Implement `handleFeaturesHealthList` handler in `internal/mcp/handlers.go`: extract params from request, call `getClient()`, call `c.GetList("/features", params, 0)` (fetch all for client-side filtering), then call the shared `filterAndSortHealthFeatures()` function from `internal/cli/features_health.go` with filter opts built from request params, apply `getLimit(request)` to truncate the result, return `toJSON(filtered)`
- [x] T021 Implement `handleFeaturesHealthGet` handler in `internal/mcp/handlers.go`: extract `id` param, call `getClient()`, call `c.GetSingle(fmt.Sprintf("/features/%s", id))`, return `toJSON(feature)`
- [x] T022 Build and verify: run `go build ./cmd/pboard` and manually test all commands per quickstart.md scenarios

**Checkpoint**: All CLI commands and MCP tools functional. JSON output valid. Build succeeds with no errors.

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **US1 (Phase 2)**: Depends on Setup (T001, T002) - this is the MVP
- **US2 (Phase 3)**: Depends on US1 completion (adds flags to the list command created in US1)
- **US3 (Phase 4)**: Depends on US1 completion (adds more flags to the list command)
- **US4 (Phase 5)**: Depends on Setup only (creates an independent command) - can run in parallel with US2/US3
- **US5 + MCP + Polish (Phase 6)**: Depends on US1-US4 completion

### User Story Dependencies

- **User Story 1 (P1)**: Depends on Setup (Phase 1) only
- **User Story 2 (P1)**: Depends on US1 (adds date filtering to the list command)
- **User Story 3 (P2)**: Depends on US1 (adds attribute filtering to the list command)
- **User Story 4 (P2)**: Depends on Setup only -- **can run in parallel with US2 and US3**
- **User Story 5 (P3)**: Satisfied by the output.Print() pattern used in US1 and US4 -- verification only

### Within Each User Story

- Flag definitions before filter logic
- Filter logic before output rendering
- Registration as final step

### Parallel Opportunities

- US2 (date filters) and US3 (attribute filters) can run in parallel after US1 completes (they add different flags to the same command but touch different code sections)
- US4 (health get) can run in parallel with US2 and US3 (separate command entirely)
- T018 and T019 (MCP tool registration) can run in parallel
- T017, T018, T019 can all run in parallel

---

## Parallel Example: After US1 Completes

```bash
# These can run in parallel (different concerns, minimal overlap):
Task: "T008 [US2] Add date filtering flags in internal/cli/features_health.go"
Task: "T011 [US3] Add attribute filtering flags in internal/cli/features_health.go"
Task: "T014 [US4] Implement health get command in internal/cli/features_health.go"
```

## Parallel Example: MCP Tools

```bash
# These can run in parallel (different files):
Task: "T018 Register features_health_list MCP tool in internal/mcp/tools.go"
Task: "T019 Register features_health_get MCP tool in internal/mcp/tools.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (T001-T002)
2. Complete Phase 2: US1 (T003-T007)
3. **STOP and VALIDATE**: Run `pboard features health list` -- should display health table
4. This alone delivers the core value: consolidated health overview in one command

### Incremental Delivery

1. Setup + US1 → Health list working (MVP!)
2. Add US2 → Date filtering works → Can find stale health updates
3. Add US3 → Attribute filtering works → Can focus on specific segments
4. Add US4 → Health get works → Can drill into single feature
5. Add US5 + MCP → Full feature complete with automation support

### Single Developer (Recommended)

Since all CLI tasks modify `internal/cli/features_health.go`, sequential execution is cleanest:

1. T001-T002 (Setup)
2. T003-T007b (US1 - MVP)
3. T008-T010 (US2 - date filters)
4. T011-T013 (US3 - attribute filters)
5. T014-T016 (US4 - health get)
6. T017-T022 (US5 + MCP + polish)

---

## Notes

- All CLI commands modify the same file (`internal/cli/features_health.go`) so sequential execution avoids conflicts
- MCP handlers go in a different file (`internal/mcp/handlers.go`) so Phase 6 has true parallelism opportunities
- US5 (JSON output) requires no code changes -- it's verified for free via `output.Print()` pattern
- The `--limit` global flag is applied AFTER client-side filtering: fetch all features (limit=0 to API), filter+sort, then truncate to limitFlag. MCP uses getLimit() (default 25) at the same post-filter stage.
- Server-side filters (`ownerEmail`, `statusName`, `archived`) reduce API payload; client-side filters (`health-status`, date range) handle what the API cannot
