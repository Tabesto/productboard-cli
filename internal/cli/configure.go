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
		Use:   "configure [token]",
		Short: "Configure API token for ProductBoard",
		Long: `Configure your ProductBoard API token.
The token is stored in ~/.config/pboard/config.yaml with restricted permissions (mode 600).

Pass the token as an argument for one-command setup, or omit it for an interactive prompt.

Examples:
  pboard configure pb_your_token_here
  pboard configure`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var token string

			if len(args) == 1 {
				token = strings.TrimSpace(args[0])
			} else {
				fmt.Print("Enter your ProductBoard API token: ")
				reader := bufio.NewReader(os.Stdin)
				input, err := reader.ReadString('\n')
				if err != nil {
					return fmt.Errorf("failed to read input: %w", err)
				}
				token = strings.TrimSpace(input)
			}

			if token == "" {
				return fmt.Errorf("token cannot be empty")
			}

			if err := config.WriteConfig(token, config.DefaultAPIVersion); err != nil {
				return err
			}

			path, _ := config.ConfigFilePath()
			fmt.Printf("Token saved to %s\n", path)
			fmt.Println("You can also set PRODUCTBOARD_API_TOKEN environment variable (takes precedence).")
			return nil
		},
	}
}
