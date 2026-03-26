package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

const skillFileName = "pboard.md"

func skillPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to resolve home directory: %w", err)
	}
	return filepath.Join(home, ".claude", "commands", skillFileName), nil
}

func newSkillCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "skill",
		Short: "Manage the pboard Claude Code skill",
	}

	cmd.AddCommand(newSkillInstallCmd())
	cmd.AddCommand(newSkillUninstallCmd())

	return cmd
}

func newSkillInstallCmd() *cobra.Command {
	var (
		force  bool
		dryRun bool
	)

	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install the pboard skill for Claude Code",
		Long:  "Installs a Claude Code skill at ~/.claude/commands/pboard.md that teaches agents how to use the pboard CLI.",
		RunE: func(cmd *cobra.Command, args []string) error {
			dest, err := skillPath()
			if err != nil {
				return err
			}

			if dryRun {
				fmt.Printf("Would install skill to %s\n\n%s\n", dest, skillContent)
				return nil
			}

			// Check if already exists
			if _, err := os.Stat(dest); err == nil && !force {
				fmt.Printf("Skill already installed at %s. Use --force to overwrite.\n", dest)
				return nil
			}

			// Create directory if needed
			dir := filepath.Dir(dest)
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", dir, err)
			}

			// Write skill file
			if err := os.WriteFile(dest, []byte(skillContent), 0644); err != nil {
				return fmt.Errorf("failed to write skill file: %w", err)
			}

			fmt.Printf("Skill installed to %s\n", dest)
			return nil
		},
	}

	cmd.Flags().BoolVar(&force, "force", false, "Overwrite existing skill without prompting")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be installed without writing")

	return cmd
}

func newSkillUninstallCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "uninstall",
		Short: "Remove the pboard skill from Claude Code",
		RunE: func(cmd *cobra.Command, args []string) error {
			dest, err := skillPath()
			if err != nil {
				return err
			}

			if _, err := os.Stat(dest); os.IsNotExist(err) {
				fmt.Printf("No skill found at %s. Nothing to remove.\n", dest)
				return nil
			}

			if err := os.Remove(dest); err != nil {
				return fmt.Errorf("failed to remove skill file: %w", err)
			}

			fmt.Printf("Skill removed from %s\n", dest)
			return nil
		},
	}
}
