# Feature Specification: Configure Command Token Argument

**Feature Branch**: `002-configure-token-arg`
**Created**: 2026-03-26
**Status**: Draft
**Input**: User description: "I would like the pboard configure command to accept the Productboard API token as a copy-paste argument and configure everything automatically so that the process is simple."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - One-Command Token Setup (Priority: P1)

As a user, I want to pass my API token directly as an argument to the configure command so I can set up pboard in a single step without interactive prompts.

**Why this priority**: This is the core request — eliminating the interactive prompt in favor of a simpler, paste-and-go workflow. It also enables scripting and CI/CD usage.

**Independent Test**: Can be fully tested by running `pboard configure <token>` and verifying that subsequent commands authenticate successfully.

**Acceptance Scenarios**:

1. **Given** a valid API token, **When** the user runs the configure command with the token as an argument, **Then** the token is saved to the configuration file and a success confirmation is displayed.
2. **Given** a valid API token, **When** the user runs the configure command with the token, **Then** subsequent commands (e.g., `pboard features list`) authenticate successfully using the saved token.
3. **Given** an empty or whitespace-only token argument, **When** the user runs the configure command, **Then** a clear error message is displayed indicating the token is invalid.

---

### User Story 2 - Preserve Interactive Fallback (Priority: P2)

As a user, I want the configure command to still prompt me interactively if I don't provide a token argument, so the existing workflow remains available.

**Why this priority**: Backward compatibility ensures existing users and documentation remain valid.

**Independent Test**: Can be tested by running `pboard configure` without arguments and verifying the interactive prompt still works.

**Acceptance Scenarios**:

1. **Given** no token argument is provided, **When** the user runs the configure command, **Then** an interactive prompt asks for the token (existing behavior preserved).
2. **Given** a token argument is provided, **When** the user runs the configure command, **Then** no interactive prompt is shown and the token is saved directly.

---

### Edge Cases

- What happens when the user provides a token that contains leading or trailing whitespace? The command should trim whitespace before saving.
- What happens when the user provides multiple arguments? The command should reject extra arguments with a clear usage message.
- What happens when the configuration directory does not exist yet? The command should create it automatically (existing behavior).
- What happens when the config file already exists with a different token? The command should overwrite it with the new token silently.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The configure command MUST accept an optional positional argument for the API token.
- **FR-002**: When a token argument is provided, the command MUST save it to the configuration file without prompting the user.
- **FR-003**: When no token argument is provided, the command MUST fall back to the existing interactive prompt behavior.
- **FR-004**: The command MUST trim leading and trailing whitespace from the provided token before saving.
- **FR-005**: The command MUST display a success confirmation message after saving the token, including the path to the config file.
- **FR-006**: The command MUST reject empty or whitespace-only tokens with a clear error message.
- **FR-007**: The command MUST preserve all existing configuration file behaviors (directory creation, file permissions mode 600, YAML format).

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can configure the tool in a single command with no interactive prompts when providing a token argument.
- **SC-002**: The existing interactive configuration flow continues to work identically when no argument is provided.
- **SC-003**: Users can copy a token from a web browser and paste it directly into the command in under 10 seconds.

## Assumptions

- Users have their API token available (e.g., copied from the ProductBoard web interface) before running the command.
- The existing `pboard configure` command and configuration infrastructure (config directory, file format, permissions) are already implemented and working.
- This change is backward-compatible — no existing scripts or workflows that use `pboard configure` interactively will break.
- Only one token value is expected; the command does not need to support multiple tokens or named profiles.
