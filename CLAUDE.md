# productboard-cli Development Guidelines

Auto-generated from all feature plans. Last updated: 2026-03-26

## Active Technologies
- Go 1.22+ (existing project) + Cobra (existing), no new dependencies (002-configure-token-arg)
- Existing config file at `~/.config/pboard/config.yaml` (mode 600) (002-configure-token-arg)
- Go 1.22+ (existing project, go.mod shows 1.26.1) + Cobra (CLI), mcp-go (MCP server), viper (config), go-pretty (output) (003-pboard-mcp-server)
- N/A (stateless server; reads config from `~/.config/pboard/config.yaml`) (003-pboard-mcp-server)
- Go 1.26.1 + Cobra (CLI), go-pretty (tables), mcp-go (MCP server), Viper (config) -- all existing, no new dependencies (004-feature-health-helper)
- N/A (reads from ProductBoard API via existing client) (004-feature-health-helper)

- Go 1.22+ + Cobra (CLI framework), tablewriter (table output), viper (config management) (001-productboard-cli-read)

## Project Structure

```text
src/
tests/
```

## Commands

# Add commands for Go 1.22+

## Code Style

Go 1.22+: Follow standard conventions

## Recent Changes
- 004-feature-health-helper: Added Go 1.26.1 + Cobra (CLI), go-pretty (tables), mcp-go (MCP server), Viper (config) -- all existing, no new dependencies
- 003-pboard-mcp-server: Added Go 1.22+ (existing project, go.mod shows 1.26.1) + Cobra (CLI), mcp-go (MCP server), viper (config), go-pretty (output)
- 002-configure-token-arg: Added Go 1.22+ (existing project) + Cobra (existing), no new dependencies


<!-- MANUAL ADDITIONS START -->
<!-- MANUAL ADDITIONS END -->
