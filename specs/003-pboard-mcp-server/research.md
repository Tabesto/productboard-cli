# Research: MCP Server for pboard CLI

**Date**: 2026-03-26
**Branch**: `003-pboard-mcp-server`

## R1: Go MCP Server Library

**Decision**: Use `github.com/mark3labs/mcp-go` (community library)

**Rationale**: Most popular Go MCP library (~8.4k stars), mature API with fluent builder pattern for tool definitions, straightforward stdio transport support via `server.ServeStdio()`. Simpler ergonomics than the official `modelcontextprotocol/go-sdk`.

**Alternatives considered**:
- `modelcontextprotocol/go-sdk` (official SDK, ~4.2k stars) -- newer, more verbose API, fewer community examples. Viable but less ergonomic for a tool-heavy server.
- Custom JSON-RPC implementation -- unnecessary complexity when libraries exist.

## R2: MCP Transport Protocol

**Decision**: Stdio transport (stdin/stdout JSON-RPC)

**Rationale**: Standard transport for local MCP servers launched by Claude Desktop. The client spawns the server binary as a subprocess, writes JSON-RPC to stdin, reads responses from stdout. Stderr is available for logging. No networking needed.

**Alternatives considered**:
- HTTP/SSE transport -- requires a running server process; overkill for a local tool.

## R3: Tool Definition Pattern

**Decision**: One MCP tool per CLI subcommand, using the `mcp-go` fluent builder

**Rationale**: Each tool maps 1:1 to a pboard CLI subcommand (e.g., `list_features`, `get_feature`). Tools are defined with `mcp.NewTool()` + `mcp.WithString()`/`mcp.WithNumber()` for parameters. Handler functions call the existing `internal/client` API directly (bypassing CLI command layer).

**Alternatives considered**:
- Shelling out to `pboard` CLI binary -- adds process overhead, parsing complexity; using the client library directly is cleaner.
- Consolidated tools with action parameters -- rejected per clarification (one tool per subcommand preferred).

## R4: Server Entry Point Architecture

**Decision**: New `pboard mcp serve` subcommand in the existing binary (no separate binary)

**Rationale**: The MCP server runs as `pboard mcp serve` via stdio. This means the install command registers `pboard mcp serve` as the command in Claude Desktop config. Single binary keeps distribution simple -- no separate build target needed.

**Alternatives considered**:
- Separate `pboard-mcp` binary -- extra build/install complexity for no benefit.
- Embedding server in the root command -- would conflict with CLI usage.

## R5: Claude Desktop Configuration

**Decision**: Target `~/Library/Application Support/Claude/claude_desktop_config.json` on macOS

**Rationale**: Standard Claude Desktop config location on macOS. The install command reads existing JSON (or creates new), adds/updates the `mcpServers.pboard` entry with `{"command": "<pboard-path>", "args": ["mcp", "serve"]}`, and writes back.

**Alternatives considered**:
- Supporting multiple OS config paths -- deferred; macOS-only per current user base. Can extend later.

## R6: Install Command Pattern

**Decision**: Follow the existing `pboard skill install` pattern

**Rationale**: The codebase already has `skill.go` with `--force`, `--dry-run`, directory creation, and idempotent install. The `mcp install` command follows the same structure: detect config path, check existing entry, write/update JSON, report success.

**Alternatives considered**:
- Interactive wizard -- overengineered for a single config write.

## R7: Reusing Internal Client

**Decision**: MCP tool handlers call `internal/client` directly, returning JSON data

**Rationale**: The existing `client.GetList()`, `client.GetSingle()`, `client.GetLinkedResources()` methods already handle auth, pagination, and error mapping. Tool handlers load config, create client, call the appropriate method, and return JSON as tool result text. Default limit of 25 applied when no limit parameter provided.

**Alternatives considered**:
- Duplicating API logic in MCP handlers -- violates DRY.
