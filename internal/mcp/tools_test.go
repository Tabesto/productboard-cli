package mcp

import (
	"sort"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/server"
)

func writeToolNames(s *server.MCPServer) []string {
	var names []string
	for name := range s.ListTools() {
		if strings.HasPrefix(name, "create_") || strings.HasPrefix(name, "update_") {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	return names
}

func TestWriteToolsOnlyWhenEnabled(t *testing.T) {
	off := server.NewMCPServer("t", "0")
	registerTools(off, false)
	if got := writeToolNames(off); len(got) != 0 {
		t.Errorf("gate off: write tools registered: %v", got)
	}

	on := server.NewMCPServer("t", "0")
	registerTools(on, true)
	want := []string{"create_entity", "create_note", "update_entity", "update_note"}
	got := writeToolNames(on)
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Errorf("gate on: got %v, want %v", got, want)
	}
}
