# Data Model: ProductBoard CLI Read-Only Tool

**Date**: 2026-03-26
**Branch**: `001-productboard-cli-read`

> This data model describes the entities as consumed by the CLI from the ProductBoard API. All entities are read-only.

## Core Entities

### Feature
| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Unique identifier |
| name | string | Feature name |
| description | string | Feature description (may contain HTML) |
| status | Status (embedded) | Current feature status (id + name) |
| parent | Reference (id) | Parent feature or component |
| owner | Reference (email) | Feature owner |
| archived | boolean | Whether feature is archived |
| created_at | datetime | Creation timestamp |
| updated_at | datetime | Last update timestamp |

**Relationships**: Feature → many Initiatives (via `/features/{id}/links/initiatives`), Feature → many Objectives (via `/features/{id}/links/objectives`)

### Product
| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Unique identifier |
| name | string | Product name |
| description | string | Product description |

### Component
| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Unique identifier |
| name | string | Component name |
| description | string | Component description |
| parent | Reference (id) | Parent product or component |

### Feature Status
| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Unique identifier |
| name | string | Status name (e.g., "In Progress", "Done") |

### Note
| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Unique identifier |
| title | string | Note title |
| content | string | Note content (may contain HTML) |
| created_at | datetime | Creation timestamp |
| updated_at | datetime | Last update timestamp |
| source | string | Origin of the note |
| owner | Reference (email) | Note owner |
| company | Reference (id) | Associated company |
| tags | Tag[] | Associated tags (via `/notes/{id}/tags`) |
| links | Link[] | Associated links (via `/notes/{id}/links`) |

### Company
| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Unique identifier |
| name | string | Company name |
| domain | string | Company domain |
| custom_fields | CustomFieldValue[] | Custom field values |

### User
| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Unique identifier |
| email | string | User email |
| name | string | User display name |

### Release
| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Unique identifier |
| name | string | Release name |
| description | string | Release description |
| timeframe | Timeframe | Start and end dates |
| release_group | Reference (id) | Parent release group |
| state | string | Release state |

### Release Group
| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Unique identifier |
| name | string | Release group name |
| description | string | Release group description |

### Feature Release Assignment
| Field | Type | Description |
|-------|------|-------------|
| feature | Reference (id) | Assigned feature |
| release | Reference (id) | Target release |

### Objective
| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Unique identifier |
| name | string | Objective name |
| description | string | Objective description |
| state | string | Current state |
| owner | Reference | Objective owner |

**Relationships**: Objective → many Features (via `/objectives/{id}/links/features`), Objective → many Initiatives (via `/objectives/{id}/links/initiatives`)

### Key Result
| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Unique identifier |
| name | string | Key result name |
| description | string | Key result description |
| current_value | number | Current progress value |
| target_value | number | Target value |
| objective | Reference (id) | Parent objective |

### Initiative
| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Unique identifier |
| name | string | Initiative name |
| description | string | Initiative description |
| state | string | Current state |
| owner | Reference | Initiative owner |

**Relationships**: Initiative → many Objectives (via `/initiatives/{id}/links/objectives`), Initiative → many Features (via `/initiatives/{id}/links/features`)

### Custom Field
| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Unique identifier |
| name | string | Field name |
| description | string | Field description |
| type | enum | text, number, dropdown, member |
| options | Option[] | Dropdown options (if type=dropdown) |

### Custom Field Value
| Field | Type | Description |
|-------|------|-------------|
| custom_field | Reference (id) | The custom field definition |
| hierarchy_entity | Reference (id) | The entity this value belongs to |
| value | any | The field value (type depends on custom field type) |

### Plugin Integration
| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Unique identifier |
| name | string | Integration name |
| connections | Connection[] | Feature connections (via `/plugin-integrations/{id}/connections`) |

### Jira Integration
| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Unique identifier |
| name | string | Integration name |
| connections | Connection[] | Feature connections (via `/jira-integrations/{id}/connections`) |

### Webhook
| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Unique identifier |
| url | string | Webhook URL |
| events | string[] | Subscribed event types |

### Feedback Form Configuration
| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Unique identifier |
| name | string | Form name |
| configuration | object | Form configuration details |

## Supporting Types

### Reference
| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Referenced entity ID |

### Tag
| Field | Type | Description |
|-------|------|-------------|
| name | string | Tag name |

### Link
| Field | Type | Description |
|-------|------|-------------|
| type | string | Link type |
| target | Reference | Linked entity |

### Timeframe
| Field | Type | Description |
|-------|------|-------------|
| start_date | date | Start date |
| end_date | date | End date |

## Pagination Model

All `list` endpoints return paginated responses:

| Field | Type | Description |
|-------|------|-------------|
| data | Entity[] | Array of entities for this page |
| pageCursor | string | Cursor for next page (null if last page) |
| totalResults | integer | Total count of results (when available) |
