package mcp

import (
	"github.com/mark3labs/mcp-go/server"
	"github.com/tabesto/productboard-cli/internal/config"
)

const defaultLimit = 25

// Version can be set by the CLI layer before calling Serve.
var Version = "dev"

// writeEnabled is the effective write gate, fixed when Serve starts.
var writeEnabled bool

// Serve creates the MCP server, registers all tools, and starts stdio transport.
// Write tools are registered only when allowWrite or the config/env gate is on.
func Serve(allowWrite bool) error {
	if cfg, err := config.Load(); err == nil {
		allowWrite = allowWrite || cfg.WriteEnabled
	}
	writeEnabled = allowWrite

	s := server.NewMCPServer(
		"pboard",
		Version,
		server.WithToolCapabilities(true),
	)

	registerTools(s, allowWrite)

	return server.ServeStdio(s)
}
