# Tasks: MCP Server for pboard CLI

**Input**: Design documents from `/specs/003-pboard-mcp-server/`
**Prerequisites**: plan.md (required), spec.md (required), research.md, data-model.md, contracts/

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Phase 1: Setup

**Purpose**: Add mcp-go dependency and create package structure

- [x] T001 Add `github.com/mark3labs/mcp-go` dependency in go.mod and run `go mod tidy`
- [x] T002 Create `internal/mcp/` package directory structure

---

## Phase 2: Foundational (MCP Server Core)

**Purpose**: MCP server skeleton and CLI command group that ALL user stories depend on

**CRITICAL**: No user story work can begin until this phase is complete

- [x] T003 Create MCP server initialization with stdio transport in internal/mcp/server.go -- define `NewServer()` that creates an mcp-go server with name "pboard", registers tools, and `Serve()` that calls `server.ServeStdio()`. Include helper for loading config and creating API client. Include error-to-MCP-result mapping per contracts/mcp-tools.md error contract.
- [x] T004 Create `pboard mcp` command group with `serve` subcommand in internal/cli/mcp.go -- `serve` calls `mcp.NewServer()` then `mcp.Serve()`. Follow the pattern in internal/cli/skill.go for command structure.
- [x] T005 Create `pboard mcp install` subcommand in internal/cli/mcp.go -- reads/creates `~/Library/Application Support/Claude/claude_desktop_config.json`, merges `pboard` entry into `mcpServers` with command path resolved via `os.Executable()`, supports `--force` and `--dry-run` flags. Preserve existing MCP server entries. Follow skill.go install pattern.
- [x] T006 Create `pboard mcp uninstall` subcommand in internal/cli/mcp.go -- reads Claude Desktop config, removes only the `pboard` key from `mcpServers`, writes back. Follow skill.go uninstall pattern.
- [x] T007 Register `mcp` command group in internal/cli/root.go by adding `rootCmd.AddCommand(newMcpCmd())` alongside existing commands

**Checkpoint**: `pboard mcp serve` starts and responds to MCP initialize request. `pboard mcp install` writes correct config. Foundation ready.

---

## Phase 3: User Story 1 - Browse ProductBoard Features (Priority: P1) MVP

**Goal**: Users can list, get, and filter features from Claude Desktop

**Independent Test**: Run `pboard mcp install`, restart Claude Desktop, ask Claude to list features and get a feature by ID

### Implementation for User Story 1

- [x] T008 [P] [US1] Define `list_features` tool with parameters (status-name, status-id, parent-id, archived, owner-email, note-id, limit) and handler that calls `client.GetList("/features", ...)` with default limit 25 in internal/mcp/tools.go and internal/mcp/handlers.go
- [x] T009 [P] [US1] Define `get_feature` tool with parameter (id, required) and handler that calls `client.GetSingle("/features/{id}")` in internal/mcp/tools.go and internal/mcp/handlers.go
- [x] T010 [P] [US1] Define `list_feature_initiatives` tool with parameters (id required, limit) and handler that calls `client.GetLinkedResources("/features/{id}/links/initiatives")` in internal/mcp/tools.go and internal/mcp/handlers.go
- [x] T011 [P] [US1] Define `list_feature_objectives` tool with parameters (id required, limit) and handler that calls `client.GetLinkedResources("/features/{id}/links/objectives")` in internal/mcp/tools.go and internal/mcp/handlers.go
- [x] T012 [US1] Register all feature tools in server setup in internal/mcp/server.go

**Checkpoint**: Features fully queryable from Claude Desktop. MVP complete.

---

## Phase 4: User Story 2 - Query Notes and Feedback (Priority: P1)

**Goal**: Users can list, search, and get notes with all filters from Claude Desktop

**Independent Test**: Ask Claude to search notes by keyword, list notes for a feature, get a specific note

### Implementation for User Story 2

- [x] T013 [P] [US2] Define `list_notes` tool with parameters (term, created-from, created-to, updated-from, updated-to, feature-id, company-id, owner-email, source, tags, limit) and handler that calls `client.GetList("/notes", ...)` with default limit 25 in internal/mcp/tools.go and internal/mcp/handlers.go
- [x] T014 [P] [US2] Define `get_note` tool with parameter (id, required) and handler that calls `client.GetSingle("/notes/{id}")` in internal/mcp/tools.go and internal/mcp/handlers.go
- [x] T015 [P] [US2] Define `list_note_tags` tool with parameters (id required, limit) and handler in internal/mcp/tools.go and internal/mcp/handlers.go
- [x] T016 [P] [US2] Define `list_note_links` tool with parameters (id required, limit) and handler in internal/mcp/tools.go and internal/mcp/handlers.go
- [x] T017 [P] [US2] Define `list_feedback_forms` tool with parameter (limit) and handler in internal/mcp/tools.go and internal/mcp/handlers.go
- [x] T018 [US2] Register all notes and feedback tools in server setup in internal/mcp/server.go

**Checkpoint**: Notes and feedback fully queryable. US1 + US2 complete.

---

## Phase 5: User Story 3 - Explore Product Hierarchy and Releases (Priority: P2)

**Goal**: Users can browse products, components, releases, objectives, initiatives, and OKR linkages

**Independent Test**: Ask Claude to list releases, show initiatives linked to a feature, list objectives

### Implementation for User Story 3

- [x] T019 [P] [US3] Define `list_products` and `get_product` tools with handlers in internal/mcp/tools.go and internal/mcp/handlers.go
- [x] T020 [P] [US3] Define `list_components` tool with handler in internal/mcp/tools.go and internal/mcp/handlers.go
- [x] T021 [P] [US3] Define `list_releases`, `get_release`, `list_release_groups`, and `list_feature_release_assignments` tools (with release-group-id and release-id filters respectively) and handlers in internal/mcp/tools.go and internal/mcp/handlers.go
- [x] T022 [P] [US3] Define `list_objectives`, `get_objective`, and `list_objective_features` tools with handlers in internal/mcp/tools.go and internal/mcp/handlers.go
- [x] T023 [P] [US3] Define `list_key_results` tool with handler in internal/mcp/tools.go and internal/mcp/handlers.go
- [x] T024 [P] [US3] Define `list_initiatives`, `get_initiative`, `list_initiative_objectives`, and `list_initiative_features` tools with handlers in internal/mcp/tools.go and internal/mcp/handlers.go
- [x] T025 [P] [US3] Define `list_feature_statuses` tool with handler in internal/mcp/tools.go and internal/mcp/handlers.go
- [x] T026 [US3] Register all hierarchy and release tools in server setup in internal/mcp/server.go

**Checkpoint**: Full product hierarchy and releases queryable. US1 + US2 + US3 complete.

---

## Phase 6: User Story 4 - Look Up Companies and Users (Priority: P2)

**Goal**: Users can search companies, view custom fields, and list users

**Independent Test**: Ask Claude to search companies by name, get company custom field values, list users

### Implementation for User Story 4

- [x] T027 [P] [US4] Define `list_companies` (with term, has-notes, feature-id filters) and `get_company` tools with handlers in internal/mcp/tools.go and internal/mcp/handlers.go
- [x] T028 [P] [US4] Define `list_company_custom_fields`, `get_company_custom_field`, and `get_company_custom_field_value` tools with handlers in internal/mcp/tools.go and internal/mcp/handlers.go
- [x] T029 [P] [US4] Define `list_users` tool with handler in internal/mcp/tools.go and internal/mcp/handlers.go
- [x] T030 [P] [US4] Define `list_custom_fields` (with type required), `get_custom_field`, and `list_custom_field_values` tools with handlers in internal/mcp/tools.go and internal/mcp/handlers.go
- [x] T031 [P] [US4] Define `list_plugin_integrations`, `list_jira_integrations`, `list_webhooks`, and `get_webhook` tools with handlers in internal/mcp/tools.go and internal/mcp/handlers.go
- [x] T032 [US4] Register all company, user, custom field, and integration tools in server setup in internal/mcp/server.go

**Checkpoint**: All ~37 tools registered. Full pboard CLI coverage achieved (SC-004).

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Documentation and final validation

- [x] T033 Update README.md with MCP server section: overview, install instructions (`pboard mcp install`), available tools summary
- [x] T034 Build binary and run end-to-end validation per quickstart.md -- verify `pboard mcp serve` responds to initialize, `pboard mcp install` writes correct config, tools return data from ProductBoard

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Phase 1 - BLOCKS all user stories
- **User Stories (Phase 3-6)**: All depend on Phase 2 completion
  - US1 and US2 are both P1 and can proceed in parallel
  - US3 and US4 are both P2 and can proceed in parallel (after or alongside US1/US2)
- **Polish (Phase 7)**: Depends on all user stories being complete

### User Story Dependencies

- **US1 - Features (P1)**: Can start after Phase 2 - No dependencies on other stories
- **US2 - Notes (P1)**: Can start after Phase 2 - No dependencies on other stories
- **US3 - Hierarchy/Releases (P2)**: Can start after Phase 2 - No dependencies on other stories
- **US4 - Companies/Users (P2)**: Can start after Phase 2 - No dependencies on other stories
- **US5 - Install Command**: Covered in Phase 2 (foundational)

### Within Each User Story

- Define tools and handlers (parallelizable across different resources)
- Register tools in server (after all tools for the story are defined)

### Parallel Opportunities

- T008-T011 (US1 tools) can all run in parallel
- T013-T017 (US2 tools) can all run in parallel
- T019-T025 (US3 tools) can all run in parallel
- T027-T031 (US4 tools) can all run in parallel
- US1 and US2 can run in parallel (both P1)
- US3 and US4 can run in parallel (both P2)

---

## Parallel Example: User Story 1

```bash
# Launch all feature tools in parallel (different tool definitions, no dependencies):
Task: "T008 Define list_features tool in internal/mcp/tools.go and internal/mcp/handlers.go"
Task: "T009 Define get_feature tool in internal/mcp/tools.go and internal/mcp/handlers.go"
Task: "T010 Define list_feature_initiatives tool in internal/mcp/tools.go and internal/mcp/handlers.go"
Task: "T011 Define list_feature_objectives tool in internal/mcp/tools.go and internal/mcp/handlers.go"

# Then register all (depends on above):
Task: "T012 Register all feature tools in internal/mcp/server.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (add mcp-go dependency)
2. Complete Phase 2: Foundational (server core + install command)
3. Complete Phase 3: User Story 1 (feature tools)
4. **STOP and VALIDATE**: Install in Claude Desktop, query features
5. Deploy/demo if ready

### Incremental Delivery

1. Setup + Foundational -> Server starts, install works
2. Add US1 (features) -> Test independently -> MVP!
3. Add US2 (notes) -> Test independently -> Core product data complete
4. Add US3 (hierarchy/releases) -> Full product context
5. Add US4 (companies/users/integrations) -> 100% CLI coverage
6. Polish -> README, end-to-end validation

---

## Notes

- [P] tasks = different files or different tool definitions, no dependencies
- [Story] label maps task to specific user story for traceability
- All tool handlers follow same pattern: load config -> create client -> call API method -> return JSON as MCP text result
- Error handling follows contracts/mcp-tools.md error contract for all tools
- Default limit of 25 applied in all list tool handlers when no limit parameter provided
