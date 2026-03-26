package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	pbmcp "github.com/tabesto/productboard-cli/internal/mcp"
)

const claudeDesktopConfigName = "claude_desktop_config.json"

func claudeDesktopConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to resolve home directory: %w", err)
	}
	return filepath.Join(home, "Library", "Application Support", "Claude", claudeDesktopConfigName), nil
}

func newMcpCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mcp",
		Short: "MCP server for Claude Desktop",
		Long:  "Manage the pboard MCP server that exposes ProductBoard data as tools in Claude Desktop.",
	}

	cmd.AddCommand(newMcpServeCmd())
	cmd.AddCommand(newMcpInstallCmd())
	cmd.AddCommand(newMcpUninstallCmd())

	return cmd
}

func newMcpServeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Start the MCP server (stdio transport)",
		Long:  "Starts the pboard MCP server using stdio transport. This is called by Claude Desktop, not typically run directly.",
		RunE: func(cmd *cobra.Command, args []string) error {
			pbmcp.Version = Version
			return pbmcp.Serve()
		},
	}
}

func newMcpInstallCmd() *cobra.Command {
	var (
		force  bool
		dryRun bool
	)

	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install the pboard MCP server in Claude Desktop",
		Long:  "Registers the pboard MCP server in Claude Desktop's configuration file so ProductBoard tools are available in conversations.",
		RunE: func(cmd *cobra.Command, args []string) error {
			configPath, err := claudeDesktopConfigPath()
			if err != nil {
				return err
			}

			execPath, err := os.Executable()
			if err != nil {
				return fmt.Errorf("failed to resolve pboard binary path: %w", err)
			}
			execPath, err = filepath.EvalSymlinks(execPath)
			if err != nil {
				return fmt.Errorf("failed to resolve pboard binary path: %w", err)
			}

			entry := map[string]interface{}{
				"command": execPath,
				"args":    []string{"mcp", "serve"},
			}

			if dryRun {
				entryJSON, _ := json.MarshalIndent(entry, "    ", "  ")
				fmt.Printf("Would write to %s:\n\n", configPath)
				fmt.Printf("  \"pboard\": %s\n", string(entryJSON))
				return nil
			}

			// Read existing config or create empty
			var config map[string]interface{}
			data, err := os.ReadFile(configPath)
			if err != nil {
				if !os.IsNotExist(err) {
					return fmt.Errorf("failed to read config: %w", err)
				}
				config = map[string]interface{}{}
			} else {
				if err := json.Unmarshal(data, &config); err != nil {
					return fmt.Errorf("failed to parse config: %w", err)
				}
			}

			// Ensure mcpServers key exists
			servers, ok := config["mcpServers"].(map[string]interface{})
			if !ok {
				servers = map[string]interface{}{}
			}

			// Check if already installed
			if _, exists := servers["pboard"]; exists && !force {
				fmt.Printf("MCP server already installed in %s. Use --force to overwrite.\n", configPath)
				return nil
			}

			servers["pboard"] = entry
			config["mcpServers"] = servers

			// Write back
			dir := filepath.Dir(configPath)
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", dir, err)
			}

			out, err := json.MarshalIndent(config, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal config: %w", err)
			}

			if err := os.WriteFile(configPath, out, 0644); err != nil {
				return fmt.Errorf("failed to write config: %w", err)
			}

			fmt.Printf("MCP server installed in %s\n", configPath)
			fmt.Println("Restart Claude Desktop to activate the pboard tools.")
			return nil
		},
	}

	cmd.Flags().BoolVar(&force, "force", false, "Overwrite existing MCP server entry")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be written without modifying files")

	return cmd
}

func newMcpUninstallCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "uninstall",
		Short: "Remove the pboard MCP server from Claude Desktop",
		RunE: func(cmd *cobra.Command, args []string) error {
			configPath, err := claudeDesktopConfigPath()
			if err != nil {
				return err
			}

			data, err := os.ReadFile(configPath)
			if err != nil {
				if os.IsNotExist(err) {
					fmt.Printf("No Claude Desktop config found at %s. Nothing to remove.\n", configPath)
					return nil
				}
				return fmt.Errorf("failed to read config: %w", err)
			}

			var config map[string]interface{}
			if err := json.Unmarshal(data, &config); err != nil {
				return fmt.Errorf("failed to parse config: %w", err)
			}

			servers, ok := config["mcpServers"].(map[string]interface{})
			if !ok {
				fmt.Println("No MCP servers configured. Nothing to remove.")
				return nil
			}

			if _, exists := servers["pboard"]; !exists {
				fmt.Println("pboard MCP server not found in config. Nothing to remove.")
				return nil
			}

			delete(servers, "pboard")
			config["mcpServers"] = servers

			out, err := json.MarshalIndent(config, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal config: %w", err)
			}

			if err := os.WriteFile(configPath, out, 0644); err != nil {
				return fmt.Errorf("failed to write config: %w", err)
			}

			fmt.Printf("pboard MCP server removed from %s\n", configPath)
			fmt.Println("Restart Claude Desktop to apply changes.")
			return nil
		},
	}
}
