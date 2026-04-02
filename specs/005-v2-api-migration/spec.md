# Feature Specification: Migrate CLI to ProductBoard V2 API

**Feature Branch**: `005-v2-api-migration`  
**Created**: 2026-04-01  
**Status**: Draft  
**Input**: User description: "there is a V2 api, see how to migrate the CLI to use it: https://developer.productboard.com/v2/openapi"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Unified Entity Browsing via V2 API (Priority: P1)

A CLI user runs existing commands (`pboard features list`, `pboard features get <id>`, etc.) and the CLI fetches data from the ProductBoard V2 API instead of V1. The user experience remains the same: identical command names, flags, and output formats. Under the hood, the HTTP client targets V2 endpoints with V2 response parsing.

**Why this priority**: This is the core migration. All existing read commands (features, products, components, initiatives, objectives, key-results, releases, release-groups, companies, users, notes) must work against V2. Without this, no other V2 work matters.

**Independent Test**: Run every existing `list` and `get` command and verify results match the same data previously returned by V1.

**Acceptance Scenarios**:

1. **Given** a configured API token with V2 scopes, **When** the user runs `pboard features list`, **Then** results are fetched from V2 `/entities` endpoint filtered by type `feature` and displayed in table/JSON format.
2. **Given** a configured API token, **When** the user runs `pboard features get <id>`, **Then** the feature detail is fetched from V2 `/entities/{id}` and displayed with the same fields as before.
3. **Given** a configured API token, **When** the user runs `pboard notes list`, **Then** notes are fetched from V2 `/notes` endpoint and displayed correctly.
4. **Given** pagination is needed, **When** the user lists resources exceeding one page, **Then** the CLI auto-paginates using V2 cursor-based pagination (`links.next` / `pageCursor`).
5. **Given** the user passes `--limit 10`, **When** listing resources, **Then** only 10 results are returned regardless of total available.

---

### User Story 2 - Relationship and Linked Resource Browsing (Priority: P2)

A CLI user browses relationships between entities (e.g., `pboard features links initiatives`, `pboard objectives links features`) and the CLI fetches these via the V2 relationships endpoint instead of V1 link endpoints.

**Why this priority**: Linked resources are a commonly used feature of the CLI. V2 uses a unified `/entities/{id}/relationships` endpoint instead of separate per-type link paths.

**Independent Test**: Run relationship commands for features, objectives, and initiatives and verify linked entities are returned correctly.

**Acceptance Scenarios**:

1. **Given** a feature with linked initiatives, **When** the user runs `pboard features links initiatives <id>`, **Then** results come from V2 `/entities/{id}/relationships` filtered by target type and display correctly.
2. **Given** a note with linked features, **When** the user runs `pboard notes links <id>`, **Then** results come from V2 `/notes/{id}/relationships` and display correctly.

---

### User Story 3 - Custom Field and Status Browsing (Priority: P2)

A CLI user lists custom fields, custom field values, and feature statuses. The CLI fetches these from V2 configuration and field-values endpoints.

**Why this priority**: Custom fields and statuses are used for filtering and understanding entity metadata. V2 restructures these under `/entities/configurations` and `/entities/fields/{id}/values`.

**Independent Test**: Run `pboard custom-fields list`, `pboard custom-fields values list`, and `pboard feature-statuses list` and verify correct output.

**Acceptance Scenarios**:

1. **Given** custom fields exist, **When** the user runs `pboard custom-fields list`, **Then** fields are fetched from V2 entity configuration endpoints.
2. **Given** a custom field with allowed values, **When** the user queries its values, **Then** values come from V2 `/entities/fields/{id}/values`.
3. **Given** feature statuses exist, **When** the user runs `pboard feature-statuses list`, **Then** statuses are retrieved from V2 configuration data.

---

### User Story 4 - Integration and Webhook Browsing (Priority: P3)

A CLI user lists Jira integrations, plugin integrations, and webhooks. These V2 endpoints have similar structure to V1 but live under `/v2` path prefix.

**Why this priority**: Integrations and webhooks are less frequently used but still need migration for completeness.

**Independent Test**: Run `pboard jira-integrations list`, `pboard plugin-integrations list`, and `pboard webhooks list` and verify results.

**Acceptance Scenarios**:

1. **Given** Jira integrations exist, **When** the user runs `pboard jira-integrations list`, **Then** results come from V2 `/jira-integrations` endpoint.
2. **Given** webhooks exist, **When** the user runs `pboard webhooks list`, **Then** results come from V2 `/webhooks` endpoint.

---

### User Story 5 - MCP Server Uses V2 API (Priority: P3)

The MCP server (used for Claude Desktop integration) exposes the same tools as the CLI. After migration, all MCP tool handlers must also use V2 endpoints, returning V2 response data.

**Why this priority**: MCP server is a secondary interface that reuses the same client. Once the client is migrated, MCP should work, but handler response mapping may need updates.

**Independent Test**: Invoke MCP tools via Claude Desktop and verify correct data is returned.

**Acceptance Scenarios**:

1. **Given** the MCP server is running, **When** a tool like `list_features` is invoked, **Then** it returns data fetched from V2 API.
2. **Given** the MCP server is running, **When** any tool is invoked, **Then** the response structure matches what Claude Desktop expects.

---

### Edge Cases

- What happens when the API token lacks V2-required scopes (e.g., `entities:read`)? The CLI should display a clear authentication error.
- What happens when a V1-only endpoint has no V2 equivalent (e.g., `feedback-form-configurations`)? The CLI should display a clear message that the command is unavailable with V2.
- What happens when V2 returns a different field structure than V1 for the same entity? The CLI output columns must still show meaningful data.
- How does the CLI handle V2 error responses (which have a different structure than V1)?
- What happens when rate-limited (429) by V2? The CLI should display rate-limit messaging consistent with current behavior.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST default to V2 API (`/v2` URL path prefix) for all API calls, while supporting V1 via a `--api-version` flag or config option (`api_version` in config file). V2 is the default when no override is specified.
- **FR-002**: System MUST map existing entity-type commands (features, products, components, initiatives, objectives, key-results, releases, release-groups) to V2 unified `/entities` endpoint with appropriate type filters.
- **FR-003**: System MUST map notes commands to V2 `/notes` endpoint with updated request/response handling.
- **FR-004**: System MUST map company and user list/get commands to V2 `/entities` endpoint with `company` and `user` type filters respectively.
- **FR-005**: System MUST handle V2 cursor-based pagination using `links.next` field and `pageCursor` query parameter.
- **FR-006**: System MUST parse V2 response structures (wrapped in `data` with `fields` sub-object) and map them to existing display columns.
- **FR-007**: System MUST map relationship/link commands to V2 `/entities/{id}/relationships` and `/notes/{id}/relationships` endpoints.
- **FR-008**: System MUST map custom field commands to V2 `/entities/configurations` and `/entities/fields/{id}/values` endpoints.
- **FR-009**: System MUST map feature-status listing to V2 entity configuration data.
- **FR-010**: System MUST map integration commands (Jira, plugin, webhooks) to their V2 equivalents.
- **FR-011**: System MUST parse V2 error responses (with `id`, `errors[]` containing `code`, `title`, `detail`) and display user-friendly messages.
- **FR-012**: System MUST preserve all existing CLI command names, flags, and output format options (table/JSON).
- **FR-013**: System MUST update the MCP server handlers to work with V2 response data.
- **FR-014**: System MUST handle commands that have no V2 equivalent (e.g., feedback-form-configurations) by displaying a deprecation or unavailability message.
- **FR-015**: System MUST support V2 release-assignment queries via entity relationships instead of the V1 dedicated endpoint.
- **FR-016**: System MUST add new `members` commands (list, get, search) using V2 `/members` endpoint. These commands are only available when using V2.
- **FR-017**: System MUST add new `teams` commands (list, get) using V2 `/teams` endpoint. These commands are only available when using V2.

### Key Entities

- **Entity**: Unified V2 concept encompassing features, products, components, initiatives, objectives, key-results, releases, release-groups, companies, and users. Each has a `type` field and configurable `fields`.
- **Note**: Separate V2 resource (`/notes`) with types: textNote, conversationNote, opportunityNote.
- **Relationship**: V2 connection between entities or notes, with types: parent, child, link, isBlockedBy, isBlocking.
- **Configuration**: V2 metadata describing available entity types, their fields, and supported operations.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: All existing CLI commands produce equivalent output when backed by V2 API as they did with V1.
- **SC-002**: Users experience no change in command syntax, flags, or output format after migration.
- **SC-003**: Pagination works correctly for result sets of any size, with auto-follow and limit support.
- **SC-004**: Error messages for authentication failures, not-found, and rate-limiting are clear and actionable.
- **SC-005**: MCP server tools return correct data from V2 API without breaking Claude Desktop integration.
- **SC-006**: Commands with no V2 equivalent display a clear message rather than failing silently.

## Clarifications

### Session 2026-04-01

- Q: Migration strategy — hard cutover to V2 or dual V1/V2 support? → A: Dual support with V2 as default. Users can switch to V1 via flag or config option.
- Q: Should the CLI add new V2-only commands (members, teams, analytics)? → A: Add members and teams commands. Skip analytics.
- Q: Should the CLI add write operations (create/update/delete) now available in V2? → A: Out of scope. CLI remains read-only for this migration. Write operations deferred to a future feature.

## Assumptions

- The existing API token format (Bearer token) is compatible with V2 API authentication.
- V2 API scopes required (e.g., `entities:read`, `notes:read`) are available on existing tokens or can be obtained by the user.
- The V2 API is stable and generally available for all ProductBoard workspaces.
- The `feedback-form-configurations` endpoint has no V2 equivalent and will be marked as unavailable.
- V2 entity types cover all resource types currently supported by V1 (feature, product, component, initiative, objective, keyResult, release, releaseGroup, company, user).
- The V1 API will eventually be deprecated, making this migration necessary.
- Release assignments can be queried through V2 entity relationships rather than a dedicated endpoint.
- Write operations (create, update, delete) are explicitly out of scope for this migration. The CLI remains read-only.
