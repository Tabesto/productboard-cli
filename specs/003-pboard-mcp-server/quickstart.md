# Quickstart: MCP Server for pboard CLI

**Date**: 2026-03-26
**Branch**: `003-pboard-mcp-server`

## Prerequisites

- Go 1.22+
- pboard CLI installed and on PATH
- ProductBoard API token configured (`pboard configure <token>`)
- Claude Desktop installed on macOS

## Development Setup

```bash
# 1. Switch to feature branch
git checkout 003-pboard-mcp-server

# 2. Install new dependency
go get github.com/mark3labs/mcp-go
go mod tidy

# 3. Build
go build -o pboard ./cmd/pboard

# 4. Test MCP server locally (stdio mode -- type JSON-RPC to stdin)
./pboard mcp serve

# 5. Install into Claude Desktop
./pboard mcp install

# 6. Restart Claude Desktop and verify tools appear
```

## Key Files to Create/Modify

| File | Action | Purpose |
|------|--------|---------|
| `internal/mcp/server.go` | Create | MCP server setup, tool registration |
| `internal/mcp/tools.go` | Create | Tool definitions (all ~37 tools) |
| `internal/mcp/handlers.go` | Create | Tool handler functions |
| `internal/cli/mcp.go` | Create | `pboard mcp` command group (install/uninstall/serve) |
| `internal/cli/root.go` | Modify | Register `mcp` command |
| `go.mod` | Modify | Add `mcp-go` dependency |

## Testing the MCP Server

```bash
# Manual stdio test (sends initialize request)
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}}' | ./pboard mcp serve

# Verify install wrote correct config
cat ~/Library/Application\ Support/Claude/claude_desktop_config.json
```

## Architecture

```
pboard mcp serve  (stdio)
  └─> internal/mcp/server.go    -- Creates MCP server, registers tools
       └─> internal/mcp/tools.go     -- Tool definitions with schemas
       └─> internal/mcp/handlers.go  -- Handlers calling internal/client
            └─> internal/client/     -- Existing ProductBoard API client
            └─> internal/config/     -- Existing config (token) loading

pboard mcp install
  └─> internal/cli/mcp.go       -- Writes to Claude Desktop config JSON
```
