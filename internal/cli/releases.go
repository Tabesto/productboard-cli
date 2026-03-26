package cli

import (
	"github.com/spf13/cobra"
	"github.com/tabesto/productboard-cli/internal/output"
)

func newReleasesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "releases",
		Short: "Manage ProductBoard releases",
	}

	cmd.AddCommand(newReleasesListCmd())
	cmd.AddCommand(newReleasesGetCmd())

	return cmd
}

func newReleasesListCmd() *cobra.Command {
	var releaseGroupID string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all releases",
		Run: func(cmd *cobra.Command, args []string) {
			c, err := getClient()
			if err != nil {
				handleError(err)
			}

			params := map[string]string{
				"releaseGroupId": releaseGroupID,
			}

			items, err := c.GetList("/releases", params, limitFlag)
			if err != nil {
				handleError(err)
			}

			headers := []string{"ID", "Name", "State", "Start Date", "End Date"}
			if err := output.Print(outputFormat, items, headers, func() [][]string {
				rows := make([][]string, 0, len(items))
				for _, item := range items {
					rows = append(rows, []string{
						output.SafeStr(item, "id"),
						output.SafeStr(item, "name"),
						output.SafeStr(item, "state"),
						output.SafeNested(item, "timeframe", "startDate"),
						output.SafeNested(item, "timeframe", "endDate"),
					})
				}
				return rows
			}); err != nil {
				handleError(err)
			}
		},
	}

	cmd.Flags().StringVar(&releaseGroupID, "release-group-id", "", "Filter by release group ID")

	return cmd
}

func newReleasesGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a release by ID",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			c, err := getClient()
			if err != nil {
				handleError(err)
			}

			item, err := c.GetSingle("/releases/" + args[0])
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
				{"State", output.SafeStr(item, "state")},
				{"Start Date", output.SafeNested(item, "timeframe", "startDate")},
				{"End Date", output.SafeNested(item, "timeframe", "endDate")},
			})
		},
	}
}
