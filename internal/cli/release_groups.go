package cli

import (
	"github.com/spf13/cobra"
	"github.com/tabesto/productboard-cli/internal/output"
)

func newReleaseGroupsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "release-groups",
		Short: "Manage ProductBoard release groups",
	}

	cmd.AddCommand(newReleaseGroupsListCmd())
	cmd.AddCommand(newReleaseGroupsGetCmd())

	return cmd
}

func newReleaseGroupsListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all release groups",
		Run: func(cmd *cobra.Command, args []string) {
			c, err := getClient()
			if err != nil {
				handleError(err)
			}

			items, err := c.GetList("/release-groups", nil, limitFlag)
			if err != nil {
				handleError(err)
			}

			headers := []string{"ID", "Name", "Description"}
			if err := output.Print(outputFormat, items, headers, func() [][]string {
				rows := make([][]string, 0, len(items))
				for _, item := range items {
					rows = append(rows, []string{
						output.SafeStr(item, "id"),
						output.SafeStr(item, "name"),
						output.Truncate(output.SafeStr(item, "description"), 60),
					})
				}
				return rows
			}); err != nil {
				handleError(err)
			}
		},
	}
}

func newReleaseGroupsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a release group by ID",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			c, err := getClient()
			if err != nil {
				handleError(err)
			}

			item, err := c.GetSingle("/release-groups/" + args[0])
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
				{"Description", output.SafeStr(item, "description")},
			})
		},
	}
}
