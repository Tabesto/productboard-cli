# MCP Tool Contracts: Feature Health

## `features_health_list`

List features with their health update information, with optional filters.

### Parameters

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| updated_since | string | No | - | Show features with health updated on or after this date (YYYY-MM-DD) |
| updated_before | string | No | - | Show features with health updated before this date (YYYY-MM-DD) |
| status | string | No | - | Filter by feature status name |
| owner | string | No | - | Filter by feature owner email |
| health_status | string | No | - | Filter by health status (on-track, at-risk, off-track) |
| include_archived | boolean | No | false | Include archived features |
| include_no_health | boolean | No | false | Include features without health updates |
| limit | number | No | 25 | Maximum number of results |

### Response

JSON array of feature objects with `lastHealthUpdate` field, sorted by health update date (most recent first).

```json
[
  {
    "id": "...",
    "name": "...",
    "status": { "name": "..." },
    "owner": { "email": "..." },
    "archived": false,
    "lastHealthUpdate": {
      "status": "on-track",
      "updatedAt": "2025-11-13T10:00:00Z",
      "message": "..."
    }
  }
]
```

---

## `features_health_get`

Get the full health update details for a specific feature.

### Parameters

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| id | string | Yes | - | Feature ID (UUID) |

### Response

Single feature object with `lastHealthUpdate` field.

```json
{
  "id": "...",
  "name": "...",
  "status": { "name": "..." },
  "owner": { "email": "..." },
  "lastHealthUpdate": {
    "status": "on-track",
    "updatedAt": "2025-11-13T10:00:00Z",
    "message": "..."
  }
}
```

Returns error message if feature not found or has no health update.
