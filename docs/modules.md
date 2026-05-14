# Writing modules

Modules are **YAML** files discovered under **`~/.tuto/modules/`** and the `--modules` directory (default `./modules`). Only `.yaml` / `.yml` files are loaded.

## Scaffold a new module

```bash
tuto init my-tutorial          # creates my-tutorial.yaml in the current directory
```

Edit the generated file, then test it immediately:

```bash
tuto session start my-tutorial
tuto step next
```

---

## File schema

```yaml
name: My Tutorial                 # required — display name, used by session start
description: One-line summary     # required — shown in module list and session start

steps:
  - id: step-id                   # required — stable identifier shown in the UI
    instruction: |                # required — shown to the learner (markdown supported)
      ## Heading

      Instructions here. Supports **bold**, `inline code`, fenced blocks.

      ```bash
      echo "example command"
      ```

    command_to_run: echo "example command"   # required — run by `tuto step next`
    expected_output: '(?i)example'           # required — Go regex matched on stdout+stderr
```

### Field reference

| Field | Required | Notes |
|-------|----------|-------|
| `name` | Yes | Case-insensitive match for `session start`. Stored in session state. |
| `description` | Yes | Shown on `module list` and `session start`. |
| `steps` | Yes | Non-empty ordered list. |
| `id` | Yes (per step) | Shown in the UI. Use kebab-case. |
| `instruction` | Yes (per step) | Rendered as markdown in the terminal. |
| `command_to_run` | Yes (per step) | Run via `sh -c` (Unix) or `cmd /C` (Windows). |
| `expected_output` | Yes (per step) | [Go regex](https://pkg.go.dev/regexp/syntax) matched anywhere in combined stdout+stderr. |

---

## Markdown in instructions

Instructions are rendered with [Glamour](https://github.com/charmbracelet/glamour) when stdout is a TTY, and fall back to plain text when piped. You can use:

- `##` / `###` headings
- `**bold**`, `_italic_`
- `` `inline code` ``
- Fenced code blocks with language hints (` ```bash `)
- Bullet lists and numbered lists

---

## How validation works

When the learner runs **`tuto step next`**:

1. tuto compiles `expected_output` as a Go regex. Invalid patterns surface as an error immediately.
2. It runs `command_to_run` via shell with a configurable timeout (default 30 s; override with `--timeout N`).
3. It combines stdout + stderr and tests whether the regex matches **anywhere** in that string (not anchored unless your pattern uses `^` / `$`).
4. If the command exits non-zero but the output still matches the regex, the step **succeeds** (useful for tools that print to stderr).

On failure, the learner sees the actual output and the expected pattern, making debugging straightforward.

---

## Module resolution

The argument to `session start <name>` is matched case-insensitively against:

1. The module's **`name`** field in YAML, or
2. The **file stem** (filename without `.yaml` / `.yml`)

Search order: `~/.tuto/modules` first, then `--modules`. The first match wins.

```bash
# For git-basics.yaml with name: Git Basics
tuto session start git-basics
tuto session start "Git Basics"   # both work
```

---

## YAML quoting tips

| Situation | Recommended style |
|-----------|------------------|
| Regex with backslashes (`\s`, `\d`) | Single-quoted: `'\s+'` |
| Command with single quotes inside | Block scalar: `command_to_run: >-` then indent |
| Colon + space inside an unquoted value | Wrap in single quotes or use block scalar |

Example — command containing `'...'` shell quoting:

```yaml
command_to_run: >-
  docker run --rm -e MSG=hello alpine sh -c 'echo "$MSG"'
```

---

## Authoring tips

- **Narrow regexes** reduce false positives. Use `^…$` anchors when the output shape is stable.
- **Idempotent steps** — if a step creates a container or file, use `--replace`, a unique name, or a cleanup step so retries work cleanly.
- **Cross-platform** — remember that macOS vs Linux may differ in paths and tool availability.
- **Test your module** with `tuto step skip` to walk through all steps without executing commands.

---

## Bundled modules

| Module | Steps | What you learn |
|--------|-------|----------------|
| `git-basics` | 11 | init → stage → commit → log → diff |
| `bash-scripting` | 10 | variables, conditionals, loops, functions, exit codes |
| `docker-basics` | 16 | pull, run, exec, port mapping, env vars, clean-up |
| `podman-basics` | 18 | rootless containers, pods, exec, port mapping, clean-up |

---

## Sharing modules

1. Push your YAML to a public GitHub repository.
2. Add the topic **`tuto-module`** to the repo.
3. Learners discover it with `tuto module search` and install via raw URL:

```bash
tuto module search kubernetes
tuto module install https://raw.githubusercontent.com/you/repo/main/k8s-basics.yaml
```
