package state

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
	"github.com/codenio/tuto/internal/paths"
)

// Progress tracks the active module and step.
type Progress struct {
	ModuleName string `mapstructure:"module"       json:"module"`
	StepIndex  int    `mapstructure:"step_index"   json:"step_index"`
	TotalSteps int    `mapstructure:"total_steps"  json:"total_steps"`
	Paused     bool   `mapstructure:"paused"       json:"paused"`
}

// Load reads ~/.tuto/state.json via Viper. Migrates legacy ~/.tuto-state.json once if present.
func Load() (*Progress, error) {
	newPath, err := paths.StatePath()
	if err != nil {
		return nil, err
	}
	p, err := readStateFile(newPath)
	if err != nil {
		return nil, err
	}
	if p != nil {
		return p, nil
	}
	legacyPath, err := paths.LegacyStatePath()
	if err != nil {
		return nil, err
	}
	p, err = readStateFile(legacyPath)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, nil
	}
	if err := paths.EnsureLayout(); err != nil {
		return nil, err
	}
	if err := Save(p); err != nil {
		return nil, fmt.Errorf("migrate state from legacy file: %w", err)
	}
	_ = os.Remove(legacyPath)
	return p, nil
}

func readStateFile(path string) (*Progress, error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("json")
	if err := v.ReadInConfig(); err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read state %s: %w", path, err)
	}
	var p Progress
	if err := v.Unmarshal(&p); err != nil {
		return nil, fmt.Errorf("decode state %s: %w", path, err)
	}
	return &p, nil
}

// Save writes progress to ~/.tuto/state.json.
func Save(p *Progress) error {
	if err := paths.EnsureLayout(); err != nil {
		return err
	}
	path, err := paths.StatePath()
	if err != nil {
		return err
	}
	v := viper.New()
	v.Set("module", p.ModuleName)
	v.Set("step_index", p.StepIndex)
	v.Set("total_steps", p.TotalSteps)
	v.Set("paused", p.Paused)
	return v.WriteConfigAs(path)
}

// Clear removes state files (new layout and legacy) if they exist.
func Clear() error {
	path, err := paths.StatePath()
	if err != nil {
		return err
	}
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove state: %w", err)
	}
	legacy, err := paths.LegacyStatePath()
	if err != nil {
		return err
	}
	if err := os.Remove(legacy); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove legacy state: %w", err)
	}
	return nil
}
