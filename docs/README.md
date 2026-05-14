# tuto documentation

Welcome to the **tuto** docs. These pages cover installation, the full CLI reference, writing modules, and shell prompt integration.

## Contents

1. **[Installation](installation.md)** — Requirements, `Makefile` targets, `go install`, and PATH setup.
2. **[Commands](commands.md)** — Full CLI reference for all subcommands.
3. **[Modules](modules.md)** — YAML schema, validation rules, authoring tips, and `tuto init`.
4. **[Shell Integration](shell-integration.md)** — Show the active session in your shell prompt (`tuto session shell-setup`).

---

## Command overview

tuto groups its commands into four top-level namespaces:

```
tuto session   — lifecycle of a learning session
tuto step      — navigation within the active session
tuto module    — discover, install, and remove modules
tuto init      — scaffold a new module YAML
```

### Typical workflow

```bash
tuto module list                   # browse available modules
tuto session start git-basics      # begin a session
tuto step show                     # read the current instruction
tuto step next                     # validate and advance
tuto step skip                     # skip without validating
tuto session status                # see overall progress
tuto session end                   # quit the session
```

---

## Project overview

tuto discovers tutorials under **`~/.tuto/modules/`** and under the directory given by `--modules` (default `./modules`). For each step the learner runs real shell commands locally; `tuto step next` re-runs the step's check command and compares combined stdout+stderr to the step's `expected_output` regex.

Step instructions are rendered as **markdown** in the terminal (headings, bold, fenced code blocks) via [Glamour](https://github.com/charmbracelet/glamour).

Session state lives in **`~/.tuto/state.json`**. Older releases used `~/.tuto-state.json`; that file is migrated automatically.

For a one-screen introduction, see the repository **[README.md](../README.md)**.
