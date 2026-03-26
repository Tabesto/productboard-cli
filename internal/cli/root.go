package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tabesto/productboard-cli/internal/client"
	"github.com/tabesto/productboard-cli/internal/config"
	"github.com/tabesto/productboard-cli/internal/output"
)

var (
	// Version is set at build time by GoReleaser.
	Version = "dev"

	outputFormat string
	limitFlag    int
)

// NewRootCmd creates the root command.
func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "pboard",
		Short: "ProductBoard CLI - read-only access to ProductBoard API",
		Long:  "pboard is a command-line tool for browsing ProductBoard data.\nIt provides read-only access to features, products, notes, releases, and more.",
		Version: Version,
	}

	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", output.FormatTable, "Output format: table or json")
	rootCmd.PersistentFlags().IntVarP(&limitFlag, "limit", "l", 0, "Maximum number of results (0 = all)")

	// Register all subcommands
	rootCmd.AddCommand(newConfigureCmd())
	rootCmd.AddCommand(newFeaturesCmd())
	rootCmd.AddCommand(newProductsCmd())
	rootCmd.AddCommand(newComponentsCmd())
	rootCmd.AddCommand(newFeatureStatusesCmd())
	rootCmd.AddCommand(newNotesCmd())
	rootCmd.AddCommand(newFeedbackFormsCmd())
	rootCmd.AddCommand(newCompaniesCmd())
	rootCmd.AddCommand(newUsersCmd())
	rootCmd.AddCommand(newReleasesCmd())
	rootCmd.AddCommand(newReleaseGroupsCmd())
	rootCmd.AddCommand(newAssignmentsCmd())
	rootCmd.AddCommand(newObjectivesCmd())
	rootCmd.AddCommand(newKeyResultsCmd())
	rootCmd.AddCommand(newInitiativesCmd())
	rootCmd.AddCommand(newCustomFieldsCmd())
	rootCmd.AddCommand(newPluginIntegrationsCmd())
	rootCmd.AddCommand(newJiraIntegrationsCmd())
	rootCmd.AddCommand(newWebhooksCmd())

	return rootCmd
}

// Execute runs the root command.
func Execute() {
	if err := NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}

// getClient loads config and creates an API client.
func getClient() (*client.Client, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	return client.New(cfg)
}

// handleError prints an error and exits with the appropriate code.
func handleError(err error) {
	if apiErr, ok := err.(*client.APIError); ok {
		fmt.Fprintf(os.Stderr, "Error: %s\n", apiErr.Message)
		os.Exit(apiErr.ExitCode)
	}
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	os.Exit(client.ExitGeneralError)
}
