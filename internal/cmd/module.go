package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/codenio/tuto/internal/modstore"
	"github.com/codenio/tuto/internal/paths"
	"github.com/codenio/tuto/internal/tutorial"
	"github.com/codenio/tuto/internal/ui"
)

// fence is triple-backtick — avoids embedding backticks inside a raw string literal.
const fence = "```"

func moduleTemplate(name string) string {
	return fmt.Sprintf(`name: %s
description: A short description of what this module teaches

steps:
  - id: step-1
    instruction: |
      ## Step 1: Your first task

      Explain what the learner should do here.
      You can use **markdown** — headings, bold, code blocks, lists.

      %sbash
      echo "hello world"
      %s

    command_to_run: echo "hello world"
    expected_output: "(?i)hello"

  - id: step-2
    instruction: |
      ## Step 2: Another task

      Continue the tutorial. The command below is what `+"`tuto step next`"+` will
      run to check your work.

    command_to_run: echo "done"
    expected_output: "(?i)done"
`, name, fence, fence)
}

func moduleCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "module",
		Short: "List, install, update, remove, search, or create modules",
	}
	root.AddCommand(
		moduleListCmd(),
		moduleCreateCmd(),
		moduleInstallCmd(),
		moduleUpdateCmd(),
		moduleUninstallCmd(),
		moduleSearchCmd(),
	)
	return root
}

func moduleCreateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create [module-name]",
		Short: "Scaffold a new module YAML in the current directory",
		Long: `Create a ready-to-edit YAML file with example steps and instructions.
Edit the file, then run 'tuto session start <module-name>' to try it out.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := "my-module"
			if len(args) > 0 {
				name = args[0]
			}
			filename := strings.ToLower(strings.ReplaceAll(name, " ", "-")) + ".yaml"
			if _, err := os.Stat(filename); err == nil {
				return fmt.Errorf("%s already exists; delete it or choose a different name", filename)
			}
			if err := os.WriteFile(filename, []byte(moduleTemplate(name)), 0o644); err != nil {
				return fmt.Errorf("create module file: %w", err)
			}
			abs, _ := filepath.Abs(filename)
			fmt.Println(ui.Success("Created: " + abs))
			fmt.Println(ui.Muted("Edit the file, then run:  tuto session start " + name))
			return nil
		},
	}
}

func moduleListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List available learning modules",
		RunE: func(cmd *cobra.Command, args []string) error {
			dirs, err := moduleSearchDirs()
			if err != nil {
				return err
			}
			paths, err := tutorial.CollectModuleFiles(dirs)
			if err != nil {
				return err
			}
			if len(paths) == 0 {
				fmt.Println(ui.Muted("No modules found. Searched: " + strings.Join(dirs, ", ")))
				return nil
			}
			fmt.Println(ui.Title("Learning modules"))
			for _, p := range paths {
				name, desc, err := tutorial.Summarize(p)
				if err != nil {
					fmt.Fprintf(cmd.ErrOrStderr(), "%s: %v\n", ui.Error("skip"), err)
					continue
				}
				line := fmt.Sprintf("%s\n  %s", ui.Box(name+"\n"+desc), ui.Muted(p))
				fmt.Println(line)
				fmt.Println()
			}
			return nil
		},
	}
}

func moduleInstallCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "install <path-or-https-url>",
		Short: "Copy or download a module YAML into ~/.tuto/modules",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := modstore.Install(args[0]); err != nil {
				return err
			}
			ud, _ := paths.UserModulesDir()
			fmt.Println(ui.Success("Installed module into " + ud))
			return nil
		},
	}
}

func moduleUpdateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update <path-or-https-url>",
		Short: "Overwrite an installed module file (same destination basename as install)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := modstore.Update(args[0]); err != nil {
				return err
			}
			ud, _ := paths.UserModulesDir()
			fmt.Println(ui.Success("Updated module in " + ud))
			return nil
		},
	}
}

func moduleUninstallCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "uninstall <module-name-or-stem>",
		Aliases: []string{"remove", "rm"},
		Short:   "Remove a module file from ~/.tuto/modules",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := modstore.Remove(args[0]); err != nil {
				return err
			}
			fmt.Println(ui.Success("Uninstalled module: " + args[0]))
			return nil
		},
	}
}

type ghSearchResult struct {
	TotalCount int `json:"total_count"`
	Items      []struct {
		FullName    string `json:"full_name"`
		Description string `json:"description"`
		HTMLURL     string `json:"html_url"`
		StarCount   int    `json:"stargazers_count"`
		Topics      []string `json:"topics"`
	} `json:"items"`
}

func moduleSearchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "search [query]",
		Short: "Search GitHub for community modules (topic: tuto-module)",
		Long: `Search GitHub repositories tagged with the 'tuto-module' topic.
Publish your own module by tagging your repo with that topic.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			q := "topic:tuto-module"
			if len(args) > 0 && strings.TrimSpace(args[0]) != "" {
				q += "+" + url.QueryEscape(args[0])
			}
			apiURL := "https://api.github.com/search/repositories?q=" + q + "&sort=stars&per_page=10"

			client := &http.Client{Timeout: 15 * time.Second}
			req, err := http.NewRequest(http.MethodGet, apiURL, nil)
			if err != nil {
				return fmt.Errorf("build request: %w", err)
			}
			req.Header.Set("Accept", "application/vnd.github+json")
			req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

			resp, err := client.Do(req)
			if err != nil {
				return fmt.Errorf("search GitHub: %w", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusTooManyRequests {
				return fmt.Errorf("GitHub API rate limit reached; try again in a minute")
			}
			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("GitHub API returned %s", resp.Status)
			}

			var result ghSearchResult
			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
				return fmt.Errorf("parse response: %w", err)
			}

			if result.TotalCount == 0 {
				fmt.Println(ui.Muted("No community modules found."))
				fmt.Println(ui.Muted("Publish yours by tagging a GitHub repo with the topic 'tuto-module'."))
				return nil
			}

			fmt.Println(ui.Title(fmt.Sprintf("Community modules (%d found)", result.TotalCount)))
			fmt.Println()
			for _, item := range result.Items {
				desc := item.Description
				if desc == "" {
					desc = "(no description)"
				}
				line := fmt.Sprintf("★ %d  %s\n%s\n%s",
					item.StarCount,
					ui.Success(item.FullName),
					ui.Muted(desc),
					ui.Muted("  "+item.HTMLURL),
				)
				fmt.Println(ui.Box(line))
			}
			fmt.Println(ui.Muted("Install a module: tuto module install <raw-yaml-url>"))
			return nil
		},
	}
}
