# Commands

## Global flag

| Flag | Default | Description |
|------|---------|-------------|
| `--modules <dir>` | `./modules` | Extra directory of tutorial YAML files, searched **after** `~/.tuto/modules` |

The flag applies to every subcommand. Duplicate absolute paths are de-duplicated automatically.

Unless noted, commands print errors to stderr and exit non-zero on failure.

---

## `tuto session`

Manages the lifecycle of your active learning session.

### `session start <module-name>`

Starts a session on the given module. Resolves by **YAML `name`** (case-insensitive) or **file stem**. Saves `module`, `step_index: 0`, and `total_steps` to `~/.tuto/state.json`, then prints step 1.

```bash
tuto session start git-basics
tuto session start "Git Basics"
```

### `session status`

Shows the active module name, description, overall progress bar, and current step detail.

```bash
tuto session status
```

Alternatively, use `tuto show` for the current step instruction only.

### `session pause`

Hides the shell prompt token while preserving session state.

```bash
tuto session pause
```

### `session resume`

Restores the shell prompt token for a paused session.

```bash
tuto session resume
```

### `session restart` (`session reset`)

Resets `step_index` back to 0 **without** clearing the session. The module stays active; you restart from step 1.

```bash
tuto session restart
tuto session reset
```

### `session stop` (`session end`)

Clears all session state (`~/.tuto/state.json` and legacy `~/.tuto-state.json`). The module is no longer active.

```bash
tuto session stop
tuto session end
```

### `session prompt`

Prints a compact token like `(git-basics: 3/11)` when a session is active, and **nothing** when there is none. Designed to be embedded in shell prompts — see [Shell Integration](shell-integration.md).

```bash
tuto session prompt
```

---

## Step navigation

These commands operate on the active session. Requires `tuto session start <module>` first.

### `next [--timeout N]`

Validates the **current** step, then advances.

1. Runs `command_to_run` via shell and captures stdout + stderr.
2. If the combined output matches `expected_output` (Go regex), the step succeeds: `step_index` is incremented and saved.
3. Prints the next step instruction, or a completion message if all steps are done.

On failure, prints the actual output and the expected pattern. Re-run after fixing your environment.

```bash
tuto next
tuto next --timeout 60    # kill check command after 60 s (default: 30)
```

### `previous` (`prev`)

Moves `step_index` back by one, saves, and prints that step. Does not re-run validation or undo any system changes.

```bash
tuto previous
tuto prev
```

### `skip`

Advances past the current step **without** running its check command. Useful if you already completed the action or want to explore ahead.

```bash
tuto skip
```

### `show` (`current`)

Displays the current step instruction and check command without changing state.

```bash
tuto show
tuto current
```

---

## `tuto module`

Manages tutorial YAML files.

### `module list`

Lists all modules found under `~/.tuto/modules` and `--modules`. Shows name, description, and file path. Invalid files are skipped with a warning.

```bash
tuto module list
```

### `module create [module-name]`

Scaffolds a starter module YAML in the **current directory**. The file contains example steps with markdown instructions and regex validation.

```bash
tuto module create                      # creates my-module.yaml
tuto module create kubernetes-basics    # creates kubernetes-basics.yaml
```

Fails if the file already exists.

### `module search [query]`

Searches GitHub for repositories tagged with the topic **`tuto-module`**. Results are sorted by stars. Requires internet access.

```bash
tuto module search
tuto module search kubernetes
```

### `module install <path-or-https-url>`

Validates YAML, then copies it into `~/.tuto/modules/` using the source basename. Fails if the filename already exists — use `update` to overwrite.

```bash
tuto module install ./my-tutorial.yaml
tuto module install https://raw.githubusercontent.com/you/repo/main/k8s-basics.yaml
```

### `module update <path-or-https-url>`

Same as `install` but the destination file **must already exist**.

```bash
tuto module update ./my-tutorial.yaml
```

### `module remove <name>` (`rm`, `uninstall`)

Removes a module from `~/.tuto/modules` by YAML `name` or file stem. Never touches files in your `--modules` directory.

```bash
tuto module remove git-basics
tuto module remove "Git Basics"
```

---

## `tuto init`

One-time setup: creates the `~/.tuto` layout and injects the shell prompt snippet.

```bash
tuto init
tuto init shell-setup    # inject shell prompt only (zsh, bash, fish)
```

---

## Built-in help

```bash
tuto --help
tuto session --help
tuto next --help
tuto module --help
tuto <command> --help
```
