package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tabesto/productboard-cli/internal/output"
)

func newInitiativesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "initiatives",
		Short: "Manage ProductBoard initiatives",
		Long:  "List and inspect ProductBoard initiatives, and explore their linked objectives and features.",
	}

	cmd.AddCommand(newInitiativesListCmd())
	cmd.AddCommand(newInitiativesGetCmd())
	cmd.AddCommand(newInitiativesLinksCmd())

	return cmd
}

func newInitiativesListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all initiatives",
		Long:  "Fetch a list of initiatives from ProductBoard.",
		Run: func(cmd *cobra.Command, args []string) {
			c, err := getClient()
			if err != nil {
				handleError(err)
			}

			items, err := c.GetList("/initiatives", nil, limitFlag)
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

func newInitiativesGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get an initiative by ID",
		Long:  "Fetch a single initiative by its ID from ProductBoard.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			c, err := getClient()
			if err != nil {
				handleError(err)
			}

			item, err := c.GetSingle(fmt.Sprintf("/initiatives/%s", args[0]))
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

func newInitiativesLinksCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "links",
		Short: "List linked resources for an initiative",
		Long:  "Browse resources linked to an initiative, such as objectives and features.",
	}

	cmd.AddCommand(newInitiativesLinksObjectivesCmd())
	cmd.AddCommand(newInitiativesLinksFeaturesCmd())

	return cmd
}

func newInitiativesLinksObjectivesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "objectives <id>",
		Short: "List objectives linked to an initiative",
		Long:  "Fetch all objectives linked to the specified initiative.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			c, err := getClient()
			if err != nil {
				handleError(err)
			}

			items, err := c.GetLinkedResources(
				fmt.Sprintf("/initiatives/%s/links/objectives", args[0]),
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

func newInitiativesLinksFeaturesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "features <id>",
		Short: "List features linked to an initiative",
		Long:  "Fetch all features linked to the specified initiative.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			c, err := getClient()
			if err != nil {
				handleError(err)
			}

			items, err := c.GetLinkedResources(
				fmt.Sprintf("/initiatives/%s/links/features", args[0]),
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
