package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/codenio/tuto/internal/ui"
)

// shellProfile holds detection info for a supported shell.
type shellProfile struct {
	name       string
	configFile string
	snippet    string
	sourceCmd  string
	evalHint   string
}

func detectShell() (*shellProfile, error) {
	shell := filepath.Base(os.Getenv("SHELL"))
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("cannot determine home directory: %w", err)
	}

	marker := "# tuto shell integration"

	switch shell {
	case "zsh":
		snippet := marker + "\n" +
			`RPROMPT='$(tuto session prompt 2>/dev/null)'`
		return &shellProfile{
			name:       "zsh",
			configFile: filepath.Join(home, ".zshrc"),
			snippet:    snippet,
			sourceCmd:  "source ~/.zshrc",
			evalHint:   `RPROMPT='$(tuto session prompt 2>/dev/null)'`,
		}, nil

	case "bash":
		cfg := filepath.Join(home, ".bashrc")
		if _, err := os.Stat(filepath.Join(home, ".bash_profile")); err == nil {
			cfg = filepath.Join(home, ".bash_profile")
		}
		snippet := marker + "\n" +
			`PS1='$(tuto session prompt 2>/dev/null) '$PS1`
		return &shellProfile{
			name:       "bash",
			configFile: cfg,
			snippet:    snippet,
			sourceCmd:  "source " + cfg,
			evalHint:   `PS1='$(tuto session prompt 2>/dev/null) '$PS1`,
		}, nil

	case "fish":
		cfg := filepath.Join(home, ".config", "fish", "config.fish")
		snippet := marker + "\n" +
			"function fish_right_prompt\n" +
			"    tuto session prompt 2>/dev/null\n" +
			"end"
		return &shellProfile{
			name:       "fish",
			configFile: cfg,
			snippet:    snippet,
			sourceCmd:  "source ~/.config/fish/config.fish",
			evalHint:   "function fish_right_prompt; tuto session prompt 2>/dev/null; end",
		}, nil

	default:
		return nil, fmt.Errorf("unsupported shell %q — add the prompt manually (see docs/shell-integration.md)", shell)
	}
}

// runShellSetup detects the current shell, appends the prompt snippet to the
// shell config file if not already present, and prints activation instructions.
func runShellSetup() error {
	profile, err := detectShell()
	if err != nil {
		return err
	}

	fmt.Printf("%s  Detected shell: %s\n", ui.Success("✓"), profile.name)
	fmt.Printf("%s  Config file:    %s\n", ui.Success("✓"), profile.configFile)

	existing, readErr := os.ReadFile(profile.configFile)
	if readErr != nil && !os.IsNotExist(readErr) {
		return fmt.Errorf("read config: %w", readErr)
	}
	if strings.Contains(string(existing), "# tuto shell integration") {
		fmt.Println()
		fmt.Println(ui.Muted("Shell prompt already installed — no changes made."))
		fmt.Println(ui.Muted("Activate in this session:  " + profile.evalHint))
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(profile.configFile), 0o755); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}

	f, err := os.OpenFile(profile.configFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("open config: %w", err)
	}
	defer func() { _ = f.Close() }()
	if _, err := fmt.Fprintf(f, "\n%s\n", profile.snippet); err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	fmt.Println()
	fmt.Println(ui.Success("✓  Shell prompt integration added."))
	fmt.Println()
	fmt.Println(ui.Muted("Reload your shell config:"))
	fmt.Println("   " + profile.sourceCmd)
	fmt.Println()
	fmt.Println(ui.Muted("Or activate just this session:"))
	fmt.Println("   " + profile.evalHint)
	return nil
}
