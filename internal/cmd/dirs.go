package cmd

import (
	"path/filepath"

	"github.com/codenio/tuto/internal/paths"
)

// moduleSearchDirs returns ~/.tuto/modules first, then the --modules directory (de-duplicated).
func moduleSearchDirs() ([]string, error) {
	u, err := paths.UserModulesDir()
	if err != nil {
		return nil, err
	}
	uabs, err := filepath.Abs(u)
	if err != nil {
		return nil, err
	}
	proj, err := filepath.Abs(filepath.Clean(modulesDir))
	if err != nil {
		return nil, err
	}
	seen := make(map[string]struct{})
	var out []string
	for _, d := range []string{uabs, proj} {
		if _, ok := seen[d]; ok {
			continue
		}
		seen[d] = struct{}{}
		out = append(out, d)
	}
	return out, nil
}
