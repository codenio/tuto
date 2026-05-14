package cmd

import (
	"github.com/spf13/cobra"
	"github.com/codenio/tuto/internal/paths"
	"github.com/codenio/tuto/internal/version"
)

var modulesDir string

// Execute runs the CLI.
func Execute() error {
	if err := paths.EnsureLayout(); err != nil {
		return err
	}
	root := &cobra.Command{
		Use:     "tuto",
		Short:   "Interactive step-by-step terminal tutorials",
		Long:    "tuto loads YAML modules from ~/.tuto/modules and your --modules directory, and guides you through tasks with validation.",
		Version: version.String(),
	}
	root.PersistentFlags().StringVar(&modulesDir, "modules", "./modules", "additional directory containing tutorial YAML files (after ~/.tuto/modules)")

	root.AddCommand(newSessionCmd(), newStepGroupCmd(), newInitCmd(), newModuleCmd())
	return root.Execute()
}
