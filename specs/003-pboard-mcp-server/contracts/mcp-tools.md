# MCP Tool Contracts

**Date**: 2026-03-26
**Protocol**: MCP (Model Context Protocol) via stdio transport

## Server Identity

- **Name**: `pboard`
- **Version**: matches pboard CLI version
- **Transport**: stdio (JSON-RPC over stdin/stdout)

## Tool Contract Pattern

All tools follow a consistent contract:

### List Tools

```
Tool: list_{resource}
Parameters:
  - limit (number, optional, default: 25): Maximum results to return
  - [resource-specific filters] (string, optional): Filter parameters
Returns: JSON array of resource objects
Errors: Authentication error, rate limit, API error
```

### Get Tools

```
Tool: get_{resource}
Parameters:
  - id (string, required): Resource identifier
Returns: JSON object with resource details
Errors: Not found, authentication error, API error
```

### Link/Relationship Tools

```
Tool: list_{parent}_{children}
Parameters:
  - id (string, required): Parent resource identifier
  - limit (number, optional, default: 25): Maximum results
Returns: JSON array of linked resource objects
Errors: Not found, authentication error, API error
```

## Error Response Contract

All errors are returned as MCP tool error results with descriptive messages:

| Error Condition | Message Pattern |
|----------------|-----------------|
| Missing token | "ProductBoard API token not configured. Run 'pboard configure <token>' to set it up." |
| Invalid/expired token | "Authentication failed: invalid or expired API token." |
| Resource not found | "{Resource} with ID '{id}' not found." |
| Rate limited | "ProductBoard API rate limit reached. Please try again later." |
| API error | "ProductBoard API error: {status} - {message}" |

## Claude Desktop Configuration Contract

`~/Library/Application Support/Claude/claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "pboard": {
      "command": "{absolute-path-to-pboard-binary}",
      "args": ["mcp", "serve"]
    }
  }
}
```

- Install: Merges `pboard` entry into existing `mcpServers` object (preserves other servers)
- Uninstall: Removes only `pboard` key from `mcpServers` (preserves other servers)
- If file doesn't exist: Creates with `{"mcpServers": {"pboard": {...}}}`
