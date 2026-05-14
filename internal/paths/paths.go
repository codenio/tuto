package paths

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	dotDirName      = ".tuto"
	modulesDirName  = "modules"
	stateFileName   = "state.json"
	configFileName  = "config"
	legacyStateFile = ".tuto-state.json"
)

// DotDir returns ~/.tuto (expanded).
func DotDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, dotDirName), nil
}

// UserModulesDir returns ~/.tuto/modules.
func UserModulesDir() (string, error) {
	d, err := DotDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(d, modulesDirName), nil
}

// StatePath returns ~/.tuto/state.json.
func StatePath() (string, error) {
	d, err := DotDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(d, stateFileName), nil
}

// LegacyStatePath returns ~/.tuto-state.json (pre-layout migration).
func LegacyStatePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, legacyStateFile), nil
}

// ConfigPath returns ~/.tuto/config.
func ConfigPath() (string, error) {
	d, err := DotDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(d, configFileName), nil
}

// EnsureLayout creates ~/.tuto and ~/.tuto/modules if missing.
func EnsureLayout() error {
	ud, err := UserModulesDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(ud, 0o755); err != nil {
		return fmt.Errorf("create user modules dir: %w", err)
	}
	return nil
}

const defaultConfig = `# tuto configuration

# Timeout in seconds for step validation commands.
timeout: 30
`

// EnsureConfig writes ~/.tuto/config with defaults if the file does not exist.
func EnsureConfig() error {
	p, err := ConfigPath()
	if err != nil {
		return err
	}
	if _, err := os.Stat(p); err == nil {
		return nil // already exists
	}
	if err := os.WriteFile(p, []byte(defaultConfig), 0o644); err != nil {
		return fmt.Errorf("create config file: %w", err)
	}
	return nil
}
