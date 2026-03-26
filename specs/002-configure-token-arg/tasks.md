# Tasks: Configure Command Token Argument

**Input**: Design documents from `/specs/002-configure-token-arg/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, contracts/

**Tests**: Not explicitly requested in the feature specification. Tests are not included.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2)
- Include exact file paths in descriptions

## Phase 1: User Story 1 - One-Command Token Setup (Priority: P1)

**Goal**: Users can pass their API token as a positional argument to `pboard configure` for instant setup.

**Independent Test**: Run `pboard configure <token>` and verify the token is saved and subsequent commands authenticate.

### Implementation for User Story 1

- [x] T001 [US1] Modify configure command to accept optional positional argument using `cobra.MaximumNArgs(1)` and save token directly when provided (skip interactive prompt) in internal/cli/configure.go
- [x] T002 [US1] Add whitespace trimming and empty-token validation for the positional argument in internal/cli/configure.go

**Checkpoint**: `pboard configure pb_token_here` saves the token and displays confirmation. No interactive prompt shown.

---

## Phase 2: User Story 2 - Preserve Interactive Fallback (Priority: P2)

**Goal**: Existing interactive prompt behavior is preserved when no argument is provided.

**Independent Test**: Run `pboard configure` without arguments and verify the interactive prompt still works.

### Implementation for User Story 2

- [x] T003 [US2] Verify interactive fallback works when no argument is provided (existing behavior preserved) in internal/cli/configure.go

**Checkpoint**: `pboard configure` (no args) still prompts interactively. `pboard configure <token>` skips the prompt.

---

## Phase 3: Polish & Cross-Cutting Concerns

**Purpose**: Documentation and validation

- [x] T004 Update README.md to document the new `pboard configure <token>` usage
- [x] T005 Run quickstart.md validation — verify `pboard configure <token>` example works

---

## Dependencies & Execution Order

### Phase Dependencies

- **US1 (Phase 1)**: No dependencies — can start immediately
- **US2 (Phase 2)**: Depends on US1 completion (same file modification)
- **Polish (Phase 3)**: Depends on US1 and US2 being complete

### Within Each User Story

- T001 before T002 (same file, T002 refines T001's behavior)
- T003 is a verification task after T001+T002

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete T001 + T002: Token argument support
2. **STOP and VALIDATE**: `pboard configure pb_token` works
3. Complete T003: Verify backward compatibility
4. Complete T004 + T005: Documentation

### Notes

- This is a single-file change (~20 lines modified in internal/cli/configure.go)
- All tasks are sequential (same file)
- Commit after T002 for a complete, working feature
