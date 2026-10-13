# Interface Contracts: Gated Write Support

## Global

| Flag | Scope | Meaning |
|------|-------|---------|
| `--allow-write` | persistent (all commands) | Enable writes for this run. Never persisted. |

Exit codes: existing 0–5, plus **6 = writes disabled**.
Refusal message (stderr): `Error: write operations are disabled. Enable with --allow-write, \`pboard configure --write-enabled\`, or PRODUCTBOARD_WRITE_ENABLED=1.`

## CLI

### `pboard entities create`
```
pboard entities create --type <type> [--field key=value]... [--json '<data object>'] [-o table|json]
```
- `--type` required unless `--json` contains `type`.
- Sends `POST /v2/entities`. Prints `ID | Type | Link` (table) or `{id,type,links}` (json).

### `pboard entities update <id>`
```
pboard entities update <id> (--field key=value... | --patch '<json array>' | --json '<data object>')
```
- Exactly one of `--field`, `--patch`, `--json`. Sends `PATCH /v2/entities/<id>`.

### `pboard notes create`
```
pboard notes create --name <name> [--type textNote|conversationNote] [--content <text>] [--field key=value]... [--json '<data object>']
```
- `--type` default `textNote`. `--name`/`--content`/`--field` are shortcuts for `fields.*` and are merged into a `--json` object too (flags win). A non-empty `fields.name` is always required locally (exit 5). Sends `POST /v2/notes`.

### `pboard notes update <id>`
Same rule as entities update. Sends `PATCH /v2/notes/<id>`.

### `pboard configure`
```
pboard configure [token] [--write-enabled[=true|false]]
```
- With only `--write-enabled`: updates that key, no prompt. With token (arg or prompt): updates token. Other keys always preserved.

### `pboard mcp serve [--allow-write]` / `pboard mcp install [--allow-write] [--force] [--dry-run]`
- `install --allow-write` registers args `["mcp","serve","--allow-write"]`.

## MCP tools (registered only when gate is on)

| Tool | Required args | Optional args |
|------|---------------|---------------|
| `create_entity` | `type` (string) | `fields` (object) |
| `update_entity` | `id` (string) | `fields` (object) **xor** `patch` (array) |
| `create_note` | `name` (string) | `type` (string, default `textNote`), `content` (string), `fields` (object) |
| `update_note` | `id` (string) | `fields` (object) **xor** `patch` (array) |

Result text: JSON `{"id","type","links":{"self"}}`. Errors: `isError` result with the client message.
Descriptions MUST tell the agent to read `/entities/configurations` / `/notes/configurations` data (via existing list tools) before writing.

## HTTP (outbound)

| Call | Headers |
|------|---------|
| `POST {base}/v2/entities`, `PATCH {base}/v2/entities/{id}`, `POST {base}/v2/notes`, `PATCH {base}/v2/notes/{id}` | `Authorization: Bearer <token>`, `Accept: application/json`, `Content-Type: application/json`; **no** `X-Version` |
