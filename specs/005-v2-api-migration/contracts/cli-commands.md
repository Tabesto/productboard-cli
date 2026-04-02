# CLI Command Contracts: V2 API Migration

**Branch**: `005-v2-api-migration` | **Date**: 2026-04-01

## Global Changes

### New Persistent Flag

```
--api-version string   API version to use: 1 or 2 (default "2")
```

### Config File Addition

```yaml
# ~/.config/pboard/config.yaml
api_token: "pb_..."
api_url: "https://api.productboard.com"
api_version: "2"    # NEW: "1" or "2", default "2"
```

Precedence: `--api-version` flag > `PRODUCTBOARD_API_VERSION` env var > config file > default "2".

## Existing Commands (unchanged interface, different backend)

All existing commands preserve their flags, arguments, and output columns. Only the underlying API calls change when `api_version=2`.

### Example: `pboard features list`

**V1 backend**: `GET /features?statusName=...` with `X-Version: 1` header
**V2 backend**: `GET /v2/entities?type[]=feature&status[name]=...`
**Output columns**: ID, Name, Status, Owner, Archived (unchanged)

### V1-to-V2 Filter Mapping

| V1 Parameter | V2 Parameter |
|-------------|-------------|
| `statusId` | `status[id]` |
| `statusName` | `status[name]` |
| `parentId` | `parent[id]` |
| `ownerEmail` | `owner[email]` |
| `archived` | `archived` |
| `noteId` | *(not supported in V2 entity query)* |
| `term` (notes) | *(use POST /v2/notes/search)* |
| `featureId` (notes) | *(use relationships)* |
| `companyId` (notes) | *(use relationships)* |

## New Commands (V2 only)

### `pboard members`

```
pboard members list [--role <role>] [--query <search>] [--limit <n>] [--output table|json]
pboard members get <id> [--output table|json]
```

**List output columns**: ID, Name, Email, Role
**Get output**: ID, Name, Email, Role, Disabled

Error when `api_version=1`: "Error: The 'members' command requires API V2. Use --api-version 2 or update your config."

### `pboard teams`

```
pboard teams list [--query <search>] [--limit <n>] [--output table|json]
pboard teams get <id> [--output table|json]
```

**List output columns**: ID, Name, Handle, Description
**Get output**: ID, Name, Handle, Description

Error when `api_version=1`: "Error: The 'teams' command requires API V2. Use --api-version 2 or update your config."

## Deprecated Commands

### `pboard feedback-forms`

When `api_version=2`:
```
Error: The 'feedback-forms' command is not available with API V2. Use --api-version 1 to access this command.
```

When `api_version=1`: Works as before.

## Error Response Contract

### V2 Error Display

V2 errors are parsed from `errors[0].detail` (falling back to `errors[0].title`, then raw body).

Example outputs:
```
Error: Authentication failed. Missing required scope: entities:read
Error: Resource not found.
Error: Rate limit exceeded. Wait and retry, or reduce request frequency.
```
