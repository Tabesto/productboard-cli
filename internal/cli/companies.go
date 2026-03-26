package cli

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tabesto/productboard-cli/internal/output"
)

func newCompaniesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "companies",
		Short: "Manage ProductBoard companies",
	}

	cmd.AddCommand(newCompaniesListCmd())
	cmd.AddCommand(newCompaniesGetCmd())
	cmd.AddCommand(newCompaniesCustomFieldsCmd())
	cmd.AddCommand(newCompaniesCustomFieldValueCmd())

	return cmd
}

func newCompaniesListCmd() *cobra.Command {
	var (
		term      string
		hasNotes  string
		featureID string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all companies",
		Run: func(cmd *cobra.Command, args []string) {
			c, err := getClient()
			if err != nil {
				handleError(err)
			}

			params := map[string]string{
				"term":      term,
				"hasNotes":  hasNotes,
				"featureId": featureID,
			}

			items, err := c.GetList("/companies", params, limitFlag)
			if err != nil {
				handleError(err)
			}

			headers := []string{"ID", "Name", "Domain"}
			if err := output.Print(outputFormat, items, headers, func() [][]string {
				rows := make([][]string, 0, len(items))
				for _, item := range items {
					rows = append(rows, []string{
						output.SafeStr(item, "id"),
						output.SafeStr(item, "name"),
						output.SafeStr(item, "domain"),
					})
				}
				return rows
			}); err != nil {
				handleError(err)
			}
		},
	}

	cmd.Flags().StringVar(&term, "term", "", "Filter by search term")
	cmd.Flags().StringVar(&hasNotes, "has-notes", "", "Filter companies that have notes (true/false)")
	cmd.Flags().StringVar(&featureID, "feature-id", "", "Filter by feature ID")

	return cmd
}

func newCompaniesGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a company by ID",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			c, err := getClient()
			if err != nil {
				handleError(err)
			}

			item, err := c.GetSingle("/companies/" + args[0])
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
				{"Domain", output.SafeStr(item, "domain")},
			})
		},
	}
}

func newCompaniesCustomFieldsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "custom-fields",
		Short: "Manage custom fields for companies",
	}

	cmd.AddCommand(newCompaniesCustomFieldsListCmd())
	cmd.AddCommand(newCompaniesCustomFieldsGetCmd())

	return cmd
}

func newCompaniesCustomFieldsListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List custom fields for companies",
		Run: func(cmd *cobra.Command, args []string) {
			c, err := getClient()
			if err != nil {
				handleError(err)
			}

			items, err := c.GetList("/companies/custom-fields", nil, limitFlag)
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
}

func newCompaniesCustomFieldsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a company custom field by ID",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			c, err := getClient()
			if err != nil {
				handleError(err)
			}

			item, err := c.GetSingle("/companies/custom-fields/" + args[0])
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

func newCompaniesCustomFieldValueCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "custom-field-value <companyId> <fieldId>",
		Short: "Get a custom field value for a company",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			c, err := getClient()
			if err != nil {
				handleError(err)
			}

			path := fmt.Sprintf("/companies/%s/custom-fields/%s/value", args[0], args[1])
			body, err := c.Get(path, nil)
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

			var resp map[string]interface{}
			if err := json.Unmarshal(body, &resp); err != nil {
				handleError(err)
			}

			data, _ := resp["data"].(map[string]interface{})
			if data == nil {
				data = resp
			}

			output.PrintSingleTable([][]string{
				{"Company ID", args[0]},
				{"Field ID", args[1]},
				{"Value", fmt.Sprintf("%v", data["value"])},
			})
		},
	}
}
