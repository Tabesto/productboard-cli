package cli

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/tabesto/productboard-cli/internal/client"
	"github.com/tabesto/productboard-cli/internal/output"
)

func invalidInput(format string, a ...any) error {
	return &client.APIError{Message: fmt.Sprintf(format, a...), ExitCode: client.ExitInvalidInput}
}

// parseFieldPairs turns ["name=x", "score=3"] into a fields map.
// Each value is JSON-decoded when valid (numbers, booleans, objects, arrays), else kept as a string.
func parseFieldPairs(pairs []string) (map[string]any, error) {
	fields := make(map[string]any, len(pairs))
	for _, p := range pairs {
		k, v, ok := strings.Cut(p, "=")
		if !ok || k == "" {
			return nil, invalidInput("invalid --field %q: expected key=value", p)
		}
		var val any
		if err := json.Unmarshal([]byte(v), &val); err != nil {
			val = v
		}
		fields[k] = val
	}
	return fields, nil
}

// parseRawData decodes --json as a JSON object (the request's "data" object).
func parseRawData(raw string) (map[string]any, error) {
	var data map[string]any
	if err := json.Unmarshal([]byte(raw), &data); err != nil || data == nil {
		return nil, invalidInput("invalid --json: expected a JSON object")
	}
	return data, nil
}

// buildUpdateBody builds {"data": ...} for an update from exactly one of --field, --patch, --json.
func buildUpdateBody(fieldPairs []string, patch, rawJSON string) (map[string]any, error) {
	n := 0
	for _, set := range []bool{len(fieldPairs) > 0, patch != "", rawJSON != ""} {
		if set {
			n++
		}
	}
	if n != 1 {
		return nil, invalidInput("provide exactly one of --field, --patch or --json")
	}

	switch {
	case rawJSON != "":
		data, err := parseRawData(rawJSON)
		if err != nil {
			return nil, err
		}
		return map[string]any{"data": data}, nil
	case patch != "":
		var ops []any
		if err := json.Unmarshal([]byte(patch), &ops); err != nil || ops == nil {
			return nil, invalidInput("invalid --patch: expected a JSON array")
		}
		return map[string]any{"data": map[string]any{"patch": ops}}, nil
	default:
		fields, err := parseFieldPairs(fieldPairs)
		if err != nil {
			return nil, err
		}
		return map[string]any{"data": map[string]any{"fields": fields}}, nil
	}
}

// buildEntityCreateBody builds {"data": {"type", "fields"}} for an entity create.
func buildEntityCreateBody(typ string, fieldPairs []string, rawJSON string) (map[string]any, error) {
	if rawJSON != "" {
		data, err := parseRawData(rawJSON)
		if err != nil {
			return nil, err
		}
		if typ != "" {
			data["type"] = typ
		}
		if data["type"] == nil {
			return nil, invalidInput("--type is required")
		}
		return map[string]any{"data": data}, nil
	}
	if typ == "" {
		return nil, invalidInput("--type is required")
	}
	fields, err := parseFieldPairs(fieldPairs)
	if err != nil {
		return nil, err
	}
	return map[string]any{"data": map[string]any{"type": typ, "fields": fields}}, nil
}

// buildNoteCreateBody builds {"data": {"type", "fields"}} for a note create.
// --name and --content are shortcuts for fields.name and fields.content.
func buildNoteCreateBody(typ, name, content string, fieldPairs []string, rawJSON string) (map[string]any, error) {
	if rawJSON != "" {
		data, err := parseRawData(rawJSON)
		if err != nil {
			return nil, err
		}
		if data["type"] == nil {
			data["type"] = typ
		}
		return map[string]any{"data": data}, nil
	}
	fields, err := parseFieldPairs(fieldPairs)
	if err != nil {
		return nil, err
	}
	if name != "" {
		fields["name"] = name
	}
	if content != "" {
		fields["content"] = content
	}
	if fields["name"] == nil {
		return nil, invalidInput("note name is required (--name or --field name=...)")
	}
	return map[string]any{"data": map[string]any{"type": typ, "fields": fields}}, nil
}

// printWriteResult prints the minimal reference returned by a create/update.
func printWriteResult(body []byte) error {
	ref, err := client.ParseWriteRef(body)
	if err != nil {
		return err
	}
	if outputFormat == output.FormatJSON {
		return output.PrintJSON(ref)
	}
	output.PrintSingleTable([][]string{{"ID", ref.ID}, {"Type", ref.Type}, {"Link", ref.Link}})
	return nil
}
