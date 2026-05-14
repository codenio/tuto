package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/codenio/tuto/internal/paths"
	"github.com/codenio/tuto/internal/ui"
)

func newInitCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "init",
		Short: "One-time setup: create ~/.tuto layout and inject shell prompt",
		Long: `Creates ~/.tuto/config and ~/.tuto/modules/ if they do not exist,
then injects the tuto prompt into your shell config file.

Safe to run multiple times — existing files and shell snippets are not overwritten.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := paths.EnsureLayout(); err != nil {
				return err
			}
			ud, _ := paths.UserModulesDir()
			fmt.Printf("%s  Modules dir:  %s\n", ui.Success("✓"), ud)

			if err := paths.EnsureConfig(); err != nil {
				return err
			}
			cp, _ := paths.ConfigPath()
			fmt.Printf("%s  Config file:  %s\n", ui.Success("✓"), cp)

			fmt.Println()
			return runShellSetup()
		},
	}
	root.AddCommand(newInitShellSetupCmd())
	return root
}

func newInitShellSetupCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "shell-setup",
		Short: "Inject tuto prompt into your shell config (zsh, bash, fish)",
		Long: `Detects your current shell, appends the prompt snippet to your shell config
file, and prints the command to activate it in your current session.
Safe to run multiple times — will not add duplicates.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runShellSetup()
		},
	}
}
