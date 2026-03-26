# Data Model: MCP Server for pboard CLI

**Date**: 2026-03-26
**Branch**: `003-pboard-mcp-server`

## Entities

### MCP Tool Registration

Each tool maps 1:1 to a pboard CLI subcommand. No persistent data model -- tools are registered in memory at server startup.

| Tool Name | API Method | Parameters |
|-----------|-----------|------------|
| `list_features` | `GetList("/features", ...)` | status-name, status-id, parent-id, archived, owner-email, note-id, limit |
| `get_feature` | `GetSingle("/features/{id}")` | id (required) |
| `list_feature_initiatives` | `GetLinkedResources("/features/{id}/links/initiatives")` | id (required), limit |
| `list_feature_objectives` | `GetLinkedResources("/features/{id}/links/objectives")` | id (required), limit |
| `list_notes` | `GetList("/notes", ...)` | term, created-from, created-to, updated-from, updated-to, feature-id, company-id, owner-email, source, tags, limit |
| `get_note` | `GetSingle("/notes/{id}")` | id (required) |
| `list_note_tags` | `GetList("/notes/{id}/tags", ...)` | id (required), limit |
| `list_note_links` | `GetList("/notes/{id}/links", ...)` | id (required), limit |
| `list_products` | `GetList("/products", ...)` | limit |
| `get_product` | `GetSingle("/products/{id}")` | id (required) |
| `list_components` | `GetList("/components", ...)` | limit |
| `list_releases` | `GetList("/releases", ...)` | release-group-id, limit |
| `get_release` | `GetSingle("/releases/{id}")` | id (required) |
| `list_release_groups` | `GetList("/release-groups", ...)` | limit |
| `list_feature_release_assignments` | `GetList("/feature-release-assignments", ...)` | release-id, limit |
| `list_objectives` | `GetList("/objectives", ...)` | limit |
| `get_objective` | `GetSingle("/objectives/{id}")` | id (required) |
| `list_objective_features` | `GetLinkedResources("/objectives/{id}/links/features")` | id (required), limit |
| `list_key_results` | `GetList("/key-results", ...)` | limit |
| `list_initiatives` | `GetList("/initiatives", ...)` | limit |
| `get_initiative` | `GetSingle("/initiatives/{id}")` | id (required) |
| `list_initiative_objectives` | `GetLinkedResources("/initiatives/{id}/links/objectives")` | id (required), limit |
| `list_initiative_features` | `GetLinkedResources("/initiatives/{id}/links/features")` | id (required), limit |
| `list_companies` | `GetList("/companies", ...)` | term, has-notes, feature-id, limit |
| `get_company` | `GetSingle("/companies/{id}")` | id (required) |
| `list_company_custom_fields` | `GetList("/companies/custom-fields", ...)` | limit |
| `get_company_custom_field` | `GetSingle("/companies/custom-fields/{id}")` | id (required) |
| `get_company_custom_field_value` | `GetSingle("/companies/{companyId}/custom-fields/{fieldId}")` | company-id (required), field-id (required) |
| `list_users` | `GetList("/users", ...)` | limit |
| `list_custom_fields` | `GetList("/hierarchy-entities/custom-fields", ...)` | type (required), limit |
| `get_custom_field` | `GetSingle("/hierarchy-entities/custom-fields/{id}")` | id (required) |
| `list_custom_field_values` | `GetList("/hierarchy-entities/custom-fields/{id}/values", ...)` | custom-field-id (required), limit |
| `list_feature_statuses` | `GetList("/feature-statuses", ...)` | limit |
| `list_plugin_integrations` | `GetList("/plugin-integrations", ...)` | limit |
| `list_jira_integrations` | `GetList("/jira-integrations", ...)` | limit |
| `list_webhooks` | `GetList("/webhooks", ...)` | limit |
| `get_webhook` | `GetSingle("/webhooks/{id}")` | id (required) |
| `list_feedback_forms` | `GetList("/feedback-forms", ...)` | limit |

### Claude Desktop Configuration Entry

Written to `~/Library/Application Support/Claude/claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "pboard": {
      "command": "/path/to/pboard",
      "args": ["mcp", "serve"]
    }
  }
}
```

- `command`: Absolute path to the pboard binary (resolved at install time)
- `args`: Fixed arguments to start MCP server mode

### Default Behavior

- All `list_*` tools default to `limit=25` when no limit parameter is provided
- All tools return JSON-formatted results as MCP text content
- Error responses use MCP error result format with descriptive messages
