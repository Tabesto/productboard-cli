# Feature Specification: MCP Server for pboard CLI

**Feature Branch**: `003-pboard-mcp-server`
**Created**: 2026-03-26
**Status**: Draft
**Input**: User description: "create an mcp server for using pboard cli tool from Claude desktop"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Browse ProductBoard Features from Claude Desktop (Priority: P1)

A product manager or developer using Claude Desktop wants to query ProductBoard features without leaving the conversation. They ask Claude to list features, get details on a specific feature, or filter features by status, and Claude uses the MCP server to retrieve the data from ProductBoard and present it directly in the chat.

**Why this priority**: Feature browsing is the most common ProductBoard interaction and the core value proposition of having an MCP server -- bringing ProductBoard data directly into the AI assistant workflow.

**Independent Test**: Can be fully tested by configuring the MCP server in Claude Desktop, asking Claude to list or get features, and verifying correct data is returned from ProductBoard.

**Acceptance Scenarios**:

1. **Given** a configured MCP server with a valid ProductBoard token, **When** the user asks Claude to list features, **Then** Claude retrieves and displays the features with their names, statuses, and IDs.
2. **Given** a configured MCP server, **When** the user asks Claude to get details for a specific feature by ID, **Then** Claude retrieves and displays the full feature details.
3. **Given** a configured MCP server, **When** the user asks Claude to filter features by status or owner, **Then** Claude applies the filters and returns matching features.

---

### User Story 2 - Query Notes and Feedback (Priority: P1)

A product manager wants to review customer feedback and notes associated with features. They ask Claude to list notes, search notes by keyword, or get notes linked to a specific feature or company.

**Why this priority**: Notes and feedback are the primary input for product decisions in ProductBoard. Accessing them from Claude Desktop enables AI-assisted analysis of customer feedback.

**Independent Test**: Can be fully tested by asking Claude to list notes with various filters (term, date range, company, feature) and verifying the returned data matches ProductBoard.

**Acceptance Scenarios**:

1. **Given** a configured MCP server, **When** the user asks Claude to search notes containing a keyword, **Then** Claude returns matching notes with their content summaries.
2. **Given** a configured MCP server, **When** the user asks Claude to list notes for a specific feature, **Then** Claude returns notes linked to that feature.
3. **Given** a configured MCP server, **When** the user asks Claude to list notes from a specific company, **Then** Claude returns notes associated with that company.

---

### User Story 3 - Explore Product Hierarchy and Releases (Priority: P2)

A team member wants to understand the product structure (products, components, objectives, initiatives) or check release status. They ask Claude to list products, components, releases, or explore OKR linkages.

**Why this priority**: Product hierarchy and release data provide context for product decisions but are queried less frequently than features and notes.

**Independent Test**: Can be fully tested by asking Claude to list products, components, releases, objectives, or initiatives and verifying accurate data retrieval.

**Acceptance Scenarios**:

1. **Given** a configured MCP server, **When** the user asks Claude to list releases, **Then** Claude returns the releases with their details.
2. **Given** a configured MCP server, **When** the user asks Claude to show initiatives linked to a feature, **Then** Claude returns the linked initiatives.
3. **Given** a configured MCP server, **When** the user asks Claude to list objectives and their linked features, **Then** Claude retrieves objectives and follows the linkage to show related features.

---

### User Story 4 - Look Up Companies and Users (Priority: P2)

A team member wants to look up company information or find users in ProductBoard. They ask Claude to list or search companies, view company custom fields, or list users.

**Why this priority**: Company and user lookups support customer-centric workflows but are secondary to core product data.

**Independent Test**: Can be fully tested by asking Claude to list companies, search by term, or get company details and custom field values.

**Acceptance Scenarios**:

1. **Given** a configured MCP server, **When** the user asks Claude to search for a company by name, **Then** Claude returns matching companies.
2. **Given** a configured MCP server, **When** the user asks Claude to get custom field values for a company, **Then** Claude returns the field values.

---

### User Story 5 - One-Command Install via pboard CLI (Priority: P1)

A user wants to set up the MCP server to work with Claude Desktop with zero manual configuration. They run `pboard mcp install` and the command automatically registers the MCP server in Claude Desktop's configuration file. After restarting Claude Desktop, the ProductBoard tools are immediately available. This follows the same pattern as the existing `pboard skill install` command.

**Why this priority**: Without easy setup, no other user story can be realized. A single-command install removes all friction and makes adoption trivial.

**Independent Test**: Can be fully tested by running `pboard mcp install`, restarting Claude Desktop, and verifying the ProductBoard tools appear and are usable.

**Acceptance Scenarios**:

1. **Given** the user has pboard CLI installed and a ProductBoard API token already configured, **When** they run `pboard mcp install`, **Then** the MCP server entry is automatically added to Claude Desktop's configuration file and the user is told to restart Claude Desktop.
2. **Given** the MCP server is already installed, **When** the user runs `pboard mcp install` again, **Then** they are informed it is already installed (or can use `--force` to overwrite).
3. **Given** the user wants to preview the changes, **When** they run `pboard mcp install --dry-run`, **Then** the command shows what configuration would be written without modifying any files.
4. **Given** the user wants to remove the MCP server, **When** they run `pboard mcp uninstall`, **Then** the MCP server entry is removed from Claude Desktop's configuration file.
5. **Given** no ProductBoard API token is configured, **When** the MCP server starts via Claude Desktop, **Then** it provides a clear error message indicating the token is missing and how to configure it via `pboard configure`.

---

### Edge Cases

- What happens when the ProductBoard API token is invalid or expired? The server should return a clear, user-friendly error message.
- What happens when the ProductBoard API rate limit is reached? The server should communicate the rate limit condition back to Claude.
- What happens when a requested entity (feature, note, company) does not exist? The server should return a clear "not found" message.
- What happens when the API returns paginated results and the user hasn't specified a limit? The server returns up to 25 results by default. Users can override via the limit parameter.
- What happens when the MCP server process crashes or loses connectivity? Claude Desktop should be able to detect the disconnection and inform the user.
- What happens when Claude Desktop's configuration file doesn't exist yet? The install command should create it with the correct structure.
- What happens when the Claude Desktop configuration file has been manually edited with other MCP servers? The install command should add the pboard entry without disturbing existing entries.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST expose each pboard CLI resource type (features, notes, products, components, releases, release-groups, objectives, key-results, initiatives, companies, users, custom-fields, feature-statuses, feature-release-assignments, plugin-integrations, jira-integrations, webhooks, feedback-forms) as one or more MCP tools.
- **FR-002**: System MUST support all filtering options available in the pboard CLI for each resource (e.g., status-name, owner-email, term, date ranges for notes, etc.).
- **FR-003**: System MUST support retrieving individual entities by ID where the pboard CLI supports it (get subcommands).
- **FR-004**: System MUST support relationship/linkage queries (e.g., features linked to an objective, initiatives linked to a feature).
- **FR-005**: System MUST return data in a structured format that Claude can interpret and present clearly to users.
- **FR-006**: System MUST reuse the existing pboard configuration (token from `~/.config/pboard/config.yaml` or `PRODUCTBOARD_API_TOKEN` environment variable).
- **FR-007**: System MUST provide meaningful error messages when the API token is missing, invalid, or when API calls fail.
- **FR-008**: System MUST support result limiting to control the volume of data returned. The default limit MUST be 25 results when no explicit limit is provided by the user.
- **FR-009**: System MUST communicate using the MCP protocol (stdio transport) so Claude Desktop can connect to it as a tool provider.
- **FR-010**: System MUST provide a `pboard mcp install` command that automatically registers the MCP server in Claude Desktop's configuration file, requiring no manual editing.
- **FR-011**: The `pboard mcp install` command MUST support `--force` to overwrite an existing installation and `--dry-run` to preview changes without writing.
- **FR-012**: System MUST provide a `pboard mcp uninstall` command that removes the MCP server entry from Claude Desktop's configuration file.
- **FR-013**: The install command MUST detect the correct Claude Desktop configuration file location for the user's operating system. The install targets Claude Desktop only (not Claude Code or other MCP clients).

### Key Entities

- **MCP Tool**: A callable function exposed by the server, with a one-to-one mapping to each pboard CLI subcommand (e.g., `list_features`, `get_feature`, `list_notes`, `get_note`). Each tool has a unique name, description, and input schema. Expect ~30+ individual tools.
- **Tool Result**: The structured response returned by a tool invocation, containing the data retrieved from ProductBoard or an error message.
- **Server Configuration**: The connection settings (transport type, command path, environment variables) needed to register the server with Claude Desktop.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can query any ProductBoard resource type supported by pboard CLI directly from Claude Desktop without switching applications.
- **SC-002**: Server setup requires a single command (`pboard mcp install`) after initial token configuration -- no manual file editing needed.
- **SC-003**: All tool invocations return results or meaningful error messages within the same response time as running the equivalent pboard CLI command.
- **SC-004**: 100% of pboard CLI read commands are accessible as MCP tools.
- **SC-005**: Users with an existing pboard CLI configuration can start using the MCP server without re-entering their API token.

## Clarifications

### Session 2026-03-26

- Q: Should `pboard mcp install` support Claude Desktop only, or also Claude Code? → A: Claude Desktop only.
- Q: Tool granularity -- one tool per CLI subcommand or consolidated? → A: One tool per CLI subcommand (~30+ individual tools).
- Q: Default result limit when no limit specified? → A: 25 results.

## Assumptions

- Users already have or can obtain a ProductBoard API token with appropriate read permissions.
- Users have Claude Desktop installed (the install command handles configuration automatically).
- The pboard CLI is installed and available on the user's system PATH (the install command references the pboard binary location).
- The MCP server will be read-only, matching the current pboard CLI capabilities (no write operations).
- The server uses stdio transport, which is the standard for local MCP servers in Claude Desktop.
- The existing Go codebase and internal packages (client, CLI) can be reused to build the MCP server.
