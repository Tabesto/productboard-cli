# Feature Specification: Feature Health Check Helper

**Feature Branch**: `004-feature-health-helper`
**Created**: 2026-03-26
**Status**: Draft
**Input**: User description: "Create helper commands that aggregate multiple API calls to retrieve specific entities. Starting with a feature health check helper that retrieves health updates across features with filtering by last updated date, status, owner, and other feature attributes."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - List Health Updates Across All Features (Priority: P1)

As a product manager, I want to quickly see the latest health updates for all my active features in a single command, so I can get a quick overview of project health without manually checking each feature.

**Why this priority**: This is the core value proposition -- replacing multiple manual API calls with a single aggregated view. Without this, the feature has no purpose.

**Independent Test**: Can be fully tested by running the health command with no filters and verifying it displays a table of features with their health status, health update date, progress, problems, and plans.

**Acceptance Scenarios**:

1. **Given** the user has a valid ProductBoard token configured, **When** they run the health list command with no filters, **Then** the system displays all non-archived features that have a health update, showing: feature name, status, owner, health status (e.g., "On Track", "At Risk", "Off Track"), health update date, and a summary of the latest health message.
2. **Given** the user has a valid token configured, **When** they run the health list command, **Then** features without any health update are excluded from the output by default.
3. **Given** the user has a valid token configured, **When** they run the health list command with the `--include-no-health` flag, **Then** features without health updates are also shown, with an empty/placeholder health status.

---

### User Story 2 - Filter Health Updates by Date (Priority: P1)

As a product manager, I want to filter health updates by date range so I can focus on features that were recently updated or find stale features that haven't been updated in a while.

**Why this priority**: Date filtering is the primary use case mentioned -- finding recently updated features and identifying stale health checks. Equal priority to the base listing since the list alone may be too large to be useful.

**Independent Test**: Can be tested by running the health command with `--updated-since` and `--updated-before` flags and verifying only features matching the date range appear.

**Acceptance Scenarios**:

1. **Given** features with various health update dates exist, **When** the user runs the command with `--updated-since 2026-01-01`, **Then** only features with a health update date on or after 2026-01-01 are shown.
2. **Given** features with various health update dates exist, **When** the user runs the command with `--updated-before 2026-01-01`, **Then** only features with a health update date before 2026-01-01 are shown.
3. **Given** features with various health update dates exist, **When** the user combines `--updated-since` and `--updated-before`, **Then** only features within the specified date range are shown.
4. **Given** the user provides an invalid date format, **When** they run the command, **Then** a clear error message is displayed explaining the expected date format (YYYY-MM-DD).

---

### User Story 3 - Filter Health Updates by Feature Attributes (Priority: P2)

As a product manager, I want to filter the health overview by feature attributes like status, owner, or health status so I can focus on specific segments of my feature portfolio.

**Why this priority**: Attribute filtering is valuable for large boards but secondary to the core list + date filtering. The command is still useful without these filters.

**Independent Test**: Can be tested by running the health command with `--status`, `--owner`, or `--health-status` flags and verifying only matching features appear.

**Acceptance Scenarios**:

1. **Given** features with various statuses exist, **When** the user runs the command with `--status "In Progress"`, **Then** only features with that status are shown.
2. **Given** features with various owners exist, **When** the user runs the command with `--owner "amel.meradi@deliverect.com"`, **Then** only features owned by that person are shown.
3. **Given** features with various health statuses exist, **When** the user runs the command with `--health-status "at-risk"`, **Then** only features whose latest health update has that status are shown.
4. **Given** the user combines multiple filters, **When** they run the command, **Then** all filters are applied together (AND logic).

---

### User Story 4 - View Detailed Health Update for a Single Feature (Priority: P2)

As a product manager, I want to get the full health update details for a specific feature, including progress, problems, and plans sections, so I can understand the current state without navigating to the ProductBoard UI.

**Why this priority**: Complements the list view by providing drill-down capability. Useful but the list view alone already delivers significant value.

**Independent Test**: Can be tested by running the health get command with a feature ID and verifying it displays the full health update content.

**Acceptance Scenarios**:

1. **Given** a feature ID with a health update, **When** the user runs the health get command with that ID, **Then** the system displays: feature name, status, owner, health status, health update date, and the full health message text.
2. **Given** a feature ID without a health update, **When** the user runs the health get command, **Then** a message indicates no health update exists for this feature.
3. **Given** an invalid or non-existent feature ID, **When** the user runs the health get command, **Then** a clear error message is displayed.

---

### User Story 5 - Output in JSON Format (Priority: P3)

As a developer or automation script, I want to get health data in JSON format so I can pipe it to other tools for further processing or integration.

**Why this priority**: Follows the existing CLI pattern where all commands support `--output json`. Important for scriptability but not the primary interactive use case.

**Independent Test**: Can be tested by running any health command with `-o json` and verifying valid JSON output.

**Acceptance Scenarios**:

1. **Given** any health command, **When** the user adds `-o json` flag, **Then** the output is valid JSON matching the displayed data.
2. **Given** the health list command with filters and `-o json`, **When** the user runs it, **Then** the JSON output contains the same filtered results as the table view.

---

### Edge Cases

- What happens when the ProductBoard API rate limit is hit while paginating through features? The command should handle rate-limit errors gracefully and inform the user.
- What happens when a feature has a health update with empty/null progress, problems, or plans fields? The command should display those as empty rather than crashing.
- What happens when the user has no features at all? The command should display an empty result with an informational message.
- What happens when date filter values result in zero matching features? The command should display an empty result, not an error.
- What happens when the `lastHealthUpdate` field structure changes or is missing from the API response? The command should degrade gracefully.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST provide a `pboard features health` subcommand group for health-related operations.
- **FR-002**: System MUST provide a `pboard features health list` command that fetches all features and displays their latest health update information in a consolidated view.
- **FR-003**: System MUST support `--updated-since <YYYY-MM-DD>` flag to filter features by health update date (on or after the given date).
- **FR-004**: System MUST support `--updated-before <YYYY-MM-DD>` flag to filter features by health update date (before the given date).
- **FR-005**: System MUST support `--status <status-name>` flag to filter by feature status.
- **FR-006**: System MUST support `--owner <email>` flag to filter by feature owner email.
- **FR-007**: System MUST support `--health-status <status>` flag to filter by health update status (e.g., "on-track", "at-risk", "off-track").
- **FR-008**: System MUST exclude archived features from results by default.
- **FR-009**: System MUST support `--include-archived` flag to include archived features.
- **FR-010**: System MUST exclude features without health updates by default.
- **FR-011**: System MUST support `--include-no-health` flag to include features without health updates.
- **FR-012**: System MUST provide a `pboard features health get <feature-id>` command that displays the full health update details for a single feature.
- **FR-013**: System MUST support both table and JSON output formats, consistent with existing CLI commands.
- **FR-014**: System MUST sort results by health update date (most recent first) by default.
- **FR-015**: System MUST respect the existing `--limit` global flag for controlling result count.
- **FR-016**: All filters MUST combine with AND logic when multiple are specified.
- **FR-017**: System MUST expose `features_health_list` and `features_health_get` as MCP tools with parameters matching the CLI flags, consistent with the existing MCP server pattern.

### Key Entities

- **Feature**: A ProductBoard feature with attributes: id, name, status, owner, archived, createdAt, updatedAt.
- **Health Update**: The latest health status of a feature, containing: status (on-track/at-risk/off-track), updatedAt, and a single free-text message field. Authors may follow a progress/problems/plans convention within the message, but the API does not enforce or separate these sections.
- **Health Summary**: An aggregated view combining Feature metadata with its latest Health Update for display purposes.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can retrieve a health overview of all active features in a single command invocation, replacing the need for multiple manual lookups.
- **SC-002**: Users can identify features with stale health updates (e.g., older than 30 days) in under 10 seconds using date filters.
- **SC-003**: Users can filter the health overview by any combination of status, owner, health status, and date range.
- **SC-004**: The command output is parseable by scripts when using JSON output mode, enabling automation workflows.
- **SC-005**: The health list command completes within a reasonable time for boards with up to 200 features.

## Clarifications

### Session 2026-03-26

- Q: Does the features list API endpoint include `lastHealthUpdate` data, or only the individual GET endpoint? → A: List endpoint includes `lastHealthUpdate` -- single paginated call is sufficient.
- Q: Is `lastHealthUpdate.message` a single free-text field or structured sub-fields (progress/problems/plans)? → A: Single `message` text field -- progress/problems/plans are an authoring convention, not separate API fields.
- Q: Should health commands be exposed as MCP tools in addition to CLI? → A: Yes -- expose `health list` and `health get` as MCP tools, matching CLI flags.

## Assumptions

- The ProductBoard API returns `lastHealthUpdate` data as part of both the features list and individual feature GET responses (confirmed -- the field includes status, updatedAt, and message). A single paginated list call is sufficient to retrieve health data for all features.
- Health update statuses follow a known set of values (e.g., "on-track", "at-risk", "off-track") that can be filtered against.
- The existing API client pagination mechanism is sufficient to fetch all features without changes.
- The health update message is a single free-text field. Progress/problems/plans structure within it is an authoring convention, not enforced by the API. The command displays the message as-is without parsing sub-sections.
- This is the first "helper" command, establishing a pattern for future aggregation commands under the existing command tree rather than a new top-level command.
- Date filtering is performed client-side after fetching features, since the ProductBoard API does not support server-side sorting or filtering by health update date.
