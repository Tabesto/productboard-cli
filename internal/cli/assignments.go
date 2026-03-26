package cli

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tabesto/productboard-cli/internal/output"
)

func newAssignmentsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "feature-release-assignments",
		Short: "Manage ProductBoard feature-release assignments",
	}

	cmd.AddCommand(newAssignmentsListCmd())
	cmd.AddCommand(newAssignmentsGetCmd())

	return cmd
}

func newAssignmentsListCmd() *cobra.Command {
	var featureID, releaseID, releaseState, endDateFrom, endDateTo string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all feature-release assignments",
		Run: func(cmd *cobra.Command, args []string) {
			c, err := getClient()
			if err != nil {
				handleError(err)
			}

			params := map[string]string{
				"featureId":    featureID,
				"releaseId":    releaseID,
				"releaseState": releaseState,
				"endDateFrom":  endDateFrom,
				"endDateTo":    endDateTo,
			}

			items, err := c.GetList("/feature-release-assignments", params, limitFlag)
			if err != nil {
				handleError(err)
			}

			headers := []string{"Feature ID", "Release ID"}
			if err := output.Print(outputFormat, items, headers, func() [][]string {
				rows := make([][]string, 0, len(items))
				for _, item := range items {
					rows = append(rows, []string{
						output.SafeNested(item, "feature", "id"),
						output.SafeNested(item, "release", "id"),
					})
				}
				return rows
			}); err != nil {
				handleError(err)
			}
		},
	}

	cmd.Flags().StringVar(&featureID, "feature-id", "", "Filter by feature ID")
	cmd.Flags().StringVar(&releaseID, "release-id", "", "Filter by release ID")
	cmd.Flags().StringVar(&releaseState, "release-state", "", "Filter by release state")
	cmd.Flags().StringVar(&endDateFrom, "end-date-from", "", "Filter by end date from (inclusive)")
	cmd.Flags().StringVar(&endDateTo, "end-date-to", "", "Filter by end date to (inclusive)")

	return cmd
}

func newAssignmentsGetCmd() *cobra.Command {
	var featureID, releaseID string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a feature-release assignment by feature ID and release ID",
		Run: func(cmd *cobra.Command, args []string) {
			c, err := getClient()
			if err != nil {
				handleError(err)
			}

			params := map[string]string{
				"featureId": featureID,
				"releaseId": releaseID,
			}

			body, err := c.Get("/feature-release-assignments/assignment", params)
			if err != nil {
				handleError(err)
			}

			var resp struct {
				Data map[string]interface{} `json:"data"`
			}
			if err := json.Unmarshal(body, &resp); err != nil {
				handleError(fmt.Errorf("failed to parse response: %w", err))
			}

			item := resp.Data

			if outputFormat == output.FormatJSON {
				if err := output.PrintJSON(item); err != nil {
					handleError(err)
				}
				return
			}

			output.PrintSingleTable([][]string{
				{"Feature ID", output.SafeNested(item, "feature", "id")},
				{"Release ID", output.SafeNested(item, "release", "id")},
			})
		},
	}

	cmd.Flags().StringVar(&featureID, "feature-id", "", "Feature ID (required)")
	cmd.Flags().StringVar(&releaseID, "release-id", "", "Release ID (required)")
	cmd.MarkFlagRequired("feature-id") //nolint:errcheck
	cmd.MarkFlagRequired("release-id") //nolint:errcheck

	return cmd
}
