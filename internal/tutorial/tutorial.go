package tutorial

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Step is one checkpoint in a module.
type Step struct {
	ID             string `yaml:"id"`
	Instruction    string `yaml:"instruction"`
	CommandToRun   string `yaml:"command_to_run"`
	ExpectedOutput string `yaml:"expected_output"`
}

// Module is a tutorial loaded from YAML.
type Module struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Steps       []Step `yaml:"steps"`
}

// Load reads and parses a module YAML file.
func Load(path string) (*Module, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read module %q: %w", path, err)
	}
	return LoadFromBytesWithContext(data, path)
}

// LoadFromBytes parses and validates YAML module bytes (no path in error strings).
func LoadFromBytes(data []byte) (*Module, error) {
	return LoadFromBytesWithContext(data, "module")
}

// LoadFromBytesWithContext parses YAML and validates; context labels parse errors.
func LoadFromBytesWithContext(data []byte, context string) (*Module, error) {
	var m Module
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("parse %s: %w", context, err)
	}
	if err := validateModule(&m, context); err != nil {
		return nil, err
	}
	return &m, nil
}

func validateModule(m *Module, context string) error {
	if m.Name == "" {
		return fmt.Errorf("%s: missing name", context)
	}
	if len(m.Steps) == 0 {
		return fmt.Errorf("module %q (%s): no steps", m.Name, context)
	}
	for i := range m.Steps {
		if m.Steps[i].ID == "" {
			return fmt.Errorf("module %q: step %d missing id", m.Name, i)
		}
		if m.Steps[i].CommandToRun == "" || m.Steps[i].ExpectedOutput == "" {
			return fmt.Errorf("module %q: step %q missing command_to_run or expected_output", m.Name, m.Steps[i].ID)
		}
	}
	return nil
}

// listYAMLFiles returns paths to .yaml/.yml files in dir.
// When mustExist is true, a missing dir is returned as an error.
// When mustExist is false, a missing dir yields an empty slice.
func listYAMLFiles(dir string, mustExist bool) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			if mustExist {
				return nil, fmt.Errorf("modules directory %q does not exist; create it and add YAML tutorials", dir)
			}
			return nil, nil
		}
		return nil, err
	}
	var paths []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(e.Name()))
		if ext == ".yaml" || ext == ".yml" {
			paths = append(paths, filepath.Join(dir, e.Name()))
		}
	}
	return paths, nil
}

// ListModuleFiles returns paths to .yaml/.yml files in dir.
// Returns an error if dir does not exist.
func ListModuleFiles(dir string) ([]string, error) {
	return listYAMLFiles(dir, true)
}

// ListModuleFilesIn returns YAML paths in dir.
// A missing dir yields an empty slice without error.
func ListModuleFilesIn(dir string) ([]string, error) {
	return listYAMLFiles(dir, false)
}

// CollectModuleFiles lists YAML files from several roots (order preserved, paths de-duplicated).
func CollectModuleFiles(searchDirs []string) ([]string, error) {
	seen := make(map[string]struct{})
	var out []string
	for _, dir := range searchDirs {
		abs, err := filepath.Abs(dir)
		if err != nil {
			return nil, err
		}
		ps, err := ListModuleFilesIn(abs)
		if err != nil {
			return nil, err
		}
		for _, p := range ps {
			rp, err := filepath.Abs(p)
			if err != nil {
				return nil, err
			}
			if _, ok := seen[rp]; ok {
				continue
			}
			seen[rp] = struct{}{}
			out = append(out, rp)
		}
	}
	return out, nil
}

// Summarize returns a short line for list output (loads file for name/description).
func Summarize(path string) (name, description string, err error) {
	m, err := Load(path)
	if err != nil {
		return "", "", err
	}
	return m.Name, m.Description, nil
}
