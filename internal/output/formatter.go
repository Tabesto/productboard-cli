package output

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
)

// Format constants.
const (
	FormatTable = "table"
	FormatJSON  = "json"
)

// PrintJSON outputs data as formatted JSON to stdout.
func PrintJSON(data interface{}) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}

// PrintTable outputs data as a formatted table to stdout.
func PrintTable(headers []string, rows [][]string) {
	if len(rows) == 0 {
		fmt.Println("No results found.")
		return
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	headerRow := table.Row{}
	for _, h := range headers {
		headerRow = append(headerRow, h)
	}
	t.AppendHeader(headerRow)

	for _, row := range rows {
		tableRow := table.Row{}
		for _, cell := range row {
			tableRow = append(tableRow, cell)
		}
		t.AppendRow(tableRow)
	}

	t.SetStyle(table.StyleLight)
	t.Render()
}

// PrintSingleTable outputs a single resource as a key-value table.
func PrintSingleTable(fields [][]string) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Field", "Value"})

	for _, field := range fields {
		t.AppendRow(table.Row{field[0], field[1]})
	}

	t.SetStyle(table.StyleLight)
	t.Render()
}

// Print outputs data in the specified format.
func Print(format string, data interface{}, headers []string, toRows func() [][]string) error {
	switch format {
	case FormatJSON:
		return PrintJSON(data)
	case FormatTable:
		PrintTable(headers, toRows())
		return nil
	default:
		return fmt.Errorf("unknown output format: %s", format)
	}
}

// Truncate shortens a string to maxLen, appending "..." if truncated.
func Truncate(s string, maxLen int) string {
	s = strings.ReplaceAll(s, "\n", " ")
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

// SafeStr extracts a string from a map, returning "" if not found.
func SafeStr(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok && v != nil {
		return fmt.Sprintf("%v", v)
	}
	return ""
}

// SafeNested extracts a nested string value from a map.
func SafeNested(m map[string]interface{}, keys ...string) string {
	current := m
	for i, key := range keys {
		v, ok := current[key]
		if !ok || v == nil {
			return ""
		}
		if i == len(keys)-1 {
			return fmt.Sprintf("%v", v)
		}
		nested, ok := v.(map[string]interface{})
		if !ok {
			return ""
		}
		current = nested
	}
	return ""
}
