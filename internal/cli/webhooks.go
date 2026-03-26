package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tabesto/productboard-cli/internal/output"
)

func newWebhooksCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "webhooks",
		Short: "Manage ProductBoard webhooks",
	}

	cmd.AddCommand(newWebhooksListCmd())
	cmd.AddCommand(newWebhooksGetCmd())

	return cmd
}

func newWebhooksListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all webhooks",
		Run: func(cmd *cobra.Command, args []string) {
			c, err := getClient()
			if err != nil {
				handleError(err)
			}

			items, err := c.GetList("/webhooks", nil, limitFlag)
			if err != nil {
				handleError(err)
			}

			headers := []string{"ID", "URL", "Events"}
			if err := output.Print(outputFormat, items, headers, func() [][]string {
				rows := make([][]string, 0, len(items))
				for _, item := range items {
					var events string
					if v, ok := item["events"]; ok && v != nil {
						if slice, ok := v.([]interface{}); ok {
							parts := make([]string, 0, len(slice))
							for _, e := range slice {
								parts = append(parts, fmt.Sprintf("%v", e))
							}
							events = strings.Join(parts, ", ")
						} else {
							events = output.SafeStr(item, "events")
						}
					}
					rows = append(rows, []string{
						output.SafeStr(item, "id"),
						output.SafeStr(item, "url"),
						events,
					})
				}
				return rows
			}); err != nil {
				handleError(err)
			}
		},
	}
}

func newWebhooksGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a webhook by ID",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			c, err := getClient()
			if err != nil {
				handleError(err)
			}

			item, err := c.GetSingle(fmt.Sprintf("/webhooks/%s", args[0]))
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
				{"URL", output.SafeStr(item, "url")},
				{"Events", output.SafeStr(item, "events")},
			})
		},
	}
}
