# CLI Command Contracts: Feature Health

## `pboard features health`

Parent subcommand group. No direct action; displays help.

```
Usage: pboard features health [command]

Commands:
  list    List health updates across features
  get     Get health update for a specific feature
```

---

## `pboard features health list`

Fetch all features and display a consolidated health overview with optional filters.

```
Usage: pboard features health list [flags]

Flags:
  --updated-since    string   Show features with health updated on or after this date (YYYY-MM-DD)
  --updated-before   string   Show features with health updated before this date (YYYY-MM-DD)
  --status           string   Filter by feature status name (e.g., "In Progress")
  --owner            string   Filter by feature owner email
  --health-status    string   Filter by health status (on-track, at-risk, off-track)
  --include-archived          Include archived features (excluded by default)
  --include-no-health         Include features without health updates

Global Flags:
  -o, --output string   Output format: table or json (default "table")
  -l, --limit  int      Maximum number of results (0 = all, default 0)
```

### Table Output

```
+-------------------+--------------+---------------------------+----------+---------------------+-----------------------------+
| FEATURE NAME      | STATUS       | OWNER                     | HEALTH   | HEALTH UPDATED      | MESSAGE                     |
+-------------------+--------------+---------------------------+----------+---------------------+-----------------------------+
| PRM accessibility | Live         | amel.meradi@deliverect... | on-track | 2025-11-13T10:00:00 | Reverted Nepting from th... |
| Smart discounts   | Backlog      | maeva.bonin@deliverect... | at-risk  | 2025-10-01T08:00:00 | Blocked on pricing API...   |
+-------------------+--------------+---------------------------+----------+---------------------+-----------------------------+
```

Sorted by health update date (most recent first). Message truncated to 50 characters in table view.

### JSON Output

```json
[
  {
    "id": "7043709b-...",
    "name": "PRM accessibility standard",
    "status": { "name": "Live" },
    "owner": { "email": "amel.meradi@deliverect.com" },
    "archived": false,
    "lastHealthUpdate": {
      "status": "on-track",
      "updatedAt": "2025-11-13T10:00:00Z",
      "message": "Reverted Nepting from the release"
    }
  }
]
```

### Exit Codes

| Code | Condition |
|------|-----------|
| 0 | Success (including empty results) |
| 1 | General error |
| 2 | Authentication error (invalid/missing token) |
| 4 | Rate limited |
| 5 | Invalid input (bad date format) |

---

## `pboard features health get <feature-id>`

Fetch a single feature and display its full health update details.

```
Usage: pboard features health get <feature-id>

Arguments:
  feature-id   UUID of the feature (required)

Global Flags:
  -o, --output string   Output format: table or json (default "table")
```

### Table Output (single resource view)

```
+------------------+-------------------------------------------+
| FIELD            | VALUE                                     |
+------------------+-------------------------------------------+
| ID               | 7043709b-f020-43f0-bebb-8ae7649acbf7      |
| Name             | PRM accessibility standard                |
| Status           | Live                                      |
| Owner            | amel.meradi@deliverect.com                |
| Health Status    | on-track                                  |
| Health Updated   | 2025-11-13T10:00:00Z                      |
| Message          | Reverted Nepting from the release         |
+------------------+-------------------------------------------+
```

When no health update exists:
```
+------------------+-------------------------------------------+
| FIELD            | VALUE                                     |
+------------------+-------------------------------------------+
| ID               | 7043709b-f020-43f0-bebb-8ae7649acbf7      |
| Name             | PRM accessibility standard                |
| Status           | Live                                      |
| Owner            | amel.meradi@deliverect.com                |
| Health Status    | (none)                                    |
+------------------+-------------------------------------------+
```

### Exit Codes

| Code | Condition |
|------|-----------|
| 0 | Success |
| 1 | General error |
| 2 | Authentication error |
| 3 | Feature not found |
| 4 | Rate limited |
