package cli

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tabesto/productboard-cli/internal/output"
)

func newCustomFieldsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "custom-fields",
		Short: "Manage ProductBoard hierarchy-entity custom fields",
	}

	cmd.AddCommand(newCustomFieldsListCmd())
	cmd.AddCommand(newCustomFieldsGetCmd())
	cmd.AddCommand(newCustomFieldsValuesCmd())

	return cmd
}

func newCustomFieldsListCmd() *cobra.Command {
	var fieldTypes []string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List custom fields",
		Run: func(cmd *cobra.Command, args []string) {
			c, err := getClient()
			if err != nil {
				handleError(err)
			}

			params := map[string]string{
				"type": strings.Join(fieldTypes, ","),
			}

			items, err := c.GetList("/hierarchy-entities/custom-fields", params, limitFlag)
			if err != nil {
				handleError(err)
			}

			headers := []string{"ID", "Name", "Type"}
			if err := output.Print(outputFormat, items, headers, func() [][]string {
				rows := make([][]string, 0, len(items))
				for _, item := range items {
					rows = append(rows, []string{
						output.SafeStr(item, "id"),
						output.SafeStr(item, "name"),
						output.SafeStr(item, "type"),
					})
				}
				return rows
			}); err != nil {
				handleError(err)
			}
		},
	}

	cmd.Flags().StringSliceVar(&fieldTypes, "type", nil, "Filter by type(s)")
	if err := cmd.MarkFlagRequired("type"); err != nil {
		panic(err)
	}

	return cmd
}

func newCustomFieldsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a custom field by ID",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			c, err := getClient()
			if err != nil {
				handleError(err)
			}

			item, err := c.GetSingle("/hierarchy-entities/custom-fields/" + args[0])
			if err != nil {
				handleError(err)
			}

			if outputFormat == output.FormatJSON {
				if err := output.PrintJSON(item); err != nil {
					handleError(err)
				}
				return
			}

			output.PrintSingleTable([][]string{
				{"ID", output.SafeStr(item, "id")},
				{"Name", output.SafeStr(item, "name")},
				{"Type", output.SafeStr(item, "type")},
			})
		},
	}
}

func newCustomFieldsValuesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "values",
		Short: "Manage custom field values",
	}

	cmd.AddCommand(newCustomFieldsValuesListCmd())
	cmd.AddCommand(newCustomFieldsValuesGetCmd())

	return cmd
}

func newCustomFieldsValuesListCmd() *cobra.Command {
	var (
		fieldTypes        []string
		customFieldID     string
		hierarchyEntityID string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List custom field values",
		Run: func(cmd *cobra.Command, args []string) {
			c, err := getClient()
			if err != nil {
				handleError(err)
			}

			params := map[string]string{
				"type":              strings.Join(fieldTypes, ","),
				"customFieldId":     customFieldID,
				"hierarchyEntityId": hierarchyEntityID,
			}

			items, err := c.GetList("/hierarchy-entities/custom-fields-values", params, limitFlag)
			if err != nil {
				handleError(err)
			}

			headers := []string{"Custom Field ID", "Entity ID", "Value"}
			if err := output.Print(outputFormat, items, headers, func() [][]string {
				rows := make([][]string, 0, len(items))
				for _, m := range items {
					rows = append(rows, []string{
						output.SafeNested(m, "customField", "id"),
						output.SafeNested(m, "hierarchyEntity", "id"),
						fmt.Sprintf("%v", m["value"]),
					})
				}
				return rows
			}); err != nil {
				handleError(err)
			}
		},
	}

	cmd.Flags().StringSliceVar(&fieldTypes, "type", nil, "Filter by type(s)")
	cmd.Flags().StringVar(&customFieldID, "custom-field-id", "", "Filter by custom field ID")
	cmd.Flags().StringVar(&hierarchyEntityID, "hierarchy-entity-id", "", "Filter by hierarchy entity ID")

	return cmd
}

func newCustomFieldsValuesGetCmd() *cobra.Command {
	var (
		customFieldID     string
		hierarchyEntityID string
	)

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a custom field value",
		Run: func(cmd *cobra.Command, args []string) {
			c, err := getClient()
			if err != nil {
				handleError(err)
			}

			params := map[string]string{
				"customFieldId":     customFieldID,
				"hierarchyEntityId": hierarchyEntityID,
			}

			body, err := c.Get("/hierarchy-entities/custom-fields-values/value", params)
			if err != nil {
				handleError(err)
			}

			if outputFormat == output.FormatJSON {
				var raw interface{}
				if err := json.Unmarshal(body, &raw); err != nil {
					handleError(err)
				}
				if err := output.PrintJSON(raw); err != nil {
					handleError(err)
				}
				return
			}

			var resp struct {
				Data map[string]interface{} `json:"data"`
			}
			if err := json.Unmarshal(body, &resp); err != nil {
				handleError(err)
			}

			data := resp.Data
			output.PrintSingleTable([][]string{
				{"Custom Field ID", output.SafeNested(data, "customField", "id")},
				{"Entity ID", output.SafeNested(data, "hierarchyEntity", "id")},
				{"Value", fmt.Sprintf("%v", data["value"])},
			})
		},
	}

	cmd.Flags().StringVar(&customFieldID, "custom-field-id", "", "Custom field ID")
	cmd.Flags().StringVar(&hierarchyEntityID, "hierarchy-entity-id", "", "Hierarchy entity ID")
	if err := cmd.MarkFlagRequired("custom-field-id"); err != nil {
		panic(err)
	}
	if err := cmd.MarkFlagRequired("hierarchy-entity-id"); err != nil {
		panic(err)
	}

	return cmd
}
