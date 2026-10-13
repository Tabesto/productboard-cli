package cli

import (
	"github.com/spf13/cobra"
)

func newEntitiesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "entities",
		Short: "Create or update ProductBoard entities (requires write access)",
		Long: `Create or update entities (features, initiatives, objectives, ...) through API V2.

Writes are disabled by default. Enable with --allow-write, ` + "`pboard configure --write-enabled`" + `,
or PRODUCTBOARD_WRITE_ENABLED=1. Field names depend on your workspace configuration.`,
	}

	cmd.AddCommand(newEntitiesCreateCmd())
	cmd.AddCommand(newEntitiesUpdateCmd())

	return cmd
}

func newEntitiesCreateCmd() *cobra.Command {
	var (
		typ    string
		fields []string
		raw    string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create an entity",
		Example: `  pboard --allow-write entities create --type feature --field name="My feature"
  pboard --allow-write entities create --json '{"type":"feature","fields":{"name":"X"}}'`,
		Args: cobra.NoArgs,
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
			body, err := buildEntityCreateBody(typ, fields, raw)
			if err != nil {
				handleError(err)
				return nil
			}
			resp, err := c.Post("/entities", body)
			if err != nil {
				handleError(err)
				return nil
			}
			return printWriteResult(resp)
		},
	}

	cmd.Flags().StringVar(&typ, "type", "", "Entity type (e.g. feature, initiative); see entity configurations")
	cmd.Flags().StringArrayVar(&fields, "field", nil, "Field value as key=value (repeatable; value read as JSON when valid)")
	cmd.Flags().StringVar(&raw, "json", "", "Full request \"data\" object as JSON (advanced)")

	return cmd
}

func newEntitiesUpdateCmd() *cobra.Command {
	var (
		fields []string
		patch  string
		raw    string
	)

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update an entity",
		Long:  "Update an entity. Provide exactly one of --field (replace values), --patch (set/addItems/removeItems/clear operations) or --json.",
		Example: `  pboard --allow-write entities update <id> --field name="Renamed"
  pboard --allow-write entities update <id> --patch '[{"op":"clear","path":"owner"}]'`,
		Args: cobra.ExactArgs(1),
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
			resp, err := c.Patch("/entities/"+args[0], body)
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
