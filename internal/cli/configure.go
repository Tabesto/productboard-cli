package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tabesto/productboard-cli/internal/config"
)

func newConfigureCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "configure",
		Short: "Configure API token for ProductBoard",
		Long:  "Interactively configure your ProductBoard API token.\nThe token is stored in ~/.config/pboard/config.yaml with restricted permissions (mode 600).",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Print("Enter your ProductBoard API token: ")
			reader := bufio.NewReader(os.Stdin)
			token, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("failed to read input: %w", err)
			}

			token = strings.TrimSpace(token)
			if token == "" {
				return fmt.Errorf("token cannot be empty")
			}

			if err := config.WriteToken(token); err != nil {
				return err
			}

			path, _ := config.ConfigFilePath()
			fmt.Printf("Token saved to %s\n", path)
			fmt.Println("You can also set PRODUCTBOARD_API_TOKEN environment variable (takes precedence).")
			return nil
		},
	}
}
