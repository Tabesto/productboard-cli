# Data Model: Feature Health Check Helper

**Date**: 2026-03-26 | **Branch**: `004-feature-health-helper`

## Entities

### Feature (existing -- no changes)

The existing `map[string]interface{}` representation returned by `client.GetList("/features", ...)`. Relevant fields for health commands:

| Field | Type | Source | Notes |
|-------|------|--------|-------|
| id | string | API | UUID |
| name | string | API | Feature name |
| status.name | string (nested) | API | e.g., "In Progress", "Backlog", "Live" |
| owner.email | string (nested) | API | Owner email address |
| archived | bool | API | Archived flag |
| createdAt | string (ISO 8601) | API | Creation timestamp |
| updatedAt | string (ISO 8601) | API | Last modification timestamp |
| lastHealthUpdate | object or null | API | Health update data (see below) |

### lastHealthUpdate (new field to extract)

Nested object within a Feature. May be `null` if no health update has been posted.

| Field | Type | Notes |
|-------|------|-------|
| status | string | Enum: "on-track", "at-risk", "off-track" |
| updatedAt | string (ISO 8601) | When the health update was last modified |
| message | string | Free-text field; may contain progress/problems/plans by convention |

### Health Summary (derived -- for display only)

A flattened view combining Feature metadata with Health Update data for table output. Not persisted; constructed in-memory during command execution.

| Column | Source |
|--------|--------|
| Feature Name | `feature.name` |
| Status | `feature.status.name` |
| Owner | `feature.owner.email` |
| Health | `feature.lastHealthUpdate.status` |
| Health Updated | `feature.lastHealthUpdate.updatedAt` |
| Message (truncated) | `feature.lastHealthUpdate.message` (truncated to 50 chars for table) |

## Relationships

```
Feature 1 ── 0..1 lastHealthUpdate
```

A feature has zero or one health update. Features without a health update have `lastHealthUpdate: null`.

## Validation Rules

- Date filter values (`--updated-since`, `--updated-before`) must parse as `YYYY-MM-DD` format. Invalid dates produce exit code 5 (InvalidInput).
- `--health-status` values are matched case-insensitively against `lastHealthUpdate.status`.
- `--owner` is matched as an exact string against `owner.email`.
- `--status` is matched case-insensitively against `status.name`.

## State Transitions

N/A -- this feature is read-only. No state is modified.
