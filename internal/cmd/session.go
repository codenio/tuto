package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/codenio/tuto/internal/state"
	"github.com/codenio/tuto/internal/tutorial"
	"github.com/codenio/tuto/internal/ui"
)

func newSessionCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "session",
		Short: "Start, pause, resume, restart and stop the learning session",
	}
	root.AddCommand(
		newSessionStartCmd(),
		newSessionStopCmd(),
		newSessionPauseCmd(),
		newSessionResumeCmd(),
		newSessionRestartCmd(),
		newSessionStatusCmd(),
		newSessionPromptCmd(),
	)
	return root
}

func newSessionStartCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "start <module-name>",
		Short: "Begin a new session on a module",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dirs, err := moduleSearchDirs()
			if err != nil {
				return err
			}
			_, m, err := tutorial.FindModulePath(dirs, args[0])
			if err != nil {
				return err
			}
			p := &state.Progress{ModuleName: m.Name, StepIndex: 0, TotalSteps: len(m.Steps)}
			if err := state.Save(p); err != nil {
				return fmt.Errorf("save state: %w", err)
			}
			fmt.Println(ui.Title("Session started: " + m.Name))
			fmt.Println(ui.Muted(m.Description))
			fmt.Println()
			step := m.Steps[0]
			ui.PrintStep(1, len(m.Steps), step.ID, step.Instruction, step.CommandToRun)
			return nil
		},
	}
}

func newSessionStopCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "stop",
		Aliases: []string{"end"},
		Short:   "Stop and discard the active session",
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := state.Load()
			if err != nil {
				return err
			}
			if p == nil || p.ModuleName == "" {
				fmt.Println(ui.Muted("No active session."))
				return nil
			}
			moduleName := p.ModuleName
			if err := state.Clear(); err != nil {
				return err
			}
			fmt.Println(ui.Success("Session stopped: " + moduleName + "."))
			fmt.Println(ui.Muted("Run `tuto session start <module>` to begin a new one."))
			return nil
		},
	}
}

func newSessionPauseCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "pause",
		Short: "Pause the session (hides shell prompt token)",
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := state.Load()
			if err != nil {
				return err
			}
			if p == nil || p.ModuleName == "" {
				fmt.Println(ui.Muted("No active session to pause."))
				return nil
			}
			if p.Paused {
				fmt.Println(ui.Muted("Session is already paused.  Run `tuto session resume` to continue."))
				return nil
			}
			p.Paused = true
			if err := state.Save(p); err != nil {
				return fmt.Errorf("save state: %w", err)
			}
			fmt.Println(ui.Success("Session paused: " + p.ModuleName))
			fmt.Println(ui.Muted("Shell prompt token hidden.  Run `tuto session resume` to restore it."))
			return nil
		},
	}
}

func newSessionResumeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "resume",
		Short: "Resume a paused session (restores shell prompt token)",
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := state.Load()
			if err != nil {
				return err
			}
			if p == nil || p.ModuleName == "" {
				fmt.Println(ui.Muted("No session to resume.  Run `tuto session start <module>` to begin."))
				return nil
			}
			if !p.Paused {
				fmt.Println(ui.Muted("Session is already active.  Use `tuto step show` to see the current step."))
				return nil
			}
			p.Paused = false
			if err := state.Save(p); err != nil {
				return fmt.Errorf("save state: %w", err)
			}
			fmt.Println(ui.Success("Session resumed: " + p.ModuleName))
			fmt.Println(ui.Muted("Shell prompt token restored.  Run `tuto step show` to see the current step."))
			return nil
		},
	}
}

func newSessionRestartCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "restart",
		Aliases: []string{"reset"},
		Short:   "Restart the current module from step 1",
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := state.Load()
			if err != nil {
				return err
			}
			if p == nil || p.ModuleName == "" {
				fmt.Println(ui.Muted("No active session to restart."))
				return nil
			}
			moduleName := p.ModuleName
			p.StepIndex = 0
			p.Paused = false
			if err := state.Save(p); err != nil {
				return fmt.Errorf("save state: %w", err)
			}
			fmt.Println(ui.Success("Session restarted: " + moduleName + " — back to step 1."))
			fmt.Println(ui.Muted("Run `tuto step show` to see the first step."))
			return nil
		},
	}
}

// newSessionStatusCmd merges show + status: module overview + current step detail.
func newSessionStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show session progress and current step",
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := state.Load()
			if err != nil {
				return err
			}
			if p == nil || p.ModuleName == "" {
				fmt.Println(ui.Muted("No active session.  Run `tuto session start <module>` to begin."))
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
			total := len(m.Steps)

			label := "Active session"
			if p.Paused {
				label = "Session (paused)"
			}
			fmt.Println(ui.Title(label))
			fmt.Println(ui.Box(fmt.Sprintf(
				"Module:  %s\n%s\n%s",
				ui.Success(m.Name),
				ui.Muted(m.Description),
				ui.ProgressLine("Progress", p.StepIndex, total),
			)))
			fmt.Println()

			step := m.Steps[p.StepIndex]
			ui.PrintStep(p.StepIndex+1, total, step.ID, step.Instruction, step.CommandToRun)

			if p.Paused {
				fmt.Println(ui.Muted("Session is paused.  Run `tuto session resume` to continue."))
			}
			return nil
		},
	}
}

// newSessionPromptCmd is hidden — called by the shell snippet injected by `tuto init shell-setup`.
// Outputs nothing when no session is active or when the session is paused.
func newSessionPromptCmd() *cobra.Command {
	return &cobra.Command{
		Use:    "prompt",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := state.Load()
			if err != nil || p == nil || p.ModuleName == "" || p.Paused {
				return nil
			}
			slug := strings.ToLower(strings.ReplaceAll(p.ModuleName, " ", "-"))
			step := p.StepIndex + 1
			total := p.TotalSteps
			if total == 0 {
				fmt.Printf("(%s: %d)", slug, step)
			} else {
				fmt.Printf("(%s: %d/%d)", slug, step, total)
			}
			return nil
		},
	}
}
