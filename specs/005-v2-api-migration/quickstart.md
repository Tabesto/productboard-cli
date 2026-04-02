# Quickstart: V2 API Migration Implementation

**Branch**: `005-v2-api-migration` | **Date**: 2026-04-01

## Prerequisites

- Go 1.22+ installed
- ProductBoard API token with V2 scopes (`entities:read`, `notes:read`, `members:read`, `teams:read`)
- Existing codebase checked out on `005-v2-api-migration` branch

## Implementation Order

### Step 1: Config + Client Foundation

1. Add `APIVersion` field to `internal/config/config.go` (`Config` struct)
2. Read `api_version` from config file and `PRODUCTBOARD_API_VERSION` env var
3. Add `apiVersion` field to `internal/client/client.go` (`Client` struct)
4. Add `--api-version` persistent flag to `internal/cli/root.go`

### Step 2: V2 Request Building

1. In `Client.Get()`, branch on `apiVersion`:
   - V1: keep existing `X-Version: 1` header + direct path
   - V2: prepend `/v2` to path, omit version header
2. Add path translation for entity types (e.g., `/features` → `/v2/entities?type[]=feature`)
3. Add query parameter mapping (e.g., `statusName` → `status[name]`)

### Step 3: V2 Response Parsing

1. In `GetSingle()` and `GetList()`, detect V2 response shape (`data.fields` exists)
2. Flatten `fields` sub-object to top level for display compatibility
3. Extract pagination cursor from `links.next` URL instead of root `pageCursor`

### Step 4: V2 Error Handling

1. Update `NewAPIError()` to attempt parsing V2 error format
2. Extract `errors[0].detail` for user-facing message
3. Fall back to existing behavior if V2 parsing fails

### Step 5: New Commands

1. Add `internal/cli/members.go` with list/get subcommands
2. Add `internal/cli/teams.go` with list/get subcommands
3. Add V2-only guard that errors when `apiVersion == "1"`
4. Register in `root.go`

### Step 6: Deprecated Command Handling

1. Add V2 guard to `feedback-forms` command
2. Error with message pointing to `--api-version 1`

### Step 7: MCP Server Update

1. Update MCP handlers to pass API version through
2. Verify all tool handlers work with V2 response shapes

### Step 8: Constitution Amendment

1. Update `.specify/memory/constitution.md` API versioning constraint
2. Bump version to 1.1.0

## Verification

```bash
# Build
go build -o pboard ./cmd/pboard

# Test V2 (default)
./pboard features list --limit 5
./pboard features get <id>
./pboard notes list --limit 5
./pboard members list --limit 5
./pboard teams list --limit 5

# Test V1 fallback
./pboard features list --api-version 1 --limit 5
./pboard feedback-forms list --api-version 1

# Test V2 error for deprecated command
./pboard feedback-forms list
# Expected: Error: The 'feedback-forms' command is not available with API V2.

# Test V1 error for new command
./pboard members list --api-version 1
# Expected: Error: The 'members' command requires API V2.
```
