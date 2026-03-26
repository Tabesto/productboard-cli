# CLI Command Contract: pboard configure (updated)

**Date**: 2026-03-26

## `pboard configure [token]`

Configure API token for ProductBoard. Writes to `~/.config/pboard/config.yaml` with mode 600.

### Usage

```
# Non-interactive (new): provide token as argument
pboard configure pb_your_token_here

# Interactive (existing): prompted for token
pboard configure
```

### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `token`  | No       | ProductBoard API token. If omitted, an interactive prompt is shown. |

### Behavior

| Scenario | Behavior |
|----------|----------|
| Token argument provided | Save token directly, display confirmation, no prompt |
| No argument | Prompt user interactively (existing behavior) |
| Empty/whitespace token | Display error: "token cannot be empty" |
| Multiple arguments | Rejected by Cobra with usage message |

### Output

On success:
```
Token saved to /Users/<user>/.config/pboard/config.yaml
You can also set PRODUCTBOARD_API_TOKEN environment variable (takes precedence).
```

On error (empty token):
```
Error: token cannot be empty
```
