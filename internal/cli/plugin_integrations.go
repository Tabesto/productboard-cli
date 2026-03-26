package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tabesto/productboard-cli/internal/output"
)

func newPluginIntegrationsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plugin-integrations",
		Short: "Manage ProductBoard plugin integrations",
	}

	cmd.AddCommand(newPluginIntegrationsListCmd())
	cmd.AddCommand(newPluginIntegrationsGetCmd())
	cmd.AddCommand(newPluginIntegrationsConnectionsCmd())

	return cmd
}

func newPluginIntegrationsListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all plugin integrations",
		Run: func(cmd *cobra.Command, args []string) {
			c, err := getClient()
			if err != nil {
				handleError(err)
			}

			items, err := c.GetList("/plugin-integrations", nil, limitFlag)
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

func newPluginIntegrationsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a plugin integration by ID",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			c, err := getClient()
			if err != nil {
				handleError(err)
			}

			item, err := c.GetSingle(fmt.Sprintf("/plugin-integrations/%s", args[0]))
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

func newPluginIntegrationsConnectionsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "connections",
		Short: "Manage plugin integration connections",
	}

	cmd.AddCommand(newPluginIntegrationsConnectionsListCmd())
	cmd.AddCommand(newPluginIntegrationsConnectionsGetCmd())

	return cmd
}

func newPluginIntegrationsConnectionsListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list <id>",
		Short: "List connections for a plugin integration",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			c, err := getClient()
			if err != nil {
				handleError(err)
			}

			items, err := c.GetList(fmt.Sprintf("/plugin-integrations/%s/connections", args[0]), nil, limitFlag)
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
}

func newPluginIntegrationsConnectionsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id> <featureId>",
		Short: "Get a plugin integration connection by integration ID and feature ID",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			c, err := getClient()
			if err != nil {
				handleError(err)
			}

			item, err := c.GetSingle(fmt.Sprintf("/plugin-integrations/%s/connections/%s", args[0], args[1]))
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
