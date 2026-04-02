package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tabesto/productboard-cli/internal/client"
	"github.com/tabesto/productboard-cli/internal/output"
)

func newMembersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "members",
		Short: "List and inspect ProductBoard workspace members (V2 only)",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			c, err := getClient()
			if err != nil {
				handleError(err)
			}
			if !c.IsV2() {
				handleError(&client.APIError{StatusCode: 0, Message: "The 'members' command requires API V2. Use --api-version 2 or update your config.", ExitCode: client.ExitInvalidInput})
			}
		},
	}

	cmd.AddCommand(newMembersListCmd())
	cmd.AddCommand(newMembersGetCmd())

	return cmd
}

func newMembersListCmd() *cobra.Command {
	var (
		role  string
		query string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List workspace members",
		Run: func(cmd *cobra.Command, args []string) {
			c, err := getClient()
			if err != nil {
				handleError(err)
			}

			params := map[string]string{
				"roles[]": role,
				"query":   query,
			}

			members, err := c.GetList("/members", params, limitFlag)
			if err != nil {
				handleError(err)
			}

			headers := []string{"ID", "Name", "Email", "Role"}
			toRows := func() [][]string {
				rows := make([][]string, 0, len(members))
				for _, m := range members {
					rows = append(rows, []string{
						output.SafeStr(m, "id"),
						output.SafeStr(m, "name"),
						output.SafeStr(m, "email"),
						output.SafeStr(m, "role"),
					})
				}
				return rows
			}

			if err := output.Print(outputFormat, members, headers, toRows); err != nil {
				handleError(err)
			}
		},
	}

	cmd.Flags().StringVar(&role, "role", "", "Filter by role (admin, maker, viewer, contributor)")
	cmd.Flags().StringVar(&query, "query", "", "Search by name or email")

	return cmd
}

func newMembersGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a workspace member by ID",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			c, err := getClient()
			if err != nil {
				handleError(err)
			}

			member, err := c.GetSingle(fmt.Sprintf("/members/%s", args[0]))
			if err != nil {
				handleError(err)
			}

			if outputFormat == output.FormatJSON {
				if err := output.PrintJSON(member); err != nil {
					handleError(err)
				}
				return
			}

			output.PrintSingleTable([][]string{
				{"ID", output.SafeStr(member, "id")},
				{"Name", output.SafeStr(member, "name")},
				{"Email", output.SafeStr(member, "email")},
				{"Role", output.SafeStr(member, "role")},
				{"Disabled", output.SafeStr(member, "disabled")},
			})
		},
	}
}
