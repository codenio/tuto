package tutorial

import (
	"fmt"
	"path/filepath"
	"strings"
)

// FindModulePath returns the first module YAML whose name or file stem matches moduleName,
// searching searchDirs in order.
func FindModulePath(searchDirs []string, moduleName string) (string, *Module, error) {
	want := strings.ToLower(strings.TrimSpace(moduleName))
	for _, dir := range searchDirs {
		abs, err := filepath.Abs(dir)
		if err != nil {
			return "", nil, err
		}
		paths, err := ListModuleFilesIn(abs)
		if err != nil {
			return "", nil, err
		}
		for _, p := range paths {
			m, err := Load(p)
			if err != nil {
				continue
			}
			stem := strings.ToLower(strings.TrimSuffix(strings.TrimSuffix(filepath.Base(p), ".yaml"), ".yml"))
			if strings.ToLower(m.Name) == want || stem == want {
				return p, m, nil
			}
		}
	}
	return "", nil, fmt.Errorf("no module named %q (searched: %s)", moduleName, strings.Join(searchDirs, ", "))
}
