package cli

import (
	"github.com/spf13/cobra"
	"github.com/tabesto/productboard-cli/internal/output"
)

func newFeatureStatusesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "feature-statuses",
		Short: "Manage ProductBoard feature statuses",
	}

	cmd.AddCommand(newFeatureStatusesListCmd())

	return cmd
}

func newFeatureStatusesListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all feature statuses",
		Run: func(cmd *cobra.Command, args []string) {
			c, err := getClient()
			if err != nil {
				handleError(err)
			}

			items, err := c.GetList("/feature-statuses", nil, limitFlag)
			if err != nil {
				handleError(err)
			}

			headers := []string{"ID", "Name"}
			if err := output.Print(outputFormat, items, headers, func() [][]string {
				rows := make([][]string, 0, len(items))
				for _, item := range items {
					rows = append(rows, []string{
						output.SafeStr(item, "id"),
						output.SafeStr(item, "name"),
					})
				}
				return rows
			}); err != nil {
				handleError(err)
			}
		},
	}
}
