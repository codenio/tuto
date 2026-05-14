package modstore

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/codenio/tuto/internal/paths"
	"github.com/codenio/tuto/internal/tutorial"
)

const maxDownloadBytes = 2 << 20 // 2 MiB

// Install copies or downloads src into ~/.tuto/modules/. Fails if the destination file already exists.
func Install(src string) error {
	return installOrUpdate(src, false)
}

// Update is like Install but requires the destination file to already exist (overwrite).
func Update(src string) error {
	return installOrUpdate(src, true)
}

func installOrUpdate(src string, overwrite bool) error {
	if err := paths.EnsureLayout(); err != nil {
		return err
	}
	destDir, err := paths.UserModulesDir()
	if err != nil {
		return err
	}
	data, base, err := readSource(src)
	if err != nil {
		return err
	}
	if _, err := tutorial.LoadFromBytes(data); err != nil {
		return fmt.Errorf("invalid module YAML: %w", err)
	}
	base = sanitizeFilename(base)
	if base == "" {
		return fmt.Errorf("could not determine destination filename for %q", src)
	}
	dest := filepath.Join(destDir, base)
	if _, err := os.Stat(dest); err == nil && !overwrite {
		return fmt.Errorf("module file already exists: %s (use `tuto module update` to overwrite)", dest)
	}
	if overwrite {
		if _, err := os.Stat(dest); err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("no installed file to update: %s (use `tuto module install` first)", dest)
			}
			return err
		}
	}
	if err := os.WriteFile(dest, data, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", dest, err)
	}
	return nil
}

// Remove deletes an installed module from ~/.tuto/modules by YAML name or file stem.
func Remove(match string) error {
	destDir, err := paths.UserModulesDir()
	if err != nil {
		return err
	}
	p, err := findInUserDir(destDir, match)
	if err != nil {
		return err
	}
	if err := os.Remove(p); err != nil {
		return fmt.Errorf("remove %s: %w", p, err)
	}
	return nil
}

func findInUserDir(userDir, moduleName string) (string, error) {
	ps, err := tutorial.ListModuleFilesIn(userDir)
	if err != nil {
		return "", err
	}
	want := strings.ToLower(strings.TrimSpace(moduleName))
	if want == "" {
		return "", fmt.Errorf("empty module name")
	}
	for _, p := range ps {
		m, err := tutorial.Load(p)
		if err != nil {
			continue
		}
		stem := strings.ToLower(strings.TrimSuffix(strings.TrimSuffix(filepath.Base(p), ".yaml"), ".yml"))
		if strings.ToLower(m.Name) == want || stem == want {
			return p, nil
		}
	}
	return "", fmt.Errorf("no installed module named %q under %s", moduleName, userDir)
}

func readSource(src string) (data []byte, filename string, err error) {
	if strings.HasPrefix(strings.ToLower(src), "http://") || strings.HasPrefix(strings.ToLower(src), "https://") {
		return fetchURL(src)
	}
	src = filepath.Clean(src)
	st, err := os.Stat(src)
	if err != nil {
		return nil, "", fmt.Errorf("read source: %w", err)
	}
	if st.IsDir() {
		return nil, "", fmt.Errorf("source is a directory: %s", src)
	}
	data, err = os.ReadFile(src)
	if err != nil {
		return nil, "", err
	}
	return data, filepath.Base(src), nil
}

func fetchURL(raw string) ([]byte, string, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return nil, "", err
	}
	base := pathBaseOrDefault(u)
	if base == "" {
		return nil, "", fmt.Errorf("could not derive filename from URL %q", raw)
	}
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Get(raw)
	if err != nil {
		return nil, "", fmt.Errorf("download: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("download: HTTP %s", resp.Status)
	}
	lim := io.LimitReader(resp.Body, maxDownloadBytes+1)
	data, err := io.ReadAll(lim)
	if err != nil {
		return nil, "", err
	}
	if int64(len(data)) > maxDownloadBytes {
		return nil, "", fmt.Errorf("download exceeds %d bytes", maxDownloadBytes)
	}
	return data, base, nil
}

func pathBaseOrDefault(u *url.URL) string {
	b := filepath.Base(u.Path)
	if b == "/" || b == "." {
		return ""
	}
	if ext := strings.ToLower(filepath.Ext(b)); ext == ".yaml" || ext == ".yml" {
		return b
	}
	return ""
}

func sanitizeFilename(base string) string {
	base = filepath.Base(base)
	if base == "." || base == "/" {
		return ""
	}
	return base
}
