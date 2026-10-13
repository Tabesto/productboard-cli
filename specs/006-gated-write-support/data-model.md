# Data Model: Gated Write Support

## Config (extended) — `internal/config`
| Field | Source (low → high precedence) | Notes |
|-------|-------------------------------|-------|
| `APIToken` | config `api_token` → env `PRODUCTBOARD_API_TOKEN` | unchanged |
| `BaseURL` | config `api_url` | unchanged |
| `APIVersion` | config `api_version` → env → `--api-version` | unchanged |
| **`WriteEnabled`** (bool, default false) | config `write_enabled` → env `PRODUCTBOARD_WRITE_ENABLED` (`1`/`true` on, other non-empty off) → `--allow-write` (CLI) / `mcp serve --allow-write` (can only turn on) | new |

Config file after `configure --write-enabled`:
```yaml
api_token: "…"
api_url: "https://api.productboard.com"
api_version: "2"
write_enabled: true
```

## Write gate (derived)
`effective = cfg.WriteEnabled || flagAllowWrite`. Stored in `Client.allowWrite`. State: off (default) / on. No transitions at runtime; fixed for the process lifetime (MCP: fixed at server start).

## WriteRef (response) — `internal/client/write.go`
| Field | Type | From |
|-------|------|------|
| `ID` | string | `data.id` |
| `Type` | string | `data.type` |
| `Link` | string | `data.links.self` |

## Request bodies (JSON sent to API)
- Create entity: `{"data":{"type":T,"fields":F}}`
- Create note: `{"data":{"type":T,"fields":F}}` with `F.name` required
- Update (entity/note): `{"data":{"fields":F}}` **xor** `{"data":{"patch":P}}`
- `--json` replaces the `data` object wholesale (user-supplied).

## Validation rules
- Gate off → `APIError{ExitCode: 6}` before any I/O.
- API version ≠ 2 → `APIError{ExitCode: 5}`.
- `--field` + `--patch` together → local error (exit 5).
- `--field` not of form `key=value` → local error (exit 5).
- `--patch` / `--json` not valid JSON (array / object respectively) → local error (exit 5).
- Note create without `name` (via `--name`, `--field name=…` or `--json`) → local error (exit 5).
- Everything else (unknown type/field, forbidden note type) is validated by Productboard; its `errors[0].detail` is surfaced via existing `NewAPIError`.
