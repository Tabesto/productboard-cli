# Data Model: Migrate CLI to ProductBoard V2 API

**Branch**: `005-v2-api-migration` | **Date**: 2026-04-01

## Modified Entities

### Config

Existing `Config` struct gains a new field.

| Field | Type | Description | Default |
|-------|------|-------------|---------|
| APIToken | string | Bearer token for API auth | "" |
| BaseURL | string | API base URL | "https://api.productboard.com" |
| **APIVersion** | **string** | **API version: "1" or "2"** | **"2"** |

Config file (`~/.config/pboard/config.yaml`) adds optional `api_version` key. Absent = default "2".

### Client

Existing `Client` struct gains version awareness.

| Field | Type | Description |
|-------|------|-------------|
| httpClient | *http.Client | Standard HTTP client |
| baseURL | string | API base URL |
| token | string | Bearer token |
| **apiVersion** | **string** | **"1" or "2"** |

### V2 API Response Structures

These are not new Go structs — the project uses `map[string]interface{}` throughout. This documents the V2 JSON shapes the client must handle.

#### V2 Entity Response (single)
```json
{
  "data": {
    "id": "uuid",
    "type": "feature",
    "fields": {
      "name": "...",
      "description": "...",
      "status": {"id": "...", "name": "..."},
      "owner": {"id": "...", "email": "..."},
      "archived": false
    },
    "links": {"self": "..."}
  }
}
```

#### V2 Entity Response (list)
```json
{
  "data": [
    {"id": "uuid", "type": "feature", "fields": {...}, "links": {...}}
  ],
  "links": {"next": "...?pageCursor=abc"}
}
```

#### V2 Relationship Response
```json
{
  "data": [
    {
      "type": "link",
      "source": {"id": "uuid", "type": "feature"},
      "target": {"id": "uuid", "type": "initiative"}
    }
  ],
  "links": {"next": "..."}
}
```

#### V2 Error Response
```json
{
  "id": "error-uuid",
  "errors": [
    {"code": "FORBIDDEN", "title": "Forbidden", "detail": "Missing scope entities:read"}
  ]
}
```

## New Entities

### Member (V2 only)

| Field | Type | Description |
|-------|------|-------------|
| id | string (UUID) | Member identifier |
| type | string | Always "member" |
| fields.name | string | Display name |
| fields.email | string | Email (requires `members:pii:read` scope) |
| fields.role | string | admin, maker, viewer, contributor |
| fields.disabled | bool | Account disabled |

### Team (V2 only)

| Field | Type | Description |
|-------|------|-------------|
| id | string (UUID) | Team identifier |
| type | string | Always "team" |
| fields.name | string | Team name |
| fields.handle | string | Lowercase identifier |
| fields.description | string | Team description |

## V2 Response Flattening

To maintain backward compatibility with existing display code (`SafeStr`, `SafeNested`), the V2 client must flatten entity responses:

**Before flattening** (raw V2):
```
data.id = "uuid"
data.type = "feature"
data.fields.name = "My Feature"
data.fields.status.name = "In Progress"
```

**After flattening** (what CLI commands see):
```
id = "uuid"
type = "feature"
name = "My Feature"
status.name = "In Progress"
```

This flattening happens in the client's response parsing methods (`GetSingle`, `GetList`) when `apiVersion == "2"`.

## V2 Pagination Extraction

V1 pagination: `{"pageCursor": "abc"}` at response root.
V2 pagination: `{"links": {"next": "https://api.productboard.com/v2/entities?pageCursor=abc"}}`.

The client extracts `pageCursor` from the `links.next` URL query parameter when in V2 mode.
