package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/tabesto/productboard-cli/internal/client"
)

const writeHint = " Read the workspace configuration first (field names and types are workspace-specific) and confirm intent with the user before writing."

func registerWriteTools(s *server.MCPServer) {
	s.AddTool(mcp.NewTool("create_entity",
		mcp.WithDescription("Create a ProductBoard entity (feature, initiative, objective, ...). WRITES DATA."+writeHint),
		mcp.WithString("type", mcp.Required(), mcp.Description("Entity type, e.g. feature or initiative")),
		mcp.WithObject("fields", mcp.Description("Field values keyed by field name, e.g. {\"name\": \"My feature\"}")),
	), handleCreateEntity)

	s.AddTool(mcp.NewTool("update_entity",
		mcp.WithDescription("Update a ProductBoard entity. Provide either fields (replace values) or patch (set/addItems/removeItems/clear operations), not both. WRITES DATA."+writeHint),
		mcp.WithString("id", mcp.Required(), mcp.Description("Entity ID")),
		mcp.WithObject("fields", mcp.Description("Field values to replace, keyed by field name")),
		mcp.WithArray("patch", mcp.Description("Patch operations (set, addItems, removeItems, clear)")),
	), handleUpdateEntity)

	s.AddTool(mcp.NewTool("create_note",
		mcp.WithDescription("Create a ProductBoard note. WRITES DATA."+writeHint),
		mcp.WithString("name", mcp.Required(), mcp.Description("Note name")),
		mcp.WithString("type", mcp.Description("Note type: textNote (default) or conversationNote")),
		mcp.WithString("content", mcp.Description("Note content")),
		mcp.WithObject("fields", mcp.Description("Additional field values keyed by field name")),
	), handleCreateNote)

	s.AddTool(mcp.NewTool("update_note",
		mcp.WithDescription("Update a ProductBoard note. Provide either fields or patch, not both. WRITES DATA."+writeHint),
		mcp.WithString("id", mcp.Required(), mcp.Description("Note ID")),
		mcp.WithObject("fields", mcp.Description("Field values to replace, keyed by field name")),
		mcp.WithArray("patch", mcp.Description("Patch operations (set, addItems, removeItems, clear)")),
	), handleUpdateNote)
}

func handleCreateEntity(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	typ, err := request.RequireString("type")
	if err != nil {
		return errorResult(err), nil
	}
	data := map[string]any{"type": typ, "fields": optionalObject(request, "fields")}
	return doWrite("POST", "/entities", data)
}

func handleUpdateEntity(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return updateResource(request, "/entities/")
}

func handleCreateNote(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := request.RequireString("name")
	if err != nil {
		return errorResult(err), nil
	}
	fields := optionalObject(request, "fields")
	fields["name"] = name
	if content := request.GetString("content", ""); content != "" {
		fields["content"] = content
	}
	data := map[string]any{"type": request.GetString("type", "textNote"), "fields": fields}
	return doWrite("POST", "/notes", data)
}

func handleUpdateNote(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return updateResource(request, "/notes/")
}

func updateResource(request mcp.CallToolRequest, prefix string) (*mcp.CallToolResult, error) {
	id, err := request.RequireString("id")
	if err != nil {
		return errorResult(err), nil
	}
	args := request.GetArguments()
	fields, hasFields := args["fields"].(map[string]any)
	patch, hasPatch := args["patch"].([]any)
	if hasFields == hasPatch {
		return mcp.NewToolResultError("provide exactly one of fields or patch"), nil
	}
	data := map[string]any{}
	if hasFields {
		data["fields"] = fields
	} else {
		data["patch"] = patch
	}
	return doWrite("PATCH", prefix+id, data)
}

func optionalObject(request mcp.CallToolRequest, key string) map[string]any {
	if m, ok := request.GetArguments()[key].(map[string]any); ok {
		return m
	}
	return map[string]any{}
}

// doWrite sends the request through the client, which enforces the write gate.
func doWrite(method, path string, data map[string]any) (*mcp.CallToolResult, error) {
	c, err := getClient()
	if err != nil {
		return errorResult(err), nil
	}
	if err := c.RequireWrite(); err != nil {
		return errorResult(err), nil
	}
	body := map[string]any{"data": data}
	var resp []byte
	if method == "POST" {
		resp, err = c.Post(path, body)
	} else {
		resp, err = c.Patch(path, body)
	}
	if err != nil {
		return errorResult(err), nil
	}
	ref, err := client.ParseWriteRef(resp)
	if err != nil {
		return errorResult(err), nil
	}
	s, err := toJSON(ref)
	if err != nil {
		return errorResult(err), nil
	}
	return mcp.NewToolResultText(s), nil
}
