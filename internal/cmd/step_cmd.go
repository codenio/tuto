package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/codenio/tuto/internal/runner"
	"github.com/codenio/tuto/internal/state"
	"github.com/codenio/tuto/internal/tutorial"
	"github.com/codenio/tuto/internal/ui"
)

func newStepGroupCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "step",
		Short: "Navigate and inspect steps in the session using next, previous, skip and show",
	}
	root.AddCommand(
		newStepNextCmd(),
		newStepPreviousCmd(),
		newStepSkipCmd(),
		newStepShowCmd(),
	)
	return root
}

func newStepNextCmd() *cobra.Command {
	var timeoutSec int
	cmd := &cobra.Command{
		Use:   "next",
		Short: "Validate the current step and advance",
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := state.Load()
			if err != nil {
				return err
			}
			if p == nil || p.ModuleName == "" {
				return fmt.Errorf("no active session; run `tuto session start <module>` first")
			}
			if p.Paused {
				p.Paused = false
				fmt.Println(ui.Muted("Session resumed: " + p.ModuleName))
			}
			dirs, err := moduleSearchDirs()
			if err != nil {
				return err
			}
			_, m, err := tutorial.FindModulePath(dirs, p.ModuleName)
			if err != nil {
				return err
			}
			if p.StepIndex < 0 || p.StepIndex >= len(m.Steps) {
				return fmt.Errorf("invalid step index in state; try `tuto session reset` and start again")
			}
			step := m.Steps[p.StepIndex]
			fmt.Println(ui.Muted("Checking: " + step.CommandToRun))
			out, ok, err := runner.Check(step.CommandToRun, step.ExpectedOutput, time.Duration(timeoutSec)*time.Second)
			if err != nil {
				return err
			}
			if !ok {
				fmt.Println(ui.Error("✗ Validation failed — output did not match expected pattern."))
				if strings.TrimSpace(out) != "" {
					fmt.Println(ui.Box(strings.TrimSpace(out)))
				}
				fmt.Println(ui.Muted("Expected pattern: " + step.ExpectedOutput))
				fmt.Println(ui.Muted("Complete the step, then run `tuto step next` again."))
				return nil
			}
			fmt.Println(ui.Success("✓ Step completed: " + step.ID))
			p.StepIndex++
			if p.StepIndex >= len(m.Steps) {
				if err := state.Clear(); err != nil {
					return err
				}
				fmt.Println(ui.Title("Module complete"))
				fmt.Println(ui.Box("You finished all steps in \"" + m.Name + "\".\nRun `tuto session start <module>` to begin another."))
				return nil
			}
			if err := state.Save(p); err != nil {
				return fmt.Errorf("save state: %w", err)
			}
			next := m.Steps[p.StepIndex]
			ui.PrintStep(p.StepIndex+1, len(m.Steps), next.ID, next.Instruction, next.CommandToRun)
			return nil
		},
	}
	cmd.Flags().IntVar(&timeoutSec, "timeout", 30, "seconds before the check command is killed")
	return cmd
}

func newStepPreviousCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "previous",
		Aliases: []string{"prev"},
		Short:   "Go back one step without re-running validation",
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := state.Load()
			if err != nil {
				return err
			}
			if p == nil || p.ModuleName == "" {
				return fmt.Errorf("no active session; run `tuto session start <module>` first")
			}
			if p.Paused {
				p.Paused = false
				fmt.Println(ui.Muted("Session resumed: " + p.ModuleName))
			}
			dirs, err := moduleSearchDirs()
			if err != nil {
				return err
			}
			_, m, err := tutorial.FindModulePath(dirs, p.ModuleName)
			if err != nil {
				return err
			}
			if p.StepIndex == 0 {
				fmt.Println(ui.Muted("Already on the first step — nothing to go back to."))
				return nil
			}
			p.StepIndex--
			if err := state.Save(p); err != nil {
				return fmt.Errorf("save state: %w", err)
			}
			fmt.Println(ui.Title("Moved back · " + m.Name))
			cur := m.Steps[p.StepIndex]
			ui.PrintStep(p.StepIndex+1, len(m.Steps), cur.ID, cur.Instruction, cur.CommandToRun)
			fmt.Println(ui.Muted("Run `tuto step next` when ready to validate this step again."))
			return nil
		},
	}
}

func newStepSkipCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "skip",
		Short: "Skip the current step without validating it",
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := state.Load()
			if err != nil {
				return err
			}
			if p == nil || p.ModuleName == "" {
				return fmt.Errorf("no active session; run `tuto session start <module>` first")
			}
			if p.Paused {
				p.Paused = false
				fmt.Println(ui.Muted("Session resumed: " + p.ModuleName))
			}
			dirs, err := moduleSearchDirs()
			if err != nil {
				return err
			}
			_, m, err := tutorial.FindModulePath(dirs, p.ModuleName)
			if err != nil {
				return err
			}
			skipped := m.Steps[p.StepIndex]
			fmt.Println(ui.Muted("Skipped: " + skipped.ID))
			p.StepIndex++
			if p.StepIndex >= len(m.Steps) {
				if err := state.Clear(); err != nil {
					return err
				}
				fmt.Println(ui.Title("Module complete"))
				fmt.Println(ui.Box("You finished all steps in \"" + m.Name + "\".\nRun `tuto session start <module>` to begin another."))
				return nil
			}
			if err := state.Save(p); err != nil {
				return fmt.Errorf("save state: %w", err)
			}
			next := m.Steps[p.StepIndex]
			ui.PrintStep(p.StepIndex+1, len(m.Steps), next.ID, next.Instruction, next.CommandToRun)
			return nil
		},
	}
}

func newStepShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "show",
		Aliases: []string{"current"},
		Short:   "Show the current step instruction",
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := state.Load()
			if err != nil {
				return err
			}
			if p == nil || p.ModuleName == "" {
				fmt.Println(ui.Muted("No active session."))
				return nil
			}
			dirs, err := moduleSearchDirs()
			if err != nil {
				return err
			}
			_, m, err := tutorial.FindModulePath(dirs, p.ModuleName)
			if err != nil {
				return err
			}
			if p.StepIndex < 0 || p.StepIndex >= len(m.Steps) {
				return fmt.Errorf("invalid step index in state; try `tuto session reset` and start again")
			}
			fmt.Println(ui.Title("Current step · " + m.Name))
			cur := m.Steps[p.StepIndex]
			ui.PrintStep(p.StepIndex+1, len(m.Steps), cur.ID, cur.Instruction, cur.CommandToRun)
			fmt.Println(ui.Muted("Run `tuto step next` when done to validate and continue."))
			return nil
		},
	}
}