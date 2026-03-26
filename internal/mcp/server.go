package mcp

import (
	"github.com/mark3labs/mcp-go/server"
)

const defaultLimit = 25

// Version can be set by the CLI layer before calling Serve.
var Version = "dev"

// Serve creates the MCP server, registers all tools, and starts stdio transport.
func Serve() error {
	s := server.NewMCPServer(
		"pboard",
		Version,
		server.WithToolCapabilities(true),
	)

	registerTools(s)

	return server.ServeStdio(s)
}
