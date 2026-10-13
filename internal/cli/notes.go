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
	cmd.AddCommand(newNotesCreateCmd())
	cmd.AddCommand(newNotesUpdateCmd())

	return cmd
}

func newNotesCreateCmd() *cobra.Command {
	var (
		typ     string
		name    string
		content string
		fields  []string
		raw     string
	)

	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Create a note (requires write access)",
		Example: `  pboard --allow-write notes create --name "Customer asked for X" --content "Details"`,
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := getClient()
			if err != nil {
				handleError(err)
				return nil
			}
			if err := c.RequireWrite(); err != nil {
				handleError(err)
				return nil
			}
			body, err := buildNoteCreateBody(typ, name, content, fields, raw)
			if err != nil {
				handleError(err)
				return nil
			}
			resp, err := c.Post("/notes", body)
			if err != nil {
				handleError(err)
				return nil
			}
			return printWriteResult(resp)
		},
	}

	cmd.Flags().StringVar(&typ, "type", "textNote", "Note type: textNote or conversationNote")
	cmd.Flags().StringVar(&name, "name", "", "Note name (shortcut for --field name=...)")
	cmd.Flags().StringVar(&content, "content", "", "Note content (shortcut for --field content=...)")
	cmd.Flags().StringArrayVar(&fields, "field", nil, "Field value as key=value (repeatable; value read as JSON when valid)")
	cmd.Flags().StringVar(&raw, "json", "", "Full request \"data\" object as JSON (advanced)")

	return cmd
}

func newNotesUpdateCmd() *cobra.Command {
	var (
		fields []string
		patch  string
		raw    string
	)

	cmd := &cobra.Command{
		Use:     "update <id>",
		Short:   "Update a note (requires write access)",
		Long:    "Update a note. Provide exactly one of --field (replace values), --patch (set/addItems/removeItems/clear operations) or --json.",
		Example: `  pboard --allow-write notes update <id> --field content="Edited"`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := getClient()
			if err != nil {
				handleError(err)
				return nil
			}
			if err := c.RequireWrite(); err != nil {
				handleError(err)
				return nil
			}
			body, err := buildUpdateBody(fields, patch, raw)
			if err != nil {
				handleError(err)
				return nil
			}
			resp, err := c.Patch("/notes/"+args[0], body)
			if err != nil {
				handleError(err)
				return nil
			}
			return printWriteResult(resp)
		},
	}

	cmd.Flags().StringArrayVar(&fields, "field", nil, "Field value as key=value (repeatable)")
	cmd.Flags().StringVar(&patch, "patch", "", "JSON array of patch operations")
	cmd.Flags().StringVar(&raw, "json", "", "Full request \"data\" object as JSON (advanced)")

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
