package cli

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/tabesto/productboard-cli/internal/output"
)

func newNotesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "notes",
		Short: "Manage ProductBoard notes",
	}

	cmd.AddCommand(newNotesListCmd())
	cmd.AddCommand(newNotesGetCmd())
	cmd.AddCommand(newNotesTagsCmd())
	cmd.AddCommand(newNotesLinksCmd())

	return cmd
}

func newNotesListCmd() *cobra.Command {
	var (
		dateFrom    string
		dateTo      string
		createdFrom string
		createdTo   string
		updatedFrom string
		updatedTo   string
		term        string
		featureID   string
		companyID   string
		ownerEmail  string
		source      string
		anyTag      []string
		allTags     []string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List notes",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := getClient()
			if err != nil {
				handleError(err)
				return nil
			}

			params := map[string]string{
				"dateFrom":    dateFrom,
				"dateTo":      dateTo,
				"createdFrom": createdFrom,
				"createdTo":   createdTo,
				"updatedFrom": updatedFrom,
				"updatedTo":   updatedTo,
				"term":        term,
				"featureId":   featureID,
				"companyId":   companyID,
				"ownerEmail":  ownerEmail,
				"source":      source,
			}

			if len(anyTag) > 0 {
				params["anyTag"] = strings.Join(anyTag, ",")
			}
			if len(allTags) > 0 {
				params["allTags"] = strings.Join(allTags, ",")
			}

			items, err := c.GetList("/notes", params, limitFlag)
			if err != nil {
				handleError(err)
				return nil
			}

			headers := []string{"ID", "Title", "Source", "Owner", "Created"}
			return output.Print(outputFormat, items, headers, func() [][]string {
				rows := make([][]string, 0, len(items))
				for _, item := range items {
					rows = append(rows, []string{
						output.SafeStr(item, "id"),
						output.Truncate(output.SafeStr(item, "title"), 60),
						output.SafeStr(item, "source"),
						output.SafeNested(item, "owner", "email"),
						output.SafeStr(item, "createdAt"),
					})
				}
				return rows
			})
		},
	}

	cmd.Flags().StringVar(&dateFrom, "date-from", "", "Filter by date from")
	cmd.Flags().StringVar(&dateTo, "date-to", "", "Filter by date to")
	cmd.Flags().StringVar(&createdFrom, "created-from", "", "Filter by created from")
	cmd.Flags().StringVar(&createdTo, "created-to", "", "Filter by created to")
	cmd.Flags().StringVar(&updatedFrom, "updated-from", "", "Filter by updated from")
	cmd.Flags().StringVar(&updatedTo, "updated-to", "", "Filter by updated to")
	cmd.Flags().StringVar(&term, "term", "", "Search term")
	cmd.Flags().StringVar(&featureID, "feature-id", "", "Filter by feature ID")
	cmd.Flags().StringVar(&companyID, "company-id", "", "Filter by company ID")
	cmd.Flags().StringVar(&ownerEmail, "owner-email", "", "Filter by owner email")
	cmd.Flags().StringVar(&source, "source", "", "Filter by source")
	cmd.Flags().StringSliceVar(&anyTag, "any-tag", nil, "Filter by any of these tags (comma-separated)")
	cmd.Flags().StringSliceVar(&allTags, "all-tags", nil, "Filter by all of these tags (comma-separated)")

	return cmd
}

func newNotesGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a note by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := getClient()
			if err != nil {
				handleError(err)
				return nil
			}

			item, err := c.GetSingle("/notes/" + args[0])
			if err != nil {
				handleError(err)
				return nil
			}

			if outputFormat == output.FormatJSON {
				return output.PrintJSON(item)
			}

			output.PrintSingleTable([][]string{
				{"ID", output.SafeStr(item, "id")},
				{"Title", output.SafeStr(item, "title")},
				{"Source", output.SafeStr(item, "source")},
				{"Owner", output.SafeNested(item, "owner", "email")},
				{"Created At", output.SafeStr(item, "createdAt")},
				{"Updated At", output.SafeStr(item, "updatedAt")},
				{"Content", output.Truncate(output.SafeStr(item, "content"), 200)},
			})
			return nil
		},
	}
}

func newNotesTagsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "tags <noteId>",
		Short: "List tags for a note",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := getClient()
			if err != nil {
				handleError(err)
				return nil
			}

			items, err := c.GetLinkedResources("/notes/"+args[0]+"/tags", limitFlag)
			if err != nil {
				handleError(err)
				return nil
			}

			headers := []string{"Name"}
			return output.Print(outputFormat, items, headers, func() [][]string {
				rows := make([][]string, 0, len(items))
				for _, item := range items {
					rows = append(rows, []string{
						output.SafeStr(item, "name"),
					})
				}
				return rows
			})
		},
	}
}

func newNotesLinksCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "links <noteId>",
		Short: "List links for a note",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := getClient()
			if err != nil {
				handleError(err)
				return nil
			}

			items, err := c.GetLinkedResources("/notes/"+args[0]+"/links", limitFlag)
			if err != nil {
				handleError(err)
				return nil
			}

			headers := []string{"Type", "Target ID"}
			return output.Print(outputFormat, items, headers, func() [][]string {
				rows := make([][]string, 0, len(items))
				for _, item := range items {
					rows = append(rows, []string{
						output.SafeStr(item, "type"),
						output.SafeNested(item, "target", "id"),
					})
				}
				return rows
			})
		},
	}
}
