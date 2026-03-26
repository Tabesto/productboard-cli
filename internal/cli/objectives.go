package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tabesto/productboard-cli/internal/output"
)

func newObjectivesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "objectives",
		Short: "Manage ProductBoard objectives",
		Long:  "List and inspect ProductBoard objectives, and explore their linked features and initiatives.",
	}

	cmd.AddCommand(newObjectivesListCmd())
	cmd.AddCommand(newObjectivesGetCmd())
	cmd.AddCommand(newObjectivesLinksCmd())

	return cmd
}

func newObjectivesListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all objectives",
		Long:  "Fetch a list of objectives from ProductBoard.",
		Run: func(cmd *cobra.Command, args []string) {
			c, err := getClient()
			if err != nil {
				handleError(err)
			}

			items, err := c.GetList("/objectives", nil, limitFlag)
			if err != nil {
				handleError(err)
			}

			headers := []string{"ID", "Name", "State", "Owner"}
			toRows := func() [][]string {
				rows := make([][]string, 0, len(items))
				for _, m := range items {
					rows = append(rows, []string{
						output.SafeStr(m, "id"),
						output.SafeStr(m, "name"),
						output.SafeStr(m, "state"),
						output.SafeNested(m, "owner", "email"),
					})
				}
				return rows
			}

			if err := output.Print(outputFormat, items, headers, toRows); err != nil {
				handleError(err)
			}
		},
	}
}

func newObjectivesGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get an objective by ID",
		Long:  "Fetch a single objective by its ID from ProductBoard.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			c, err := getClient()
			if err != nil {
				handleError(err)
			}

			item, err := c.GetSingle(fmt.Sprintf("/objectives/%s", args[0]))
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
				{"Owner", output.SafeNested(item, "owner", "email")},
			})
		},
	}
}

func newObjectivesLinksCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "links",
		Short: "List linked resources for an objective",
		Long:  "Browse resources linked to an objective, such as features and initiatives.",
	}

	cmd.AddCommand(newObjectivesLinksFeaturesCmd())
	cmd.AddCommand(newObjectivesLinksInitiativesCmd())

	return cmd
}

func newObjectivesLinksFeaturesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "features <id>",
		Short: "List features linked to an objective",
		Long:  "Fetch all features linked to the specified objective.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			c, err := getClient()
			if err != nil {
				handleError(err)
			}

			items, err := c.GetLinkedResources(
				fmt.Sprintf("/objectives/%s/links/features", args[0]),
				limitFlag,
			)
			if err != nil {
				handleError(err)
			}

			headers := []string{"ID", "Name"}
			toRows := func() [][]string {
				rows := make([][]string, 0, len(items))
				for _, item := range items {
					rows = append(rows, []string{
						output.SafeStr(item, "id"),
						output.SafeStr(item, "name"),
					})
				}
				return rows
			}

			if err := output.Print(outputFormat, items, headers, toRows); err != nil {
				handleError(err)
			}
		},
	}
}

func newObjectivesLinksInitiativesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "initiatives <id>",
		Short: "List initiatives linked to an objective",
		Long:  "Fetch all initiatives linked to the specified objective.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			c, err := getClient()
			if err != nil {
				handleError(err)
			}

			items, err := c.GetLinkedResources(
				fmt.Sprintf("/objectives/%s/links/initiatives", args[0]),
				limitFlag,
			)
			if err != nil {
				handleError(err)
			}

			headers := []string{"ID", "Name"}
			toRows := func() [][]string {
				rows := make([][]string, 0, len(items))
				for _, item := range items {
					rows = append(rows, []string{
						output.SafeStr(item, "id"),
						output.SafeStr(item, "name"),
					})
				}
				return rows
			}

			if err := output.Print(outputFormat, items, headers, toRows); err != nil {
				handleError(err)
			}
		},
	}
}
