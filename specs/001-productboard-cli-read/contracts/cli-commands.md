# CLI Command Contract: pboard

**Date**: 2026-03-26

## Global Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--output` | `-o` | string | `table` | Output format: `table` or `json` |
| `--limit` | `-l` | int | 0 (all) | Maximum number of results to return |
| `--help` | `-h` | bool | false | Show help for command |
| `--version` | `-v` | bool | false | Show version |

## Commands

### `pboard configure`
Interactive setup for API token. Writes to `~/.config/pboard/config.yaml` with mode 600.

```
pboard configure
```

---

### `pboard features`

| Subcommand | Description | Endpoint |
|------------|-------------|----------|
| `list` | List all features | `GET /features` |
| `get <id>` | Retrieve a feature by ID | `GET /features/{id}` |
| `links initiatives <id>` | List initiatives linked to a feature | `GET /features/{id}/links/initiatives` |
| `links objectives <id>` | List objectives linked to a feature | `GET /features/{id}/links/objectives` |

**List filters**:

| Flag | Type | Description |
|------|------|-------------|
| `--status-id` | UUID | Filter by status ID |
| `--status-name` | string | Filter by status name |
| `--parent-id` | UUID | Filter by parent ID |
| `--archived` | bool | Filter by archived state |
| `--owner-email` | string | Filter by owner email |
| `--note-id` | UUID | Filter by associated note ID |

---

### `pboard products`

| Subcommand | Description | Endpoint |
|------------|-------------|----------|
| `list` | List all products | `GET /products` |
| `get <id>` | Retrieve a product by ID | `GET /products/{id}` |

---

### `pboard components`

| Subcommand | Description | Endpoint |
|------------|-------------|----------|
| `list` | List all components | `GET /components` |
| `get <id>` | Retrieve a component by ID | `GET /components/{id}` |

---

### `pboard feature-statuses`

| Subcommand | Description | Endpoint |
|------------|-------------|----------|
| `list` | List all feature statuses | `GET /feature-statuses` |

---

### `pboard notes`

| Subcommand | Description | Endpoint |
|------------|-------------|----------|
| `list` | List all notes | `GET /notes` |
| `get <id>` | Retrieve a note by ID | `GET /notes/{id}` |
| `tags <noteId>` | List tags for a note | `GET /notes/{noteId}/tags` |
| `links <noteId>` | List links for a note | `GET /notes/{noteId}/links` |

**List filters**:

| Flag | Type | Description |
|------|------|-------------|
| `--date-from` | date | Filter notes from this date |
| `--date-to` | date | Filter notes to this date |
| `--created-from` | date | Filter by creation date start |
| `--created-to` | date | Filter by creation date end |
| `--updated-from` | date | Filter by update date start |
| `--updated-to` | date | Filter by update date end |
| `--term` | string | Search term |
| `--feature-id` | UUID | Filter by feature ID |
| `--company-id` | UUID | Filter by company ID |
| `--owner-email` | string | Filter by owner email |
| `--source` | string | Filter by source |
| `--any-tag` | string[] | Filter by any of these tags |
| `--all-tags` | string[] | Filter by all of these tags |

---

### `pboard feedback-forms`

| Subcommand | Description | Endpoint |
|------------|-------------|----------|
| `list` | List all feedback form configurations | `GET /feedback-form-configurations` |
| `get <id>` | Retrieve a feedback form configuration | `GET /feedback-form-configurations/{id}` |

---

### `pboard companies`

| Subcommand | Description | Endpoint |
|------------|-------------|----------|
| `list` | List all companies | `GET /companies` |
| `get <id>` | Retrieve a company | `GET /companies/{id}` |
| `custom-fields list` | List company custom fields | `GET /companies/custom-fields` |
| `custom-fields get <id>` | Retrieve a company custom field | `GET /companies/custom-fields/{id}` |
| `custom-field-value <companyId> <fieldId>` | Get custom field value | `GET /companies/{companyId}/custom-fields/{fieldId}/value` |

**List filters**:

| Flag | Type | Description |
|------|------|-------------|
| `--term` | string | Search term |
| `--has-notes` | string | Filter by note presence |
| `--feature-id` | string | Filter by feature ID |

---

### `pboard users`

| Subcommand | Description | Endpoint |
|------------|-------------|----------|
| `list` | List all users | `GET /users` |
| `get <id>` | Retrieve a user | `GET /users/{id}` |

---

### `pboard releases`

| Subcommand | Description | Endpoint |
|------------|-------------|----------|
| `list` | List all releases | `GET /releases` |
| `get <id>` | Retrieve a release | `GET /releases/{id}` |

**List filters**:

| Flag | Type | Description |
|------|------|-------------|
| `--release-group-id` | string | Filter by release group ID |

---

### `pboard release-groups`

| Subcommand | Description | Endpoint |
|------------|-------------|----------|
| `list` | List all release groups | `GET /release-groups` |
| `get <id>` | Retrieve a release group | `GET /release-groups/{id}` |

---

### `pboard feature-release-assignments`

| Subcommand | Description | Endpoint |
|------------|-------------|----------|
| `list` | List all assignments | `GET /feature-release-assignments` |
| `get` | Retrieve a specific assignment | `GET /feature-release-assignments/assignment` |

**List filters**:

| Flag | Type | Description |
|------|------|-------------|
| `--feature-id` | string | Filter by feature ID |
| `--release-id` | string | Filter by release ID |
| `--release-state` | string | Filter by release state |
| `--end-date-from` | date | Filter by end date start |
| `--end-date-to` | date | Filter by end date end |

**Get required flags**:

| Flag | Type | Description |
|------|------|-------------|
| `--feature-id` | UUID | Feature ID (required) |
| `--release-id` | UUID | Release ID (required) |

---

### `pboard objectives`

| Subcommand | Description | Endpoint |
|------------|-------------|----------|
| `list` | List all objectives | `GET /objectives` |
| `get <id>` | Get a specific objective | `GET /objectives/{id}` |
| `links features <id>` | List features linked to objective | `GET /objectives/{id}/links/features` |
| `links initiatives <id>` | List initiatives linked to objective | `GET /objectives/{id}/links/initiatives` |

---

### `pboard key-results`

| Subcommand | Description | Endpoint |
|------------|-------------|----------|
| `list` | List all key results | `GET /key-results` |
| `get <id>` | Retrieve a key result | `GET /key-results/{id}` |

---

### `pboard initiatives`

| Subcommand | Description | Endpoint |
|------------|-------------|----------|
| `list` | List all initiatives | `GET /initiatives` |
| `get <id>` | Get a specific initiative | `GET /initiatives/{id}` |
| `links objectives <id>` | List objectives linked to initiative | `GET /initiatives/{id}/links/objectives` |
| `links features <id>` | List features linked to initiative | `GET /initiatives/{id}/links/features` |

---

### `pboard custom-fields`

| Subcommand | Description | Endpoint |
|------------|-------------|----------|
| `list` | List all custom fields | `GET /hierarchy-entities/custom-fields` |
| `get <id>` | Retrieve a custom field | `GET /hierarchy-entities/custom-fields/{id}` |
| `values list` | List all custom field values | `GET /hierarchy-entities/custom-fields-values` |
| `values get` | Retrieve a custom field value | `GET /hierarchy-entities/custom-fields-values/value` |

**List filters (custom-fields list)**:

| Flag | Type | Description |
|------|------|-------------|
| `--type` | string[] | Filter by type (required): text, number, dropdown, member |

**List filters (values list)**:

| Flag | Type | Description |
|------|------|-------------|
| `--type` | string[] | Filter by type |
| `--custom-field-id` | UUID | Filter by custom field ID |
| `--hierarchy-entity-id` | UUID | Filter by entity ID |

**Get required flags (values get)**:

| Flag | Type | Description |
|------|------|-------------|
| `--custom-field-id` | UUID | Custom field ID (required) |
| `--hierarchy-entity-id` | UUID | Entity ID (required) |

---

### `pboard plugin-integrations`

| Subcommand | Description | Endpoint |
|------------|-------------|----------|
| `list` | List all plugin integrations | `GET /plugin-integrations` |
| `get <id>` | Retrieve a plugin integration | `GET /plugin-integrations/{id}` |
| `connections list <id>` | List connections for integration | `GET /plugin-integrations/{id}/connections` |
| `connections get <id> <featureId>` | Get a specific connection | `GET /plugin-integrations/{id}/connections/{featureId}` |

---

### `pboard jira-integrations`

| Subcommand | Description | Endpoint |
|------------|-------------|----------|
| `list` | List all Jira integrations | `GET /jira-integrations` |
| `get <id>` | Retrieve a Jira integration | `GET /jira-integrations/{id}` |
| `connections list <id>` | List connections | `GET /jira-integrations/{id}/connections` |
| `connections get <id> <featureId>` | Get a specific connection | `GET /jira-integrations/{id}/connections/{featureId}` |

**Connections list filters**:

| Flag | Type | Description |
|------|------|-------------|
| `--issue-key` | string | Filter by Jira issue key |
| `--issue-id` | string | Filter by Jira issue ID |

---

### `pboard webhooks`

| Subcommand | Description | Endpoint |
|------------|-------------|----------|
| `list` | List all webhook subscriptions | `GET /webhooks` |
| `get <id>` | Retrieve a webhook subscription | `GET /webhooks/{id}` |

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General error (network, unexpected response) |
| 2 | Authentication error (missing/invalid token) |
| 3 | Resource not found |
| 4 | Rate limited |
| 5 | Invalid input (bad flags, missing required args) |
