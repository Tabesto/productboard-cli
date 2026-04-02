# Research: Migrate CLI to ProductBoard V2 API

**Branch**: `005-v2-api-migration` | **Date**: 2026-04-01

## R1: V2 API Structure vs V1

**Decision**: V2 uses URL path prefix `/v2` instead of `X-Version: 1` header. Entity-type endpoints are unified under `/v2/entities` with type filters.

**Rationale**: The V2 API fundamentally restructures the endpoint hierarchy. Instead of per-type endpoints (`/features`, `/products`, `/components`), V2 uses a single `/v2/entities` endpoint with `type[]` query parameter. Notes remain separate at `/v2/notes`. Integrations and webhooks keep their own paths under `/v2/`.

**Alternatives considered**:
- Wrap V1 paths to V2 at the HTTP transport level — rejected because response structures differ significantly
- Use V2 OpenAPI codegen — rejected to avoid adding dependencies and maintain the existing untyped `map[string]interface{}` approach

## R2: V2 Response Structure Mapping

**Decision**: V2 entities wrap fields in a `data.fields` sub-object. The client must extract fields and flatten them for display compatibility.

**Rationale**: V1 returns `{"data": {"id": "...", "name": "...", ...}}`. V2 returns `{"data": {"id": "...", "type": "feature", "fields": {"name": "...", "status": {...}}, "links": {...}}}`. The CLI's `SafeStr` and `SafeNested` helpers read from the top level, so V2 responses must be flattened: `id` and `type` stay at top level, everything under `fields` gets merged up.

**Alternatives considered**:
- Change all CLI command display code to read from `fields` sub-object — rejected because it would touch every command file and break V1 compatibility
- Create separate V2 model structs — rejected because the project uses `map[string]interface{}` throughout, not typed structs for API responses

## R3: V2 Pagination

**Decision**: V2 uses cursor-based pagination with `pageCursor` query param and `links.next` in response, identical mechanism to V1.

**Rationale**: V1 uses `pageCursor` field at response root. V2 uses `pageCursor` query param with `links.next` containing the full URL for next page. The client can continue using the `pageCursor` approach but must extract it from the `links.next` URL or parse the response differently.

**Alternatives considered**: None — pagination mechanism is compatible, only the response field location differs.

## R4: V2 Error Response Format

**Decision**: V2 errors return `{"id": "...", "errors": [{"code": "...", "title": "...", "detail": "..."}]}`. The error handler must parse this new structure.

**Rationale**: V1 error bodies are plain text or simple JSON. V2 follows JSON:API-style error arrays. The `NewAPIError` function must attempt to parse V2 error format and extract the first error's `detail` for display, falling back to the raw body if parsing fails.

**Alternatives considered**: Ignore V2 error details and keep generic messages — rejected because V2 provides more actionable error info (e.g., which scope is missing).

## R5: Dual V1/V2 Client Architecture

**Decision**: Add an `APIVersion` field to `Config` (default `"2"`). The `Client` adapts request building based on version: V1 uses `X-Version: 1` header with direct paths, V2 uses `/v2` prefix with entity-type mapping.

**Rationale**: The user confirmed dual support with V2 default. The simplest approach is to make the client version-aware at the request level:
- V1: `GET /features` with `X-Version: 1` header
- V2: `GET /v2/entities?type[]=feature` with no version header

The `--api-version` persistent flag and `api_version` config key control this. Each CLI command continues to call `c.GetList("/features", params, limit)` — the client translates the path internally.

**Alternatives considered**:
- Two separate client implementations with an interface — rejected as over-engineered for a path/header change
- Proxy pattern wrapping V1 client — rejected because V2 response parsing also differs

## R6: V1-to-V2 Path Mapping

**Decision**: Build a path mapping table in the client that converts V1-style paths to V2 equivalents.

**Rationale**: CLI commands currently call `c.GetList("/features", ...)`, `c.GetSingle("/features/ID")`, etc. For V2, these must become:

| V1 Path | V2 Path | Notes |
|---------|---------|-------|
| `/features` | `/v2/entities?type[]=feature` | Type filter |
| `/features/{id}` | `/v2/entities/{id}` | Direct |
| `/products` | `/v2/entities?type[]=product` | Type filter |
| `/components` | `/v2/entities?type[]=component` | Type filter |
| `/initiatives` | `/v2/entities?type[]=initiative` | Type filter |
| `/objectives` | `/v2/entities?type[]=objective` | Type filter |
| `/key-results` | `/v2/entities?type[]=keyResult` | Type filter |
| `/releases` | `/v2/entities?type[]=release` | Type filter |
| `/release-groups` | `/v2/entities?type[]=releaseGroup` | Type filter |
| `/companies` | `/v2/entities?type[]=company` | Type filter |
| `/users` | `/v2/entities?type[]=user` | Type filter |
| `/notes` | `/v2/notes` | Separate endpoint |
| `/notes/{id}` | `/v2/notes/{id}` | Direct |
| `/feature-statuses` | `/v2/entities/configurations/feature` | Config endpoint |
| `/features/{id}/links/initiatives` | `/v2/entities/{id}/relationships?target[type]=initiative` | Relationships |
| `/features/{id}/links/objectives` | `/v2/entities/{id}/relationships?target[type]=objective` | Relationships |
| `/objectives/{id}/links/features` | `/v2/entities/{id}/relationships?target[type]=feature` | Relationships |
| `/objectives/{id}/links/initiatives` | `/v2/entities/{id}/relationships?target[type]=initiative` | Relationships |
| `/initiatives/{id}/links/features` | `/v2/entities/{id}/relationships?target[type]=feature` | Relationships |
| `/initiatives/{id}/links/objectives` | `/v2/entities/{id}/relationships?target[type]=objective` | Relationships |
| `/notes/{id}/links` | `/v2/notes/{id}/relationships` | Note relationships |
| `/hierarchy-entities/custom-fields` | `/v2/entities/configurations` | Config endpoint |
| `/hierarchy-entities/custom-fields-values` | `/v2/entities/fields/{id}/values` | Field values |
| `/jira-integrations` | `/v2/jira-integrations` | Prefix only |
| `/plugin-integrations` | `/v2/plugin-integrations` | Prefix only |
| `/webhooks` | `/v2/webhooks` | Prefix only |
| `/feedback-form-configurations` | N/A | No V2 equivalent |
| `/feature-release-assignments` | V2 entity relationships | Via relationships |

## R7: V2 Filter Parameter Mapping

**Decision**: V2 uses different filter parameter names for entity queries. The client must map V1-style params to V2 equivalents.

**Rationale**: V1 uses flat params like `statusId`, `statusName`, `ownerEmail`. V2 uses nested params like `status[id]`, `status[name]`, `owner[email]`. The client path translator must also map query parameters.

Key V2 entity filters: `name`, `owner[id]`, `owner[email]`, `status[id]`, `status[name]`, `archived`, `parent[id]`, `type[]`.

## R8: Constitution Amendment Required

**Decision**: The constitution constraint "All ProductBoard API requests MUST include the `X-Version: 1` header" must be amended to allow V2 API usage.

**Rationale**: This is a direct conflict. The amendment should update the constraint to: "API version is configurable (default V2). V1 requests use `X-Version: 1` header. V2 requests use `/v2` URL prefix."

This is a MINOR version bump (1.0.0 → 1.1.0) since it adds V2 support without removing V1.

## R9: New V2-Only Commands (Members, Teams)

**Decision**: Add `members` and `teams` as new top-level commands, available only when `api_version` is `2`.

**Rationale**: These endpoints don't exist in V1. Commands should error gracefully with "This command requires API V2" when V1 is active. Follows existing CLI patterns: list, get subcommands with table/JSON output.

V2 Members: `GET /v2/members`, `GET /v2/members/{id}` — fields: id, name, email, role
V2 Teams: `GET /v2/teams`, `GET /v2/teams/{id}` — fields: id, name, handle, description
