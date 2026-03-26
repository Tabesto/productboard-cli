package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tabesto/productboard-cli/internal/output"
)

func newKeyResultsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "key-results",
		Short: "Manage ProductBoard key results",
		Long:  "List and inspect ProductBoard key results.",
	}

	cmd.AddCommand(newKeyResultsListCmd())
	cmd.AddCommand(newKeyResultsGetCmd())

	return cmd
}

func newKeyResultsListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all key results",
		Long:  "Fetch a list of key results from ProductBoard.",
		Run: func(cmd *cobra.Command, args []string) {
			c, err := getClient()
			if err != nil {
				handleError(err)
			}

			items, err := c.GetList("/key-results", nil, limitFlag)
			if err != nil {
				handleError(err)
			}

			headers := []string{"ID", "Name", "Current Value", "Target Value", "Objective ID"}
			toRows := func() [][]string {
				rows := make([][]string, 0, len(items))
				for _, m := range items {
					rows = append(rows, []string{
						output.SafeStr(m, "id"),
						output.SafeStr(m, "name"),
						output.SafeStr(m, "currentValue"),
						output.SafeStr(m, "targetValue"),
						output.SafeNested(m, "objective", "id"),
					})
				}
				return rows
			}

			if err := output.Print(outputFormat, items, headers, toRows); err != nil {
				handleError(err)
			}
		},
	}
}

func newKeyResultsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a key result by ID",
		Long:  "Fetch a single key result by its ID from ProductBoard.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			c, err := getClient()
			if err != nil {
				handleError(err)
			}

			item, err := c.GetSingle(fmt.Sprintf("/key-results/%s", args[0]))
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
				{"Current Value", output.SafeStr(item, "currentValue")},
				{"Target Value", output.SafeStr(item, "targetValue")},
				{"Objective ID", output.SafeNested(item, "objective", "id")},
			})
		},
	}
}
