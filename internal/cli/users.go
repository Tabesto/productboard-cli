package cli

import (
	"github.com/spf13/cobra"
	"github.com/tabesto/productboard-cli/internal/output"
)

func newUsersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "users",
		Short: "Manage ProductBoard users",
	}

	cmd.AddCommand(newUsersListCmd())
	cmd.AddCommand(newUsersGetCmd())

	return cmd
}

func newUsersListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all users",
		Run: func(cmd *cobra.Command, args []string) {
			c, err := getClient()
			if err != nil {
				handleError(err)
			}

			items, err := c.GetList("/users", nil, limitFlag)
			if err != nil {
				handleError(err)
			}

			headers := []string{"ID", "Email", "Name"}
			if err := output.Print(outputFormat, items, headers, func() [][]string {
				rows := make([][]string, 0, len(items))
				for _, item := range items {
					rows = append(rows, []string{
						output.SafeStr(item, "id"),
						output.SafeStr(item, "email"),
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

func newUsersGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a user by ID",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			c, err := getClient()
			if err != nil {
				handleError(err)
			}

			item, err := c.GetSingle("/users/" + args[0])
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
				{"Email", output.SafeStr(item, "email")},
				{"Name", output.SafeStr(item, "name")},
			})
		},
	}
}
