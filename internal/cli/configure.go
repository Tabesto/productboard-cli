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
	var writeEnabled bool

	cmd := &cobra.Command{
		Use:   "configure [token]",
		Short: "Configure API token and write access for ProductBoard",
		Long: `Configure your ProductBoard API token and write access.
Settings are stored in ~/.config/pboard/config.yaml with restricted permissions (mode 600).
Existing settings not being changed are preserved.

Pass the token as an argument for one-command setup, or omit it for an interactive prompt.
Use --write-enabled / --write-enabled=false to turn the write gate on or off
(passing only that flag does not prompt for a token).

Examples:
  pboard configure pb_your_token_here
  pboard configure
  pboard configure --write-enabled
  pboard configure --write-enabled=false`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			writeFlagSet := cmd.Flags().Changed("write-enabled")
			values := map[string]any{}
			if writeFlagSet {
				values["write_enabled"] = writeEnabled
			}

			if len(args) == 0 && writeFlagSet {
				if err := config.Update(values); err != nil {
					return err
				}
				path, _ := config.ConfigFilePath()
				fmt.Printf("write_enabled=%t saved to %s\n", writeEnabled, path)
				return nil
			}

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

			values["api_token"] = token
			if err := config.Update(values); err != nil {
				return err
			}

			path, _ := config.ConfigFilePath()
			fmt.Printf("Token saved to %s\n", path)
			if writeFlagSet {
				fmt.Printf("write_enabled=%t\n", writeEnabled)
			}
			fmt.Println("You can also set PRODUCTBOARD_API_TOKEN environment variable (takes precedence).")
			return nil
		},
	}

	cmd.Flags().BoolVar(&writeEnabled, "write-enabled", false, "Enable (or =false disable) write operations in the saved config")

	return cmd
}
