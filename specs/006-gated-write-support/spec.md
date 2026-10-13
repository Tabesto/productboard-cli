# Feature Specification: Gated Write Support

**Feature Branch**: `006-gated-write-support`  
**Created**: 2026-10-13  
**Status**: Draft  
**Input**: User description: "Add gated write support (create/update of entities and notes) via Productboard API V2 to the pboard CLI and MCP server. Read-only by default; writes enabled by --allow-write flag, write_enabled config key, or PRODUCTBOARD_WRITE_ENABLED env. MCP write tools registered only when enabled. Scope: entities create/update, notes create/update; no deletes, no relationships, V2 only."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Write stays off unless explicitly enabled (Priority: P1)

A user (or an AI agent acting for them) runs `pboard` today with read-only expectations. After this feature ships, nothing changes by default: any attempt to create or update data is refused with a clear message explaining how to enable writes, and no request is sent to Productboard.

**Why this priority**: The gate is the safety contract of the feature. Without it, existing read-only users and agents gain write power silently.

**Independent Test**: With no flag, config key or environment variable set, run any write command; confirm it fails with the "writes disabled" message and a dedicated exit code, and that no network request was made.

**Acceptance Scenarios**:

1. **Given** writes are not enabled anywhere, **When** the user runs a create or update command, **Then** the command fails with a message naming the three ways to enable writes, and no data is sent.
2. **Given** writes are not enabled, **When** the user runs any existing read command, **Then** it behaves exactly as before.
3. **Given** writes are enabled by only one of the three mechanisms (flag, config key, or environment variable), **When** the user runs a write command, **Then** it is allowed.
4. **Given** the config enables writes, **When** the user explicitly disables it again with the configure command, **Then** write commands are refused again.

---

### User Story 2 - Create and update entities (Priority: P1)

A product manager creates a new feature, initiative or other entity, and later changes its fields (name, status, owner, custom fields), from the terminal, with writes enabled.

**Why this priority**: Entities are the core of the product hierarchy; this is the main value of write support.

**Independent Test**: With writes enabled on a test workspace, create an entity of a given type with a name, read it back with the existing read command, update one field, read it back again and see the change.

**Acceptance Scenarios**:

1. **Given** writes are enabled, **When** the user creates an entity by giving its type and field values, **Then** the entity exists in Productboard and the command prints its identifier, type and link.
2. **Given** an existing entity, **When** the user updates one or more fields, **Then** only those fields change.
3. **Given** an existing list-type field, **When** the user adds or removes items rather than replacing the whole value, **Then** only those items are added or removed.
4. **Given** the user provides both a full-replacement and an add/remove change in the same update, **Then** the command is refused before sending anything.
5. **Given** an entity type or field that the workspace does not support, **When** the user tries to create or update with it, **Then** the error returned by Productboard is shown unchanged enough to tell which type or field is wrong.

---

### User Story 3 - Create and update notes (Priority: P2)

A user captures customer feedback as a note (name, content, optional tags/owner/customer) and later edits it.

**Why this priority**: Notes are the second main write use case but depend on the same gate and mechanics as User Story 2.

**Independent Test**: With writes enabled, create a text note with a name and content, read it back, update its content, read it back.

**Acceptance Scenarios**:

1. **Given** writes are enabled, **When** the user creates a note with a type and a name, **Then** the note exists and its identifier and link are printed.
2. **Given** a note type that cannot be created through the API, **When** the user attempts it, **Then** the refusal from Productboard is surfaced clearly.
3. **Given** an existing note, **When** the user updates its fields, **Then** only those fields change.

---

### User Story 4 - AI agents see write tools only when writes are enabled (Priority: P1)

A user connects `pboard` to an AI assistant as an MCP server. By default the assistant sees only read tools. When the user starts the server with writes enabled (by flag, config or environment), the assistant additionally sees tools to create and update entities and notes.

**Why this priority**: Agents are the highest-risk write actors; hiding the tools when off is the strongest protection.

**Independent Test**: List the tools of the server started without enabling writes (no create/update tools) and again with writes enabled (four extra tools).

**Acceptance Scenarios**:

1. **Given** writes are not enabled, **When** an agent lists tools, **Then** no create or update tool is present.
2. **Given** writes are enabled, **When** an agent lists tools, **Then** create and update tools for entities and notes are present, and their descriptions tell the agent to inspect the workspace's field configuration first.
3. **Given** the server was registered with the install command and writes were requested at that time, **When** the assistant starts the server, **Then** writes are enabled for that server.

---

### User Story 5 - Turn writes on or off persistently (Priority: P2)

A user who writes regularly stores the choice in their configuration instead of passing a flag each time, and can turn it off again. Re-running configure to change the token must not lose this choice or other settings.

**Why this priority**: Convenience for frequent writers; also fixes a current defect where configure overwrites other settings.

**Independent Test**: Enable writes via configure, change the token via configure, confirm writes are still enabled and the API version setting is unchanged.

**Acceptance Scenarios**:

1. **Given** a saved configuration, **When** the user enables writes via configure, **Then** later write commands work without a flag.
2. **Given** writes enabled and an API version saved, **When** the user re-runs configure only to change the token, **Then** the write setting and API version are preserved.

---

### Edge Cases

- Writes requested while the legacy API version 1 is selected: refused with a message that writes need version 2.
- Productboard returns minimal data after create/update: the user still gets identifier, type and link; they use the existing read commands for full detail.
- Value given for a field is not valid structured data where structured data is required: refused locally with a clear message before any request.
- Invalid or expired token on write: same authentication error and exit code as reads.
- Rate limiting on write: same rate-limit error and exit code as reads.
- Write setting given as an unrecognised value in the environment variable: treated as disabled.
- Both flag and config disagree (flag given, config off): the flag enables for that single run only.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST refuse every create and update operation unless writes are explicitly enabled, and MUST NOT contact Productboard when refusing.
- **FR-002**: System MUST allow writes to be enabled by any one of: a command-line flag for the current run, a persisted configuration setting, or an environment variable.
- **FR-003**: System MUST enforce the gate in one shared place used by both the command line and the MCP server, so no write path can bypass it.
- **FR-004**: System MUST return a distinct, documented exit code and a message listing the ways to enable writes when a write is refused by the gate.
- **FR-005**: Users MUST be able to create an entity of a given type with field values, and receive its identifier, type and link.
- **FR-006**: Users MUST be able to update an existing entity by replacing field values, or by adding items, removing items, or clearing a field; both styles in one request MUST be rejected before sending.
- **FR-007**: Users MUST be able to create a note of a given type with a name and optional content and other fields, and receive its identifier and link.
- **FR-008**: Users MUST be able to update an existing note with the same replace or add/remove/clear choices as entities.
- **FR-009**: System MUST accept field values as simple `name=value` pairs (value read as structured data when valid, plain text otherwise) and also a full raw JSON payload for advanced cases.
- **FR-010**: System MUST refuse write operations when the legacy API version 1 is selected.
- **FR-011**: System MUST show errors returned by Productboard for invalid types or fields so the user can correct them.
- **FR-012**: The MCP server MUST register create/update tools for entities and notes only when writes are enabled at server start; otherwise those tools MUST NOT exist.
- **FR-013**: MCP write tool descriptions MUST instruct the agent to consult the workspace's field configuration before writing.
- **FR-014**: The MCP server install command MUST allow recording the "writes enabled" choice into the registered server command.
- **FR-015**: The configure command MUST allow enabling and disabling writes, and MUST preserve all other saved settings (token, API address, API version, write setting) when run for any purpose.
- **FR-016**: All existing read commands and read tools MUST keep their current behavior.
- **FR-017**: Documentation (README, embedded agent skill, and project constitution's read-only principle) MUST describe read-only-by-default with gated writes.
- **FR-018**: Deleting data, managing relationships, field-value option management and V1 writes MUST remain out of scope and unavailable.
- **FR-019**: System MUST append every write to a local audit log before sending it (with the record's pre-update value and a ready-to-send revert body for updated fields), record the outcome afterwards, and MUST NOT send the write if the log cannot be written or the current value cannot be read.

### Key Entities

- **Write gate**: A yes/no switch, off by default, resolved from flag, saved configuration and environment variable; applies to the whole process.
- **Entity**: A product-management object (feature, initiative, objective, release, etc.) with a type and workspace-defined fields.
- **Note**: A piece of customer feedback with a type, name, content and optional tags, owner and customer.
- **Write result**: The minimal confirmation returned after a create or update: identifier, type and link to the full record.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: With default settings, 100% of write attempts (command line and MCP) are refused and zero write requests reach Productboard.
- **SC-002**: A user can enable writes and create a first entity in under 2 minutes using only the built-in help.
- **SC-003**: An entity or note can be created and then updated, and the change verified with a read command, in 3 commands or fewer.
- **SC-004**: An AI assistant connected without write enabled sees 0 write tools; with it enabled, sees exactly the 4 documented write tools.
- **SC-005**: Re-running configure to change only the token leaves 100% of other saved settings unchanged.
- **SC-006**: Every pre-existing read command and read tool produces unchanged results after the feature ships.
- **SC-007**: For 100% of updates, the audit log holds enough information to restore the replaced field values without consulting Productboard history.

## Assumptions

- Productboard API V2 is the only supported version for writes (V1 was retired by Productboard on 2026-07-08).
- The user's access token already carries write permission in Productboard; the tool does not manage permissions.
- Field names and availability are workspace-specific; the tool does not validate them locally and relies on Productboard's error messages and its field-configuration endpoints.
- After create/update the full record is retrieved with existing read commands rather than returned automatically.
- Delete operations, relationship management and field-value option management are deferred to future features.
- The persisted setting lives in the existing user configuration file; the environment variable accepts `1` or `true` (case-insensitive) as enabled.
- A flag enables writes for the current run only and never changes saved configuration.
