package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tabesto/productboard-cli/internal/output"
)

func newJiraIntegrationsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "jira-integrations",
		Short: "Manage ProductBoard Jira integrations",
	}

	cmd.AddCommand(newJiraIntegrationsListCmd())
	cmd.AddCommand(newJiraIntegrationsGetCmd())
	cmd.AddCommand(newJiraIntegrationsConnectionsCmd())

	return cmd
}

func newJiraIntegrationsListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all Jira integrations",
		Run: func(cmd *cobra.Command, args []string) {
			c, err := getClient()
			if err != nil {
				handleError(err)
			}

			items, err := c.GetList("/jira-integrations", nil, limitFlag)
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

func newJiraIntegrationsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a Jira integration by ID",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			c, err := getClient()
			if err != nil {
				handleError(err)
			}

			item, err := c.GetSingle(fmt.Sprintf("/jira-integrations/%s", args[0]))
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
			})
		},
	}
}

func newJiraIntegrationsConnectionsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "connections",
		Short: "Manage Jira integration connections",
	}

	cmd.AddCommand(newJiraIntegrationsConnectionsListCmd())
	cmd.AddCommand(newJiraIntegrationsConnectionsGetCmd())

	return cmd
}

func newJiraIntegrationsConnectionsListCmd() *cobra.Command {
	var (
		issueKey string
		issueID  string
	)

	cmd := &cobra.Command{
		Use:   "list <id>",
		Short: "List connections for a Jira integration",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			c, err := getClient()
			if err != nil {
				handleError(err)
			}

			params := map[string]string{
				"issueKey": issueKey,
				"issueId":  issueID,
			}

			items, err := c.GetList(fmt.Sprintf("/jira-integrations/%s/connections", args[0]), params, limitFlag)
			if err != nil {
				handleError(err)
			}

			headers := []string{"Feature ID", "External ID"}
			if err := output.Print(outputFormat, items, headers, func() [][]string {
				rows := make([][]string, 0, len(items))
				for _, item := range items {
					rows = append(rows, []string{
						output.SafeStr(item, "featureId"),
						output.SafeStr(item, "externalId"),
					})
				}
				return rows
			}); err != nil {
				handleError(err)
			}
		},
	}

	cmd.Flags().StringVar(&issueKey, "issue-key", "", "Filter by Jira issue key")
	cmd.Flags().StringVar(&issueID, "issue-id", "", "Filter by Jira issue ID")

	return cmd
}

func newJiraIntegrationsConnectionsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id> <featureId>",
		Short: "Get a Jira integration connection by integration ID and feature ID",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			c, err := getClient()
			if err != nil {
				handleError(err)
			}

			item, err := c.GetSingle(fmt.Sprintf("/jira-integrations/%s/connections/%s", args[0], args[1]))
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
				{"Feature ID", output.SafeStr(item, "featureId")},
				{"External ID", output.SafeStr(item, "externalId")},
			})
		},
	}
}
