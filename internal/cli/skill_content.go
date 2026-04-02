package cli

const skillContent = `---
description: Guide for using the pboard CLI to query ProductBoard data. Use when an agent needs to browse features, notes, releases, companies, or any other ProductBoard entity.
---

# ProductBoard CLI (pboard) — Agent Usage Guide

You have access to ` + "`pboard`" + `, a **read-only** CLI for browsing ProductBoard data. It cannot create, update, or delete anything.

## Setup

The CLI is already configured. If you get an auth error, the user needs to run:
` + "```bash" + `
pboard configure <token>
` + "```" + `

## Global Flags

| Flag | Description |
|------|-------------|
| ` + "`-o json`" + ` | Output as JSON (default: table). **Always use ` + "`-o json`" + ` when you need to parse output programmatically.** |
| ` + "`-l N`" + ` | Limit results to N items (default: 0 = all) |
| ` + "`--api-version`" + ` | API version: 1 or 2 (default: 2). Use ` + "`--api-version 1`" + ` for V1-only commands. |

## Command Reference

### Features
` + "```bash" + `
pboard features list [flags]             # List features
pboard features get <id>                 # Get single feature (shows Description)
pboard features links initiatives <id>   # Linked initiatives
pboard features links objectives <id>    # Linked objectives
` + "```" + `
**List flags:** ` + "`--status-id`" + `, ` + "`--status-name`" + `, ` + "`--parent-id`" + `, ` + "`--archived`" + ` (true/false), ` + "`--owner-email`" + `, ` + "`--note-id`" + `
**List columns:** ID, Name, Status, Owner, Archived
**Get fields:** ID, Name, Status, Owner, Archived, Description

### Feature Health
` + "```bash" + `
pboard features health list [flags]      # List features with health updates
pboard features health get <id>          # Get health details for a feature
` + "```" + `
**List flags:** ` + "`--updated-since`" + ` (YYYY-MM-DD), ` + "`--updated-before`" + `, ` + "`--status`" + `, ` + "`--owner`" + `, ` + "`--health-status`" + ` (on-track/at-risk/off-track), ` + "`--include-archived`" + `, ` + "`--include-no-health`" + `
**List columns:** Feature Name, Status, Owner, Health, Health Updated, Message

### Products
` + "```bash" + `
pboard products list                     # List products
pboard products get <id>                 # Get single product
` + "```" + `
**List columns:** ID, Name, Description

### Components
` + "```bash" + `
pboard components list                   # List components
pboard components get <id>               # Get single component
` + "```" + `
**List columns:** ID, Name, Description, Parent

### Notes (Feedback)
` + "```bash" + `
pboard notes list [flags]                # List notes
pboard notes get <id>                    # Get single note (shows Content)
pboard notes tags <noteId>               # List tags for a note
pboard notes links <noteId>              # List links for a note
` + "```" + `
**List flags:** ` + "`--date-from`" + `, ` + "`--date-to`" + `, ` + "`--created-from`" + `, ` + "`--created-to`" + `, ` + "`--updated-from`" + `, ` + "`--updated-to`" + `, ` + "`--term`" + `, ` + "`--feature-id`" + `, ` + "`--company-id`" + `, ` + "`--owner-email`" + `, ` + "`--source`" + `, ` + "`--any-tag`" + ` (comma-separated), ` + "`--all-tags`" + ` (comma-separated)
**List columns:** ID, Title, Source, Owner, Created
**Get fields:** ID, Title, Source, Owner, Created At, Updated At, Content

> Dates use ISO 8601 format: ` + "`2026-03-19`" + ` or ` + "`2026-03-19T00:00:00Z`" + `

### Companies
` + "```bash" + `
pboard companies list [flags]            # List companies
pboard companies get <id>                # Get single company
pboard companies custom-fields list      # List company custom fields
pboard companies custom-fields get <id>
pboard companies custom-field-value <companyId> <fieldId>
` + "```" + `
**List flags:** ` + "`--term`" + `, ` + "`--has-notes`" + ` (true/false), ` + "`--feature-id`" + `
**List columns:** ID, Name, Domain

### Users
` + "```bash" + `
pboard users list                        # List users
pboard users get <id>                    # Get single user
` + "```" + `
**List columns:** ID, Email, Name

### Releases
` + "```bash" + `
pboard releases list [flags]             # List releases
pboard releases get <id>                 # Get single release
` + "```" + `
**List flags:** ` + "`--release-group-id`" + `
**List columns:** ID, Name, State, Start Date, End Date

### Release Groups
` + "```bash" + `
pboard release-groups list               # List release groups
pboard release-groups get <id>           # Get single release group
` + "```" + `
**List columns:** ID, Name, Description

### Feature-Release Assignments
` + "```bash" + `
pboard feature-release-assignments list [flags]
pboard feature-release-assignments get --feature-id <fid> --release-id <rid>
` + "```" + `
**List flags:** ` + "`--feature-id`" + `, ` + "`--release-id`" + `, ` + "`--release-state`" + `, ` + "`--end-date-from`" + `, ` + "`--end-date-to`" + `
**List columns:** Feature ID, Release ID

### Objectives
` + "```bash" + `
pboard objectives list                   # List objectives
pboard objectives get <id>               # Get single objective
pboard objectives links features <id>    # Linked features
pboard objectives links initiatives <id> # Linked initiatives
` + "```" + `
**List columns:** ID, Name, State, Owner

### Key Results
` + "```bash" + `
pboard key-results list                  # List key results
pboard key-results get <id>              # Get single key result
` + "```" + `
**List columns:** ID, Name, Current Value, Target Value, Objective ID

### Initiatives
` + "```bash" + `
pboard initiatives list                  # List initiatives
pboard initiatives get <id>              # Get single initiative
pboard initiatives links objectives <id> # Linked objectives
pboard initiatives links features <id>   # Linked features
` + "```" + `
**List columns:** ID, Name, State, Owner

### Feature Statuses
` + "```bash" + `
pboard feature-statuses list             # List all feature statuses
` + "```" + `
**List columns:** ID, Name

### Custom Fields (Hierarchy Entities)
` + "```bash" + `
pboard custom-fields list --type <type>  # List custom fields (--type required)
pboard custom-fields get <id>
pboard custom-fields values list [flags]
pboard custom-fields values get --custom-field-id <cfid> --hierarchy-entity-id <heid>
` + "```" + `
**Values list flags:** ` + "`--type`" + `, ` + "`--custom-field-id`" + `, ` + "`--hierarchy-entity-id`" + `
**List columns:** ID, Name, Type

### Members (V2 only)
` + "```bash" + `
pboard members list [flags]              # List workspace members
pboard members get <id>                  # Get single member
` + "```" + `
**List flags:** ` + "`--role`" + ` (admin/maker/viewer/contributor), ` + "`--query`" + `
**List columns:** ID, Name, Email, Role

### Teams (V2 only)
` + "```bash" + `
pboard teams list [flags]                # List teams
pboard teams get <id>                    # Get single team
` + "```" + `
**List flags:** ` + "`--query`" + `
**List columns:** ID, Name, Handle, Description

### Feedback Forms (V1 only)
` + "```bash" + `
pboard feedback-forms list --api-version 1    # List feedback forms
pboard feedback-forms get <id> --api-version 1
` + "```" + `
**List columns:** ID, Name
> **Note:** Feedback forms are not available with API V2. Use ` + "`--api-version 1`" + `.

### Webhooks
` + "```bash" + `
pboard webhooks list                     # List webhooks
pboard webhooks get <id>                 # Get single webhook
` + "```" + `
**List columns:** ID, URL, Events

### Jira Integrations
` + "```bash" + `
pboard jira-integrations list
pboard jira-integrations get <id>
pboard jira-integrations connections list <id> [flags]
pboard jira-integrations connections get <id> <featureId>
` + "```" + `
**Connections list flags:** ` + "`--issue-key`" + `, ` + "`--issue-id`" + `

### Plugin Integrations
` + "```bash" + `
pboard plugin-integrations list
pboard plugin-integrations get <id>
pboard plugin-integrations connections list <id>
pboard plugin-integrations connections get <id> <featureId>
` + "```" + `

## Agent Best Practices

1. **Always use ` + "`-o json`" + `** when you need to process the output programmatically or extract specific fields.
2. **Use ` + "`-l N`" + `** to limit results when you only need a sample or the first few items.
3. **Use ` + "`get`" + ` for full details** — ` + "`list`" + ` shows summary columns only. For example, ` + "`notes list`" + ` doesn't show Content, but ` + "`notes get <id>`" + ` does.
4. **Chain commands** to explore relationships:
   - Find features in a release: ` + "`pboard feature-release-assignments list --release-id <id> -o json`" + `
   - Then get each feature: ` + "`pboard features get <featureId> -o json`" + `
5. **Filter notes by date** to scope feedback queries: ` + "`--created-from 2026-03-01 --created-to 2026-03-31`" + `
6. **Traverse links** to understand relationships between entities (features <-> initiatives <-> objectives).

## Common Workflows

**Get recent feedback with content:**
` + "```bash" + `
# List recent notes
pboard notes list --created-from 2026-03-19 -o json
# Get full content of a specific note
pboard notes get <noteId> -o json
` + "```" + `

**Explore a release scope:**
` + "```bash" + `
# List releases
pboard releases list -o json
# Find features assigned to a release
pboard feature-release-assignments list --release-id <releaseId> -o json
# Get details of each feature
pboard features get <featureId> -o json
` + "```" + `

**Explore strategic alignment:**
` + "```bash" + `
# List objectives
pboard objectives list -o json
# See which features support an objective
pboard objectives links features <objectiveId> -o json
# See which initiatives support an objective
pboard objectives links initiatives <objectiveId> -o json
` + "```" + `
`
