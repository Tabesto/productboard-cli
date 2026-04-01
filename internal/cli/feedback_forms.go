package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tabesto/productboard-cli/internal/client"
	"github.com/tabesto/productboard-cli/internal/output"
)

func newFeedbackFormsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "feedback-forms",
		Short: "Manage ProductBoard feedback form configurations",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			c, err := getClient()
			if err != nil {
				return err
			}
			if c.IsV2() {
				fmt.Fprintln(cmd.ErrOrStderr(), "Error: The 'feedback-forms' command is not available with API V2. Use --api-version 1 to access this command.")
				return &client.APIError{StatusCode: 0, Message: "command not available with API V2", ExitCode: client.ExitInvalidInput}
			}
			return nil
		},
	}

	cmd.AddCommand(newFeedbackFormsListCmd())
	cmd.AddCommand(newFeedbackFormsGetCmd())

	return cmd
}

func newFeedbackFormsListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List feedback form configurations",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := getClient()
			if err != nil {
				handleError(err)
				return nil
			}

			items, err := c.GetList("/feedback-form-configurations", nil, limitFlag)
			if err != nil {
				handleError(err)
				return nil
			}

			headers := []string{"ID", "Name"}
			return output.Print(outputFormat, items, headers, func() [][]string {
				rows := make([][]string, 0, len(items))
				for _, item := range items {
					rows = append(rows, []string{
						output.SafeStr(item, "id"),
						output.SafeStr(item, "name"),
					})
				}
				return rows
			})
		},
	}
}

func newFeedbackFormsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a feedback form configuration by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := getClient()
			if err != nil {
				handleError(err)
				return nil
			}

			item, err := c.GetSingle("/feedback-form-configurations/" + args[0])
			if err != nil {
				handleError(err)
				return nil
			}

			if outputFormat == output.FormatJSON {
				return output.PrintJSON(item)
			}

			output.PrintSingleTable([][]string{
				{"ID", output.SafeStr(item, "id")},
				{"Name", output.SafeStr(item, "name")},
			})
			return nil
		},
	}
}
