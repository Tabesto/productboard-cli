package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tabesto/productboard-cli/internal/output"
)

func newFeaturesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "features",
		Short: "Manage ProductBoard features",
		Long: `List and inspect ProductBoard features, and explore their linked initiatives and objectives.

Examples:
  pboard features list
  pboard features list --status-name "In Progress" --limit 10
  pboard features get abc123-def456
  pboard features links initiatives abc123-def456`,
	}

	cmd.AddCommand(newFeaturesListCmd())
	cmd.AddCommand(newFeaturesGetCmd())
	cmd.AddCommand(newFeaturesLinksCmd())
	cmd.AddCommand(newFeaturesHealthCmd())

	return cmd
}

func newFeaturesListCmd() *cobra.Command {
	var (
		statusID    string
		statusName  string
		parentID    string
		archived    string
		ownerEmail  string
		noteID      string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List features",
		Long:  "Fetch a list of features from ProductBoard with optional filters.",
		Run: func(cmd *cobra.Command, args []string) {
			c, err := getClient()
			if err != nil {
				handleError(err)
			}

			params := map[string]string{
				"statusId":   statusID,
				"statusName": statusName,
				"parentId":   parentID,
				"archived":   archived,
				"ownerEmail": ownerEmail,
				"noteId":     noteID,
			}

			features, err := c.GetList("/features", params, limitFlag)
			if err != nil {
				handleError(err)
			}

			headers := []string{"ID", "Name", "Status", "Owner", "Archived"}
			toRows := func() [][]string {
				rows := make([][]string, 0, len(features))
				for _, f := range features {
					rows = append(rows, []string{
						output.SafeStr(f, "id"),
						output.Truncate(output.SafeStr(f, "name"), 60),
						output.SafeNested(f, "status", "name"),
						output.SafeNested(f, "owner", "email"),
						output.SafeStr(f, "archived"),
					})
				}
				return rows
			}

			if err := output.Print(outputFormat, features, headers, toRows); err != nil {
				handleError(err)
			}
		},
	}

	cmd.Flags().StringVar(&statusID, "status-id", "", "Filter by status ID")
	cmd.Flags().StringVar(&statusName, "status-name", "", "Filter by status name")
	cmd.Flags().StringVar(&parentID, "parent-id", "", "Filter by parent feature ID")
	cmd.Flags().StringVar(&archived, "archived", "", "Filter by archived state (true/false)")
	cmd.Flags().StringVar(&ownerEmail, "owner-email", "", "Filter by owner email")
	cmd.Flags().StringVar(&noteID, "note-id", "", "Filter by note ID")

	return cmd
}

func newFeaturesGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a feature by ID",
		Long:  "Fetch a single feature by its ID from ProductBoard.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			c, err := getClient()
			if err != nil {
				handleError(err)
			}

			feature, err := c.GetSingle(fmt.Sprintf("/features/%s", args[0]))
			if err != nil {
				handleError(err)
			}

			if outputFormat == output.FormatJSON {
				if err := output.PrintJSON(feature); err != nil {
					handleError(err)
				}
				return
			}

			output.PrintSingleTable([][]string{
				{"ID", output.SafeStr(feature, "id")},
				{"Name", output.SafeStr(feature, "name")},
				{"Status", output.SafeNested(feature, "status", "name")},
				{"Owner", output.SafeNested(feature, "owner", "email")},
				{"Archived", output.SafeStr(feature, "archived")},
				{"Description", output.Truncate(output.SafeStr(feature, "description"), 120)},
			})
		},
	}
}

func newFeaturesLinksCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "links",
		Short: "List linked resources for a feature",
		Long:  "Browse resources linked to a feature, such as initiatives and objectives.",
	}

	cmd.AddCommand(newFeaturesLinksInitiativesCmd())
	cmd.AddCommand(newFeaturesLinksObjectivesCmd())

	return cmd
}

func newFeaturesLinksInitiativesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "initiatives <id>",
		Short: "List initiatives linked to a feature",
		Long:  "Fetch all initiatives linked to the specified feature.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			c, err := getClient()
			if err != nil {
				handleError(err)
			}

			initiatives, err := c.GetLinkedResources(
				fmt.Sprintf("/features/%s/links/initiatives", args[0]),
				limitFlag,
			)
			if err != nil {
				handleError(err)
			}

			headers := []string{"ID", "Name"}
			toRows := func() [][]string {
				rows := make([][]string, 0, len(initiatives))
				for _, item := range initiatives {
					rows = append(rows, []string{
						output.SafeStr(item, "id"),
						output.SafeStr(item, "name"),
					})
				}
				return rows
			}

			if err := output.Print(outputFormat, initiatives, headers, toRows); err != nil {
				handleError(err)
			}
		},
	}
}

func newFeaturesLinksObjectivesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "objectives <id>",
		Short: "List objectives linked to a feature",
		Long:  "Fetch all objectives linked to the specified feature.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			c, err := getClient()
			if err != nil {
				handleError(err)
			}

			objectives, err := c.GetLinkedResources(
				fmt.Sprintf("/features/%s/links/objectives", args[0]),
				limitFlag,
			)
			if err != nil {
				handleError(err)
			}

			headers := []string{"ID", "Name"}
			toRows := func() [][]string {
				rows := make([][]string, 0, len(objectives))
				for _, item := range objectives {
					rows = append(rows, []string{
						output.SafeStr(item, "id"),
						output.SafeStr(item, "name"),
					})
				}
				return rows
			}

			if err := output.Print(outputFormat, objectives, headers, toRows); err != nil {
				handleError(err)
			}
		},
	}
}
