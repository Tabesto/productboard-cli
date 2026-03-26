package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/tabesto/productboard-cli/internal/client"
	"github.com/tabesto/productboard-cli/internal/health"
	"github.com/tabesto/productboard-cli/internal/output"
)

func newFeaturesHealthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "health",
		Short: "Feature health updates",
		Long: `View and filter health updates across ProductBoard features.

Examples:
  pboard features health list
  pboard features health list --health-status at-risk
  pboard features health list --updated-since 2025-01-01
  pboard features health get abc123-def456`,
	}

	cmd.AddCommand(newFeaturesHealthListCmd())
	cmd.AddCommand(newFeaturesHealthGetCmd())

	return cmd
}

func newFeaturesHealthListCmd() *cobra.Command {
	var (
		includeArchived bool
		includeNoHealth bool
		updatedSince    string
		updatedBefore   string
		statusFilter    string
		ownerFilter     string
		healthStatus    string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List health updates across features",
		Long:  "Fetch all features and display a consolidated health overview with optional filters.",
		Run: func(cmd *cobra.Command, args []string) {
			c, err := getClient()
			if err != nil {
				handleError(err)
			}

			// Parse date filters
			var sinceParsed, beforeParsed *time.Time
			if updatedSince != "" {
				t, err := time.Parse("2006-01-02", updatedSince)
				if err != nil {
					handleError(&client.APIError{
						StatusCode: 0,
						Message:    fmt.Sprintf("invalid --updated-since date format %q, expected YYYY-MM-DD", updatedSince),
						ExitCode:   client.ExitInvalidInput,
					})
				}
				sinceParsed = &t
			}
			if updatedBefore != "" {
				t, err := time.Parse("2006-01-02", updatedBefore)
				if err != nil {
					handleError(&client.APIError{
						StatusCode: 0,
						Message:    fmt.Sprintf("invalid --updated-before date format %q, expected YYYY-MM-DD", updatedBefore),
						ExitCode:   client.ExitInvalidInput,
					})
				}
				beforeParsed = &t
			}

			// Build API params (only archived is server-side)
			params := map[string]string{}
			if !includeArchived {
				params["archived"] = "false"
			}

			// Fetch ALL features (limit=0) for client-side filtering
			features, err := c.GetList("/features", params, 0)
			if err != nil {
				handleError(err)
			}

			// Apply client-side filters and sort
			opts := health.FilterOpts{
				IncludeNoHealth: includeNoHealth,
				IncludeArchived: includeArchived,
				UpdatedSince:    sinceParsed,
				UpdatedBefore:   beforeParsed,
				HealthStatus:    healthStatus,
				StatusName:      statusFilter,
				OwnerEmail:      ownerFilter,
			}
			filtered := health.FilterAndSort(features, opts)

			// Apply limit AFTER filtering
			finalResults := health.ApplyLimit(filtered, limitFlag)

			headers := []string{"Feature Name", "Status", "Owner", "Health", "Health Updated", "Message"}
			toRows := func() [][]string {
				rows := make([][]string, 0, len(finalResults))
				for _, f := range finalResults {
					hu := health.GetHealthUpdate(f)
					var hStatus, hUpdated, message string
					if hu != nil {
						hStatus = output.SafeStr(hu, "status")
						hUpdated = health.FormatDate(output.SafeStr(hu, "createdAt"))
						message = output.Truncate(health.StripHTML(output.SafeStr(hu, "message")), 50)
					}
					rows = append(rows, []string{
						output.Truncate(output.SafeStr(f, "name"), 40),
						output.SafeNested(f, "status", "name"),
						output.SafeNested(f, "owner", "email"),
						hStatus,
						hUpdated,
						message,
					})
				}
				return rows
			}

			if err := output.Print(outputFormat, finalResults, headers, toRows); err != nil {
				handleError(err)
			}
		},
	}

	cmd.Flags().BoolVar(&includeArchived, "include-archived", false, "Include archived features (excluded by default)")
	cmd.Flags().BoolVar(&includeNoHealth, "include-no-health", false, "Include features without health updates")
	cmd.Flags().StringVar(&updatedSince, "updated-since", "", "Show features with health updated on or after this date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&updatedBefore, "updated-before", "", "Show features with health updated before this date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&statusFilter, "status", "", "Filter by feature status name (e.g., \"In Progress\")")
	cmd.Flags().StringVar(&ownerFilter, "owner", "", "Filter by feature owner email")
	cmd.Flags().StringVar(&healthStatus, "health-status", "", "Filter by health status (on-track, at-risk, off-track)")

	return cmd
}

func newFeaturesHealthGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <feature-id>",
		Short: "Get health update for a specific feature",
		Long:  "Fetch a single feature and display its full health update details.",
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

			hu := health.GetHealthUpdate(feature)
			fields := [][]string{
				{"ID", output.SafeStr(feature, "id")},
				{"Name", output.SafeStr(feature, "name")},
				{"Status", output.SafeNested(feature, "status", "name")},
				{"Owner", output.SafeNested(feature, "owner", "email")},
			}

			if hu != nil {
				fields = append(fields,
					[]string{"Health Status", output.SafeStr(hu, "status")},
					[]string{"Health Updated", output.SafeStr(hu, "createdAt")},
					[]string{"Message", output.SafeStr(hu, "message")},
				)
			} else {
				fields = append(fields, []string{"Health Status", "(none)"})
			}

			output.PrintSingleTable(fields)
		},
	}
}
