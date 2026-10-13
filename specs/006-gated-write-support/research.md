# Research: Gated Write Support

## R1. API version and currency
- **Decision**: Writes target API V2 only.
- **Rationale**: V2 GA 2026-04-09; V1 sunset 2026-07-08 (past). No v3 exists. Changelog through 2026-10-08 has no breaking change to paths used here. `X-Version` header not needed in V2.
- **Alternatives**: Support V1 writes — rejected, V1 is retired.

## R2. V2 write endpoints (base `https://api.productboard.com/v2`)
| Op | Request | Success |
|----|---------|---------|
| Create entity | `POST /entities` `{"data":{"type":"<type>","fields":{…}}}` | 201 `{"data":{"id","type","links":{"self"}}}` |
| Update entity | `PATCH /entities/{id}` `{"data":{"fields":{…}}}` **or** `{"data":{"patch":[…]}}` | 200 minimal ref |
| Create note | `POST /notes` `{"data":{"type":"textNote\|conversationNote","fields":{"name":…,"content":…},"relationships":[…]?}}` | 201 minimal ref |
| Update note | `PATCH /notes/{id}` same fields/patch rule | 200 |

- `fields` and `patch` are mutually exclusive; patch ops: `set`, `addItems`, `removeItems`, `clear`.
- `fields.name` required for notes; `opportunityNote` not creatable.
- Field names/types are workspace-specific; discovered via `GET /entities/configurations` and `/notes/configurations`.
- Since 2026-09-22 an unavailable entity `type` returns `400 request.invalid` with the available types in `detail`.
- Exact `patch` element schema is passed through verbatim from the user (`--patch '<json array>'`) rather than modelled, to avoid drifting from the spec.
- Rate limit 50 req/s per token; 429 with `Retry-After`.

## R3. Gate semantics and precedence
- **Decision**: gate = `flag || env || config` ("any one enables"), default off. Environment variable, when set, overrides the config value (consistent with token/API version): `1`/`true` (case-insensitive) → on; any other non-empty value → off. Flag only ever turns writes on for that run.
- **Rationale**: matches user's choice "flag OR config"; env-over-config lets an operator force off with `PRODUCTBOARD_WRITE_ENABLED=0` for a session.
- **Alternatives**: flag AND config (rejected as heavy); config only (rejected, no per-run opt-in).

## R4. Where to enforce
- **Decision**: in `Client.do()` for any method other than GET, plus an exported `RequireWrite()` that commands call first so a refusal precedes local validation noise and any network I/O.
- **Rationale**: single choke point for CLI and MCP; testable with httptest by asserting zero requests.
- **Alternatives**: per-command checks (rejected: easy to forget on a new command).

## R5. Exit code
- **Decision**: new `ExitWriteDisabled = 6`. Writes under API V1 use `ExitInvalidInput` (5).

## R6. MCP exposure
- **Decision**: `registerTools(s, allowWrite)`; the four write tools are registered only if the effective gate is true at start. `mcp serve --allow-write` sets it; otherwise config/env decide. `mcp install --allow-write` appends the flag to the registered args. Handlers still go through the client gate (defense in depth).
- **Rationale**: an agent cannot call a tool it cannot see; protects against prompt-injected writes.

## R7. Config persistence
- **Problem**: `WriteConfig` rewrites the whole file via template, resetting `api_url`/`api_version`; it would erase `write_enabled`.
- **Decision**: replace with `config.Update(map[string]any)` using a fresh Viper instance: read existing file (if any), set given keys, ensure defaults for absent `api_url`/`api_version`, write with mode 0600 (`SetConfigPermissions` + chmod). No new dependency (Viper already serializes YAML).
- `configure`: token optional when `--write-enabled` is explicitly passed (no prompt); otherwise unchanged behavior. Token-only run no longer resets API version.

## R8. Input format for fields
- **Decision**: `--field name=value` repeatable; value JSON-decoded when valid (numbers, booleans, objects, arrays), else string. `--json` takes a complete `data` object for advanced cases. `--patch` takes the patch array. `--field` and `--patch` mutually exclusive; checked locally.
- **Alternatives**: one flag per field (rejected: fields are workspace-specific).

## R9. Response handling
- **Decision**: print the minimal ref (`id`, `type`, `links.self`) as table or JSON; no automatic follow-up GET.
- **Rationale**: simplest; migration guide recommends GET on `links.self` only when needed; users have `features get` etc.

## R10. Client hardening
- Add 30 s timeout to the shared `http.Client` (currently none). No retry/backoff (out of scope).

## R11. Known adjacent gaps (out of scope, to track)
- `feature-release-assignments` paths not translated to V2.
- Hardcoded `type[]` values may 400 on workspaces lacking those types (2026-09-22 change).
- Release workflow pins Go 1.22 while go.mod says 1.26.1.
