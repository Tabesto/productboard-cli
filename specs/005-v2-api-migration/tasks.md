# Tasks: Migrate CLI to ProductBoard V2 API

**Input**: Design documents from `/specs/005-v2-api-migration/`
**Prerequisites**: plan.md (required), spec.md (required), research.md, data-model.md, contracts/

**Tests**: No test tasks included (no test framework in project; manual verification via quickstart.md).

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Phase 1: Setup

**Purpose**: Add API version configuration support across config, client, and CLI flag

- [x] T001 Add `APIVersion` field to Config struct and read `api_version` from config file and `PRODUCTBOARD_API_VERSION` env var in `internal/config/config.go`
- [x] T002 Add `apiVersion` field to Client struct and accept it from Config in `internal/client/client.go`
- [x] T003 Add `--api-version` persistent flag (default "2") to root command in `internal/cli/root.go`, pass it through to config/client creation
- [x] T004 Update `WriteToken` to also persist `api_version` in config file in `internal/config/config.go`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core V2 client infrastructure that MUST be complete before ANY user story can be implemented

**CRITICAL**: No user story work can begin until this phase is complete

- [x] T005 [P] Create V2 path translation module with entity-type path mapping table (e.g., `/features` → `/v2/entities?type[]=feature`) in new file `internal/client/v2.go`
- [x] T006 [P] Add V2 query parameter mapping (e.g., `statusName` → `status[name]`, `ownerEmail` → `owner[email]`) to `internal/client/v2.go`
- [x] T007 Update `Client.Get()` to branch on `apiVersion`: V1 keeps `X-Version: 1` header + direct path; V2 uses translated path via v2.go, no version header, in `internal/client/client.go`
- [x] T008 Implement V2 response flattening in `Client.GetSingle()` — detect `data.fields` sub-object and merge fields to top level for display compatibility in `internal/client/client.go`
- [x] T009 Implement V2 response flattening in `Client.GetList()` — flatten each item's `fields` to top level, extract pagination cursor from `links.next` URL instead of root `pageCursor` in `internal/client/client.go`
- [x] T010 Update V2 error parsing in `NewAPIError()` to attempt parsing `errors[0].detail` from V2 error JSON, fall back to existing behavior in `internal/client/errors.go`

**Checkpoint**: V2 client foundation ready — entity browsing commands can now be migrated

---

## Phase 3: User Story 1 - Unified Entity Browsing via V2 API (Priority: P1) MVP

**Goal**: All existing list/get commands work against V2 API with identical output

**Independent Test**: Run `pboard features list`, `pboard features get <id>`, `pboard notes list`, `pboard products list` etc. and verify equivalent output to V1. Test with `--api-version 1` to confirm V1 fallback works.

### Implementation for User Story 1

- [x] T011 [P] [US1] Verify `pboard features list` and `pboard features get <id>` work with V2 entity path mapping — adjust field extraction in `internal/cli/features.go` if V2 response field names differ
- [x] T012 [P] [US1] Verify `pboard products list` and `pboard products get <id>` work with V2 — adjust if needed in `internal/cli/products.go`
- [x] T013 [P] [US1] Verify `pboard components list` and `pboard components get <id>` work with V2 — adjust if needed in `internal/cli/components.go`
- [x] T014 [P] [US1] Verify `pboard initiatives list/get`, `pboard objectives list/get`, `pboard key-results list/get` work with V2 — adjust if needed in `internal/cli/initiatives.go`, `internal/cli/objectives.go`, `internal/cli/key_results.go`
- [x] T015 [P] [US1] Verify `pboard releases list/get` and `pboard release-groups list/get` work with V2 — adjust if needed in `internal/cli/releases.go`, `internal/cli/release_groups.go`
- [x] T016 [P] [US1] Verify `pboard companies list/get` and `pboard users list/get` work with V2 — adjust if needed in `internal/cli/companies.go`, `internal/cli/users.go`
- [x] T017 [US1] Verify `pboard notes list` and `pboard notes get <id>` work with V2 `/notes` endpoint — notes use a separate V2 path, not unified entities — adjust path mapping and field extraction in `internal/cli/notes.go`
- [x] T018 [P] [US1] Verify `pboard notes tags <id>` works with V2 — map V1 `/notes/{id}/tags` to V2 equivalent (note relationships or configuration) or add deprecation guard if no V2 equivalent exists, in `internal/client/v2.go` and `internal/cli/notes.go`
- [x] T019 [P] [US1] Verify `pboard features health` works with V2 flattened entity data — health command aggregates feature data internally, confirm it handles V2 response shape correctly in `internal/cli/features_health.go`
- [x] T020 [US1] Test pagination with `--limit` flag across multiple entity types to confirm V2 cursor extraction works end-to-end
- [x] T021 [US1] Test `--api-version 1` fallback for all commands to confirm V1 still works

**Checkpoint**: All existing list/get commands work with V2 (default) and V1 (fallback). MVP is functional.

---

## Phase 4: User Story 2 - Relationship and Linked Resource Browsing (Priority: P2)

**Goal**: All `links` subcommands fetch from V2 relationships endpoint

**Independent Test**: Run `pboard features links initiatives <id>`, `pboard objectives links features <id>`, `pboard notes links <id>` and verify linked entities display correctly.

### Implementation for User Story 2

- [x] T022 [US2] Add V2 relationship path translation — map V1 link paths (e.g., `/features/{id}/links/initiatives`) to V2 (`/v2/entities/{id}/relationships?target[type]=initiative`) in `internal/client/v2.go`
- [x] T023 [US2] Update `Client.GetLinkedResources()` to handle V2 relationship response shape — V2 returns `source`/`target` objects instead of flat entity data — extract and flatten target entity fields in `internal/client/client.go`
- [x] T024 [P] [US2] Verify `pboard features links initiatives <id>` and `pboard features links objectives <id>` work with V2 in `internal/cli/features.go`
- [x] T025 [P] [US2] Verify `pboard objectives links features/initiatives <id>` and `pboard initiatives links features/objectives <id>` work with V2 in `internal/cli/objectives.go`, `internal/cli/initiatives.go`
- [x] T026 [US2] Verify `pboard notes links <id>` works with V2 `/notes/{id}/relationships` in `internal/cli/notes.go`
- [x] T027 [US2] Add V2 path mapping for `pboard assignments list/get` — V1 `/feature-release-assignments` maps to V2 entity relationships — update path translation and response handling in `internal/client/v2.go` and `internal/cli/assignments.go`

**Checkpoint**: All relationship/link commands work with V2

---

## Phase 5: User Story 3 - Custom Field and Status Browsing (Priority: P2)

**Goal**: Custom fields, field values, and feature statuses fetched from V2 configuration endpoints

**Independent Test**: Run `pboard custom-fields list`, `pboard custom-fields values list`, `pboard feature-statuses list` and verify output.

### Implementation for User Story 3

- [x] T030 [US3] Add V2 path mapping for custom fields — V1 `/hierarchy-entities/custom-fields` maps to V2 `/v2/entities/configurations` — add translation and response extraction in `internal/client/v2.go`
- [x] T031 [US3] Add V2 path mapping for custom field values — V1 `/hierarchy-entities/custom-fields-values` maps to V2 `/v2/entities/fields/{id}/values` — update translation in `internal/client/v2.go`
- [x] T030 [US3] Update `pboard custom-fields list/get` and `pboard custom-fields values list/get` to handle V2 configuration response shape in `internal/cli/custom_fields.go`
- [x] T031 [US3] Add V2 path mapping for feature statuses — V1 `/feature-statuses` maps to V2 `/v2/entities/configurations/feature` (extract status field from config) — update in `internal/client/v2.go` and `internal/cli/feature_statuses.go`

**Checkpoint**: Custom fields and statuses work with V2

---

## Phase 6: User Story 4 - Integration and Webhook Browsing (Priority: P3)

**Goal**: Jira integrations, plugin integrations, and webhooks fetched from V2 endpoints

**Independent Test**: Run `pboard jira-integrations list`, `pboard plugin-integrations list`, `pboard webhooks list` and verify output.

### Implementation for User Story 4

- [x] T034 [P] [US4] Add V2 path mapping for Jira integrations — prefix-only change `/jira-integrations` → `/v2/jira-integrations` — verify list/get/connections in `internal/client/v2.go` and `internal/cli/jira_integrations.go`
- [x] T033 [P] [US4] Add V2 path mapping for plugin integrations — prefix-only change `/plugin-integrations` → `/v2/plugin-integrations` — verify list/get/connections in `internal/client/v2.go` and `internal/cli/plugin_integrations.go`
- [x] T034 [P] [US4] Add V2 path mapping for webhooks — prefix-only change `/webhooks` → `/v2/webhooks` — verify list/get in `internal/client/v2.go` and `internal/cli/webhooks.go`
- [x] T035 [US4] Add V2 unavailability guard to `pboard feedback-forms` — when `apiVersion=2`, display error "The 'feedback-forms' command is not available with API V2. Use --api-version 1 to access this command." in `internal/cli/feedback_forms.go`

**Checkpoint**: All integration and webhook commands work with V2

---

## Phase 7: User Story 5 - MCP Server Uses V2 API (Priority: P3)

**Goal**: MCP server tools return correct data from V2 API

**Independent Test**: Start MCP server via `pboard mcp serve`, invoke tools via Claude Desktop, verify correct data returned.

### Implementation for User Story 5

- [x] T040 [US5] Pass `apiVersion` from config through to MCP server client creation in `internal/mcp/server.go`
- [x] T039 [US5] Add `members` MCP tools (list_members, get_member) to tool definitions in `internal/mcp/tools.go`
- [x] T040 [US5] Add `teams` MCP tools (list_teams, get_team) to tool definitions in `internal/mcp/tools.go`
- [x] T039 [US5] Add MCP handler functions for members and teams tools in `internal/mcp/handlers.go`
- [x] T040 [US5] Verify all existing MCP tool handlers work with V2 flattened response data — spot-check list_features, list_notes, get_feature handlers in `internal/mcp/handlers.go`

**Checkpoint**: MCP server fully functional with V2 API

---

## Phase 8: New V2-Only Commands (Priority: P2)

**Goal**: Add `members` and `teams` CLI commands available only with V2

**Independent Test**: Run `pboard members list`, `pboard teams list` with V2 (should work) and V1 (should error).

### Implementation

- [x] T041 [P] Create `pboard members` command with `list` (columns: ID, Name, Email, Role) and `get <id>` subcommands, with V2-only guard, in new file `internal/cli/members.go`
- [x] T042 [P] Create `pboard teams` command with `list` (columns: ID, Name, Handle, Description) and `get <id>` subcommands, with V2-only guard, in new file `internal/cli/teams.go`
- [x] T043 Register `members` and `teams` commands in root command in `internal/cli/root.go`
- [x] T044 Add V2 path entries for `/members` → `/v2/members` and `/teams` → `/v2/teams` in `internal/client/v2.go`

**Checkpoint**: New V2-only commands work; error gracefully with V1

---

## Phase 9: Polish & Cross-Cutting Concerns

**Purpose**: Constitution amendment, documentation, final validation

- [x] T045 Amend constitution API versioning constraint to support configurable V1/V2, bump version 1.0.0 → 1.1.0 in `.specify/memory/constitution.md`
- [x] T046 Run full quickstart.md validation — execute all verification commands from `specs/005-v2-api-migration/quickstart.md`
- [x] T047 Update CLI help text and root command description to mention V2 API support in `internal/cli/root.go`

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — can start immediately
- **Foundational (Phase 2)**: Depends on Setup (T001-T004) — BLOCKS all user stories
- **US1 (Phase 3)**: Depends on Foundational (T005-T010)
- **US2 (Phase 4)**: Depends on Foundational (T005-T010), independent of US1
- **US3 (Phase 5)**: Depends on Foundational (T005-T010), independent of US1/US2
- **US4 (Phase 6)**: Depends on Foundational (T005-T010), independent of others
- **US5 (Phase 7)**: Depends on Foundational (T005-T010) and ideally US1 completion (to verify flattened responses work)
- **New Commands (Phase 8)**: Depends on Foundational (T005-T010), independent of other stories
- **Polish (Phase 9)**: Depends on all phases complete

### User Story Dependencies

- **US1 (P1)**: Can start after Phase 2 — no dependencies on other stories — **MVP**
- **US2 (P2)**: Can start after Phase 2 — independent of US1
- **US3 (P2)**: Can start after Phase 2 — independent of US1/US2
- **US4 (P3)**: Can start after Phase 2 — independent of others
- **US5 (P3)**: Benefits from US1 completion but can start after Phase 2

### Within Each User Story

- Path mapping before command verification
- Core commands before edge case commands
- Story complete before moving to next priority

### Parallel Opportunities

- T001-T004: Sequential (each depends on previous)
- T005, T006: Parallel [P] (both write to v2.go but different functions)
- T007, T008, T009, T010: Sequential (T007 must come first, T008/T009 depend on T007)
- T011-T016: All parallel (different CLI command files)
- T018-T019: Parallel (notes tags + features health, different files)
- T024-T025: Parallel (different command files)
- T032-T034: All parallel (different command files)
- T037-T038: Parallel (same file but independent tool definitions)
- T041-T042: Parallel (different new files)

---

## Parallel Example: User Story 1

```bash
# After foundational phase, launch all entity verifications in parallel:
Task: "T011 [US1] Verify features list/get with V2 in internal/cli/features.go"
Task: "T012 [US1] Verify products list/get with V2 in internal/cli/products.go"
Task: "T013 [US1] Verify components list/get with V2 in internal/cli/components.go"
Task: "T014 [US1] Verify initiatives/objectives/key-results with V2"
Task: "T015 [US1] Verify releases/release-groups with V2"
Task: "T016 [US1] Verify companies/users with V2"
Task: "T018 [US1] Verify notes tags with V2 in internal/cli/notes.go"
Task: "T019 [US1] Verify features health with V2 in internal/cli/features_health.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (T001-T004)
2. Complete Phase 2: Foundational (T005-T010) — CRITICAL, blocks all stories
3. Complete Phase 3: User Story 1 (T011-T021)
4. **STOP and VALIDATE**: Test all list/get commands with V2 and V1 fallback
5. Deploy/demo if ready

### Incremental Delivery

1. Setup + Foundational → V2 client infrastructure ready
2. Add US1 → All entity browsing works → **MVP!**
3. Add US2 → Relationships work
4. Add US3 → Custom fields/statuses work
5. Add US4 → Integrations work
6. Add US5 → MCP server migrated
7. Add New Commands → Members/teams available
8. Polish → Constitution amended, docs updated

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- No test framework exists — verification is manual via quickstart.md commands
- The client-level path translation (Phase 2) is the most complex piece — once that works, most command files need minimal or no changes
- V2 response flattening ensures existing `SafeStr`/`SafeNested` display helpers work unchanged
- Commit after each task or logical group
