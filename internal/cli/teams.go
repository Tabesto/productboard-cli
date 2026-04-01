package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tabesto/productboard-cli/internal/client"
	"github.com/tabesto/productboard-cli/internal/output"
)

func newTeamsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "teams",
		Short: "List and inspect ProductBoard teams (V2 only)",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			c, err := getClient()
			if err != nil {
				return err
			}
			if !c.IsV2() {
				fmt.Fprintln(os.Stderr, "Error: The 'teams' command requires API V2. Use --api-version 2 or update your config.")
				return &client.APIError{StatusCode: 0, Message: "command requires API V2", ExitCode: client.ExitInvalidInput}
			}
			return nil
		},
	}

	cmd.AddCommand(newTeamsListCmd())
	cmd.AddCommand(newTeamsGetCmd())

	return cmd
}

func newTeamsListCmd() *cobra.Command {
	var query string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List teams",
		Run: func(cmd *cobra.Command, args []string) {
			c, err := getClient()
			if err != nil {
				handleError(err)
			}

			params := map[string]string{
				"query": query,
			}

			teams, err := c.GetList("/teams", params, limitFlag)
			if err != nil {
				handleError(err)
			}

			headers := []string{"ID", "Name", "Handle", "Description"}
			toRows := func() [][]string {
				rows := make([][]string, 0, len(teams))
				for _, t := range teams {
					rows = append(rows, []string{
						output.SafeStr(t, "id"),
						output.SafeStr(t, "name"),
						output.SafeStr(t, "handle"),
						output.Truncate(output.SafeStr(t, "description"), 60),
					})
				}
				return rows
			}

			if err := output.Print(outputFormat, teams, headers, toRows); err != nil {
				handleError(err)
			}
		},
	}

	cmd.Flags().StringVar(&query, "query", "", "Search teams by name")

	return cmd
}

func newTeamsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a team by ID",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			c, err := getClient()
			if err != nil {
				handleError(err)
			}

			team, err := c.GetSingle(fmt.Sprintf("/teams/%s", args[0]))
			if err != nil {
				handleError(err)
			}

			if outputFormat == output.FormatJSON {
				if err := output.PrintJSON(team); err != nil {
					handleError(err)
				}
				return
			}

			output.PrintSingleTable([][]string{
				{"ID", output.SafeStr(team, "id")},
				{"Name", output.SafeStr(team, "name")},
				{"Handle", output.SafeStr(team, "handle")},
				{"Description", output.Truncate(output.SafeStr(team, "description"), 120)},
			})
		},
	}
}
