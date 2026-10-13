---
description: "Task list for gated write support"
---

# Tasks: Gated Write Support

**Input**: Design documents from `/specs/006-gated-write-support/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/interfaces.md, quickstart.md

**Tests**: Included (plan.md lists test files; repo has no tests today). Keep them small: `httptest`, `t.TempDir()` + `HOME` override. No frameworks.

## Format: `[ID] [P?] [Story] Description`

Module path: `github.com/tabesto/productboard-cli`. All paths relative to repo root.

---

## Phase 1: Setup

- [x] T001 Confirm baseline: run `go build ./... && go vet ./...` on branch `006-gated-write-support` and record that it passes before changes

---

## Phase 2: Foundational (blocks all stories)

**Purpose**: gate + non-GET plumbing shared by CLI and MCP.

- [x] T002 Add `ExitWriteDisabled = 6` and `NewWriteDisabledError()` in `internal/client/errors.go`. Message (exact): `write operations are disabled. Enable with --allow-write, \`pboard configure --write-enabled\`, or PRODUCTBOARD_WRITE_ENABLED=1.`
- [x] T003 Add `WriteEnabled bool` to `Config` in `internal/config/config.go`; in `Load()` read viper key `write_enabled` then env `PRODUCTBOARD_WRITE_ENABLED` (const `EnvWriteEnabledKey`): value `1`/`true` (case-insensitive) → true, any other non-empty value → false, empty → keep config value
- [x] T004 Refactor `internal/client/client.go`: add `allowWrite bool` (from `cfg.WriteEnabled`) and `http.Client{Timeout: 30 * time.Second}` in `New`; extract the request-sending part of `Get` into `do(method, path string, params map[string]string, body []byte) ([]byte, error)` (keep V2 translation in `Get`, keep header logic, add `Content-Type: application/json` when body != nil); `Get` calls `do("GET", …)`. In `do`, for any method other than GET call `RequireWrite()` first
- [x] T005 Add to `internal/client/client.go`: exported `RequireWrite() error` (returns `NewWriteDisabledError()` if `!allowWrite`, else `ExitInvalidInput` APIError "write operations require API version 2" if `!IsV2()`), plus `Post(path string, body any) ([]byte, error)` and `Patch(path string, body any) ([]byte, error)` that JSON-marshal body and call `do` with path prefixed `/v2` (no `translateV2Path`)
- [x] T006 [P] Create `internal/client/write.go` with `type WriteRef struct{ID, Type, Link string}` and `ParseWriteRef(body []byte) (*WriteRef, error)` reading `data.id`, `data.type`, `data.links.self`
- [x] T007 Add persistent flag `--allow-write` in `internal/cli/root.go` (package var `allowWriteFlag`); in `getClient()` set `cfg.WriteEnabled = cfg.WriteEnabled || allowWriteFlag`; update root `Short`/`Long` text from "read-only" to "read-only by default; writes need --allow-write"

**Checkpoint**: build passes; no behavior change for existing read commands.

---

## Phase 3: User Story 1 - Write stays off unless enabled (P1)

**Goal**: any write refused with no network I/O unless one of flag/config/env is on.
**Independent Test**: `go test ./internal/client ./internal/config`; manual: quickstart §1.

- [x] T008 [P] [US1] Create `internal/client/client_test.go`: `httptest.Server` counting requests; with `allowWrite=false` `Post`/`Patch` return `*APIError` with `ExitCode==6` and server saw 0 requests; with `allowWrite=true` server sees method, path `/v2/entities`, `Authorization: Bearer …`, `Content-Type: application/json`, no `X-Version`, body `{"data":…}`; `apiVersion="1"` + allowWrite → `ExitInvalidInput`; GET still works without gate
- [x] T009 [P] [US1] Create `internal/config/config_test.go`: set `HOME` to `t.TempDir()`; table test of `Load()` for env values (`1`, `true`, `TRUE`, `0`, `false`, `yes`, empty) and config-file `write_enabled: true` with env `0` overriding to false. NOTE: viper is a global — call `viper.Reset()` in each case
- [x] T010 [US1] Verify existing read commands unchanged: `go build -o pboard ./cmd/pboard && ./pboard features list --help` still works; document result in PR notes (no code)

---

## Phase 4: User Story 2 - Create and update entities (P1)

**Goal**: `pboard entities create|update` with `--field`, `--patch`, `--json`.
**Independent Test**: `go test ./internal/cli`; manual: quickstart §2 on a test workspace.

- [x] T011 [P] [US2] Create `internal/cli/write.go` with helpers: `parseFieldPairs([]string) (map[string]any, error)` (`key=value`; split on first `=`; value JSON-decoded if valid else string; error → `ExitInvalidInput` APIError), `buildUpdateBody(fields []string, patch, rawJSON string) (map[string]any, error)` (exactly one of the three; patch must be a JSON array; json must be a JSON object; returns `{"data":…}`), and `printWriteResult(body []byte) error` using `client.ParseWriteRef` + `output.PrintSingleTable`/`output.PrintJSON` per `outputFormat`
- [x] T012 [P] [US2] Create `internal/cli/write_test.go`: table tests for `parseFieldPairs` (string, number, bool, object, `a=b=c`, missing `=`) and `buildUpdateBody` (fields only, patch only, json only, fields+patch rejected, none rejected, invalid patch JSON)
- [x] T013 [US2] Create `internal/cli/entities.go`: `newEntitiesCmd()` group (`entities`) with `create` (`--type`, `--field` StringArray, `--json`) and `update <id>` (`--field`, `--patch`, `--json`). Each `RunE`: `c := getClient()` → `c.RequireWrite()` first → build body → `c.Post("/entities", body)` / `c.Patch("/entities/"+id, body)` → `printWriteResult`. Errors via `handleError`
- [x] T014 [US2] Register `newEntitiesCmd()` in `internal/cli/root.go`

---

## Phase 5: User Story 4 - MCP write tools only when enabled (P1)

**Goal**: 4 tools exist only if gate on at server start.
**Independent Test**: `go test ./internal/mcp`; manual: quickstart §6.

- [x] T015 [US4] `internal/mcp/server.go`: `Serve(allowWrite bool)`; load config, `effective := cfg.WriteEnabled || allowWrite`; store in package var `writeEnabled` used by `getClient()` in `internal/mcp/handlers.go` (set `cfg.WriteEnabled = cfg.WriteEnabled || writeEnabled`); call `registerTools(s, effective)`
- [x] T016 [US4] `internal/mcp/tools.go`: change to `registerTools(s *server.MCPServer, allowWrite bool)`; when true call new `registerWriteTools(s)` registering `create_entity` (`type` req, `fields` object), `update_entity` (`id` req, `fields` object, `patch` array), `create_note` (`name` req, `type` default textNote, `content`, `fields` object), `update_note` (`id` req, `fields`, `patch`). Descriptions must tell agent to inspect entity/note configuration and existing records before writing
- [x] T017 [US4] Create `internal/mcp/handlers_write.go` with the four handlers: `getClient()` → `c.RequireWrite()` → build body (reject fields+patch together) → `c.Post/Patch` → return JSON of `ParseWriteRef` via `mcp.NewToolResultText`; errors via `errorResult`
- [x] T018 [P] [US4] Create `internal/mcp/tools_test.go`: build `server.NewMCPServer`, call `registerTools(s, false)` → no tool named `create_*`/`update_*`; `registerTools(s, true)` → exactly those 4 extra tools (use `s.ListTools()`; if unavailable in mcp-go v0.46.0 use an in-process `HandleMessage` `tools/list`)
- [x] T019 [US4] `internal/cli/mcp.go`: add `--allow-write` bool to `serve` (passes to `pbmcp.Serve(allowWrite)`) and to `install` (args become `["mcp","serve","--allow-write"]` when set, also in `--dry-run` output)

---

## Phase 6: User Story 3 - Notes create/update (P2)

**Goal**: `pboard notes create|update`.
**Independent Test**: manual quickstart §4; unit coverage via shared helpers.

- [x] T020 [US3] In `internal/cli/notes.go` add `newNotesCreateCmd()` (`--type` default `textNote`, `--name`, `--content`, `--field`, `--json`; `--name`/`--content` merge into `fields`; error exit 5 if no `name` after merge unless `--json` supplied) and `newNotesUpdateCmd()` (reuse `buildUpdateBody`); `RequireWrite()` first; `c.Post("/notes", …)` / `c.Patch("/notes/"+id, …)`; `printWriteResult`. Register both in `newNotesCmd()`
- [x] T021 [P] [US3] Extend `internal/cli/write_test.go` with a test for the note-create body builder (name/content shortcuts merged with `--field`, missing name rejected). Put the builder in `internal/cli/write.go` as `buildNoteCreateBody(typ, name, content string, fields []string, rawJSON string)`

---

## Phase 7: User Story 5 - Persistent enable/disable (P2)

**Goal**: `configure --write-enabled`, config never loses keys.
**Independent Test**: `go test ./internal/config`; manual quickstart §3 and §7.

- [x] T022 [US5] In `internal/config/config.go` replace `WriteConfig` with `Update(values map[string]any) error`: new `viper.New()` instance, `SetConfigFile(path)`, `SetConfigPermissions(0600)`, read existing file if present, set given keys, default `api_url`/`api_version` only when absent, `WriteConfigAs(path)`, then `os.Chmod(path, 0600)`; dir mode 0700
- [x] T023 [US5] `internal/cli/configure.go`: add `--write-enabled` bool flag (use `cmd.Flags().Changed`); if flag changed and no token arg → update only `write_enabled` (no prompt); otherwise get token as today and `Update({"api_token": token})` plus `write_enabled` if flag changed; stop resetting `api_version`; print which keys were saved (never the token). Update `Long` examples
- [x] T024 [P] [US5] Add to `internal/config/config_test.go`: `Update` preserves existing keys (write file with api_version "2" + write_enabled true, update token only → all keys present, mode 0600), creates file with defaults when absent

---

## Phase 8: Polish & cross-cutting

- [x] T025 [P] Update `.specify/memory/constitution.md`: Principle I → "Read-Only by Default, Gated Writes" (writes allowed only via the write gate, V2, enforced in `internal/client`, MCP write tools hidden when gate off); add exit code `6=write-disabled` to Development Workflow; bump to 2.0.0, Last Amended 2026-10-13; refresh Sync Impact Report
- [x] T026 [P] Update `README.md`: intro no longer "read-only"; document API V2 default + `api_version`/`--api-version`, write gate (flag/config/env), new commands, `mcp serve|install --allow-write`, exit code 6
- [x] T027 [P] Update `internal/cli/skill_content.go` (embedded pboard skill): section on write commands and the gate; tell agents never to enable writes on their own
- [x] T028 Run `gofmt -l .`, `go vet ./...`, `go test ./...`, `go build -o pboard ./cmd/pboard`
- [x] T029 (not run: live-API steps §2, §4, §5 — no test token) Run quickstart.md §1, §3 (config preservation, disable/enable) and §6 (MCP tools/list with/without gate) locally against the built binary; skip live-API steps (§2, §4, §5) unless the user supplies a test token, and report them as not run

---

## Dependencies & Execution Order

- Phase 1 → Phase 2 (T002→T003 independent; T004→T005→T007; T006 parallel)
- After Phase 2: US1 (T008–T010) and US2 (T011–T014) and US4 (T015–T019) can proceed; US3 needs T011 (helpers); US5 independent of the others after T003.
- Polish after all stories.
- Within a story: tests [P] can be written alongside, but must pass before the story is called done.

## Parallel examples

- After T005: T006, T008, T009, T011, T012, T022 touch different files → parallel.
- Docs T025–T027 in parallel at the end.

## Implementation strategy

1. MVP = Phase 2 + US1 + US2 (gate + entities). Stop and validate.
2. Add US4 (MCP) → US5 (configure) → US3 (notes) → Polish.
3. No live Productboard writes without a user-provided test workspace/token.
