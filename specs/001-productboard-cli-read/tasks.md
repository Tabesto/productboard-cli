# Tasks: ProductBoard CLI Read-Only Tool (pboard)

**Input**: Design documents from `/specs/001-productboard-cli-read/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: Not explicitly requested in the feature specification. Tests are not included.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [x] T001 Initialize Go module with `go mod init` and create go.mod
- [x] T002 Create project directory structure per plan.md: cmd/pboard/, internal/cli/, internal/client/, internal/config/, internal/output/, internal/models/, testdata/fixtures/
- [x] T003 Install dependencies: cobra, viper, tablewriter via `go get`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**Warning**: No user story work can begin until this phase is complete

- [x] T004 Define all API response structs (Feature, Product, Component, Note, Company, User, Release, ReleaseGroup, Objective, KeyResult, Initiative, CustomField, Webhook, etc.) in internal/models/models.go per data-model.md
- [x] T005 [P] Implement config loading (env var `PRODUCTBOARD_API_TOKEN` with precedence over `~/.config/pboard/config.yaml`, file mode 600 enforcement) in internal/config/config.go
- [x] T006 Implement HTTP client with Bearer token auth, base URL, and request helper in internal/client/client.go
- [x] T007 Implement cursor-based auto-pagination in the HTTP client (follow `pageCursor` until exhausted or `--limit` reached) in internal/client/client.go (depends on T006)
- [x] T008 [P] Implement structured error types and HTTP status code mapping (401→exit 2, 404→exit 3, 429→exit 4, etc.) with user-friendly messages in internal/client/errors.go
- [x] T009 [P] Implement output formatter with table (tablewriter) and JSON (encoding/json) modes, selectable via format string, in internal/output/formatter.go
- [x] T010 Create root command with global flags (`--output`, `--limit`, `--help`, `--version`) in internal/cli/root.go
- [x] T011 Create main.go entry point that executes root command in cmd/pboard/main.go

**Checkpoint**: Foundation ready — `pboard --version` and `pboard --help` work. User story implementation can begin.

---

## Phase 3: User Story 1 - List and Retrieve Features (Priority: P1) MVP

**Goal**: Users can list features, get feature details, and view linked initiatives/objectives.

**Independent Test**: Run `pboard features list` and `pboard features get <id>` to verify feature data is returned.

### Implementation for User Story 1

- [x] T012 [US1] Implement `features list` subcommand with filters (`--status-id`, `--status-name`, `--parent-id`, `--archived`, `--owner-email`, `--note-id`) in internal/cli/features.go
- [x] T013 [US1] Implement `features get <id>` subcommand in internal/cli/features.go
- [x] T014 [US1] Implement `features links initiatives <id>` subcommand in internal/cli/features.go
- [x] T015 [US1] Implement `features links objectives <id>` subcommand in internal/cli/features.go

**Checkpoint**: `pboard features list`, `pboard features get <id>`, and `pboard features links [initiatives|objectives] <id>` work with table and JSON output.

---

## Phase 4: User Story 2 - Browse Product Hierarchy (Priority: P1)

**Goal**: Users can list and retrieve products, components, and feature statuses.

**Independent Test**: Run `pboard products list`, `pboard components list`, `pboard feature-statuses list`.

### Implementation for User Story 2

- [x] T016 [P] [US2] Implement `products list` and `products get <id>` subcommands in internal/cli/products.go
- [x] T017 [P] [US2] Implement `components list` and `components get <id>` subcommands in internal/cli/components.go
- [x] T018 [P] [US2] Implement `feature-statuses list` subcommand in internal/cli/feature_statuses.go

**Checkpoint**: Product hierarchy commands work. Combined with US1, users can browse features with full context.

---

## Phase 5: User Story 9 - Install via Homebrew (Priority: P1)

**Goal**: Users can install pboard via `brew install`.

**Independent Test**: Run `brew install <tap>/<formula>` and verify `pboard --version` works.

### Implementation for User Story 9

- [x] T019 [US9] Create .goreleaser.yaml with cross-compilation config (macOS arm64/amd64, Linux amd64) and Homebrew tap formula generation
- [x] T020 [US9] Create GitHub Actions workflow for release automation (tag push → GoReleaser → GitHub Release → Homebrew tap update) in .github/workflows/release.yml

**Checkpoint**: Tagging a release produces binaries and updates the Homebrew formula.

---

## Phase 6: User Story 10 - Configure Authentication (Priority: P1)

**Goal**: Users can configure their API token interactively.

**Independent Test**: Run `pboard configure`, enter token, then verify `pboard features list` works.

### Implementation for User Story 10

- [x] T021 [US10] Implement `configure` command with interactive token prompt and config file write (~/.config/pboard/config.yaml, mode 600) in internal/cli/configure.go

**Checkpoint**: `pboard configure` writes token to config file. All commands authenticate using the stored token.

---

## Phase 7: User Story 3 - View Notes and Feedback (Priority: P2)

**Goal**: Users can list/retrieve notes with tags and links, and view feedback form configurations.

**Independent Test**: Run `pboard notes list` and `pboard notes get <id>`.

### Implementation for User Story 3

- [x] T022 [US3] Implement `notes list` with filters (`--date-from`, `--date-to`, `--created-from`, `--created-to`, `--updated-from`, `--updated-to`, `--term`, `--feature-id`, `--company-id`, `--owner-email`, `--source`, `--any-tag`, `--all-tags`) in internal/cli/notes.go
- [x] T023 [US3] Implement `notes get <id>`, `notes tags <noteId>`, `notes links <noteId>` subcommands in internal/cli/notes.go
- [x] T024 [P] [US3] Implement `feedback-forms list` and `feedback-forms get <id>` subcommands in internal/cli/feedback_forms.go

**Checkpoint**: Notes and feedback form commands work with all filters.

---

## Phase 8: User Story 4 - View Releases and Assignments (Priority: P2)

**Goal**: Users can view releases, release groups, and feature release assignments.

**Independent Test**: Run `pboard releases list` and `pboard release-groups list`.

### Implementation for User Story 4

- [x] T025 [P] [US4] Implement `releases list` (with `--release-group-id` filter) and `releases get <id>` in internal/cli/releases.go
- [x] T026 [P] [US4] Implement `release-groups list` and `release-groups get <id>` in internal/cli/release_groups.go
- [x] T027 [P] [US4] Implement `feature-release-assignments list` (with `--feature-id`, `--release-id`, `--release-state`, `--end-date-from`, `--end-date-to` filters) and `feature-release-assignments get` (with required `--feature-id`, `--release-id`) in internal/cli/assignments.go

**Checkpoint**: Release-related commands work with filters.

---

## Phase 9: User Story 5 - View Objectives, Key Results, and Initiatives (Priority: P2)

**Goal**: Users can view OKR data and initiative relationships.

**Independent Test**: Run `pboard objectives list`, `pboard key-results list`, `pboard initiatives list`.

### Implementation for User Story 5

- [x] T028 [P] [US5] Implement `objectives list`, `objectives get <id>`, `objectives links features <id>`, `objectives links initiatives <id>` in internal/cli/objectives.go
- [x] T029 [P] [US5] Implement `key-results list` and `key-results get <id>` in internal/cli/key_results.go
- [x] T030 [P] [US5] Implement `initiatives list`, `initiatives get <id>`, `initiatives links objectives <id>`, `initiatives links features <id>` in internal/cli/initiatives.go

**Checkpoint**: OKR and initiative commands work with linked resource navigation.

---

## Phase 10: User Story 11 - Flexible Output Formats (Priority: P2)

**Goal**: Users can switch between table and JSON output.

**Independent Test**: Run `pboard features list -o json` and `pboard features list -o table`.

### Implementation for User Story 11

- [x] T031 [US11] Verify all existing commands pass `--output` flag through to formatter, add table column definitions for each resource type in internal/output/formatter.go

**Checkpoint**: All commands support `-o json` and `-o table` consistently.

---

## Phase 11: User Story 6 - View Companies and Users (Priority: P3)

**Goal**: Users can list/retrieve companies and users with custom fields.

**Independent Test**: Run `pboard companies list` and `pboard users list`.

### Implementation for User Story 6

- [x] T032 [P] [US6] Implement `companies list` (with `--term`, `--has-notes`, `--feature-id` filters), `companies get <id>`, `companies custom-fields list`, `companies custom-fields get <id>`, `companies custom-field-value <companyId> <fieldId>` in internal/cli/companies.go
- [x] T033 [P] [US6] Implement `users list` and `users get <id>` in internal/cli/users.go

**Checkpoint**: Company and user commands work with custom field access.

---

## Phase 12: User Story 7 - View Custom Fields (Priority: P3)

**Goal**: Users can list/retrieve hierarchy entity custom field definitions and values.

**Independent Test**: Run `pboard custom-fields list --type text`.

### Implementation for User Story 7

- [x] T034 [US7] Implement `custom-fields list` (with required `--type` filter), `custom-fields get <id>`, `custom-fields values list` (with `--type`, `--custom-field-id`, `--hierarchy-entity-id` filters), `custom-fields values get` (with required `--custom-field-id`, `--hierarchy-entity-id`) in internal/cli/custom_fields.go

**Checkpoint**: Custom field commands work with all required and optional filters.

---

## Phase 13: User Story 8 - View Integrations and Webhooks (Priority: P3)

**Goal**: Users can view plugin/Jira integrations and webhook subscriptions.

**Independent Test**: Run `pboard plugin-integrations list`, `pboard jira-integrations list`, `pboard webhooks list`.

### Implementation for User Story 8

- [x] T035 [P] [US8] Implement `plugin-integrations list`, `plugin-integrations get <id>`, `plugin-integrations connections list <id>`, `plugin-integrations connections get <id> <featureId>` in internal/cli/plugin_integrations.go
- [x] T036 [P] [US8] Implement `jira-integrations list`, `jira-integrations get <id>`, `jira-integrations connections list <id>` (with `--issue-key`, `--issue-id` filters), `jira-integrations connections get <id> <featureId>` in internal/cli/jira_integrations.go
- [x] T037 [P] [US8] Implement `webhooks list` and `webhooks get <id>` in internal/cli/webhooks.go

**Checkpoint**: All integration and webhook commands work.

---

## Phase 14: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [x] T038 Add shell completion generation (bash, zsh, fish) as a `completion` subcommand in internal/cli/root.go
- [x] T039 Add usage examples to `--help` text for all commands in internal/cli/*.go
- [x] T040 Run quickstart.md validation — verify all example commands from quickstart.md work correctly
- [x] T041 Create README.md with installation, configuration, and usage instructions

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion — BLOCKS all user stories
- **User Stories (Phase 3–13)**: All depend on Foundational phase completion only — can start in any order
  - Recommended priority order: P1 (US1, US2, US9, US10) → P2 (US3, US4, US5, US11) → P3 (US6, US7, US8)
  - All stories can proceed in parallel if capacity allows
- **Polish (Phase 14)**: Depends on all user stories being complete

### User Story Dependencies

- **US1** (Features): Foundational only — no story dependencies
- **US2** (Product Hierarchy): Foundational only — no story dependencies
- **US9** (Homebrew): Foundational only — needs at least one working command for validation
- **US10** (Auth Config): Foundational only — no story dependencies
- **US3** (Notes): Foundational only — no story dependencies
- **US4** (Releases): Foundational only — no story dependencies
- **US5** (OKRs): Foundational only — no story dependencies
- **US11** (Output Formats): Depends on at least one command existing (US1 recommended)
- **US6** (Companies/Users): Foundational only — no story dependencies
- **US7** (Custom Fields): Foundational only — no story dependencies
- **US8** (Integrations): Foundational only — no story dependencies

### Within Each User Story

- Models before services (handled in Foundational)
- CLI commands use shared client and formatter
- Core implementation before integration

### Parallel Opportunities

- T005, T006, T008, T009 can all run in parallel (Foundational phase)
- T016, T017, T018 can all run in parallel (US2)
- T025, T026, T027 can all run in parallel (US4)
- T028, T029, T030 can all run in parallel (US5)
- T032, T033 can run in parallel (US6)
- T035, T036, T037 can all run in parallel (US8)
- All P1 user stories can run in parallel after Foundational
- All P2 user stories can run in parallel after Foundational
- All P3 user stories can run in parallel after Foundational

---

## Parallel Example: User Story 2

```bash
# Launch all US2 commands in parallel (different files):
Task: "T016 Implement products commands in internal/cli/products.go"
Task: "T017 Implement components commands in internal/cli/components.go"
Task: "T018 Implement feature-statuses command in internal/cli/feature_statuses.go"
```

## Parallel Example: User Story 5

```bash
# Launch all US5 commands in parallel (different files):
Task: "T028 Implement objectives commands in internal/cli/objectives.go"
Task: "T029 Implement key-results commands in internal/cli/key_results.go"
Task: "T030 Implement initiatives commands in internal/cli/initiatives.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL — blocks all stories)
3. Complete Phase 3: User Story 1 (Features)
4. **STOP and VALIDATE**: `pboard features list` and `pboard features get <id>` work
5. Deploy/demo if ready

### Incremental Delivery

1. Setup + Foundational → Foundation ready
2. Add US1 (Features) + US2 (Hierarchy) + US10 (Auth) → Core experience (MVP!)
3. Add US3 (Notes) + US4 (Releases) + US5 (OKRs) + US11 (Output) → Full read experience
4. Add US6 (Companies) + US7 (Custom Fields) + US8 (Integrations) → Complete coverage
5. Add US9 (Homebrew) + Polish → Release-ready
6. Each story adds value without breaking previous stories

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together
2. Once Foundational is done:
   - Developer A: US1 (Features) + US3 (Notes)
   - Developer B: US2 (Hierarchy) + US4 (Releases) + US5 (OKRs)
   - Developer C: US10 (Auth) + US6 (Companies) + US7 (Custom Fields) + US8 (Integrations)
3. Stories complete and integrate independently
4. US9 (Homebrew) + US11 (Output) + Polish as final pass

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story is independently completable and testable
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- All commands follow the same pattern: register Cobra command, call client, format output
