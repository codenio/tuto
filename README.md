# tuto

**tuto** is a command-line learning framework: declare step-by-step tutorials in YAML, ship them anywhere, and let learners validate their work in real-time against their own environment — no browser, no cloud VM, no account required.

```
tuto session start git-basics      # begin a tutorial
tuto step show                     # read the current instruction
tuto step next                     # validate and advance
tuto session status                # see your progress
```

---

## Why tuto?

| | tuto | Killercoda / Katacoda |
|---|---|---|
| Infrastructure | None — runs locally | Hosted VM required |
| Authoring | One YAML file | Platform account + config |
| Validation | Real commands, real env | Sandboxed VM |
| Distribution | URL or file path | Platform-locked |
| Offline | ✓ | ✗ |

---

## Features

- **YAML modules** — one file per tutorial: name, description, ordered steps
- **Real shell validation** — `step next` runs your actual check command and matches output with a regex
- **Markdown instructions** — rendered with syntax highlighting, bold, code blocks ([Glamour](https://github.com/charmbracelet/glamour))
- **Session management** — `session start / status / pause / resume / restart / stop`
- **Step navigation** — `step next / previous / skip / show`
- **Shell prompt integration** — shows `(git-basics: 3/11)` in your prompt via `tuto init`
- **Module registry** — `module search` finds GitHub repos tagged `tuto-module`
- **Remote install** — `module install <url>` fetches and installs modules from any HTTPS URL
- **Author scaffolding** — `tuto module create <name>` generates a starter module in seconds
- **Timeout protection** — `step next --timeout 60` kills hung check commands

---

## Quick start

### go install (recommended)

Requires [Go](https://go.dev/dl/) 1.22+.

```bash
go install github.com/codenio/tuto/cmd/tuto@latest
```

The binary lands in `$GOPATH/bin` (usually `~/go/bin`). Make sure that directory is on your `$PATH`.

### From source

```bash
git clone https://github.com/codenio/tuto && cd tuto
make build

./bin/tuto init                    # one-time setup (creates ~/.tuto, injects shell prompt)
./bin/tuto module list
./bin/tuto session start git-basics
./bin/tuto step show
./bin/tuto step next
```

Install system-wide from source:

```bash
make install    # places binary in $GOBIN or $GOPATH/bin
```

---

## Commands

### Session

| Command | Description |
|---------|-------------|
| `session start <name>` | Begin a tutorial session |
| `session status` | Show progress and current step |
| `session pause` | Hide shell prompt token (session preserved) |
| `session resume` | Restore shell prompt token |
| `session restart` | Restart from step 1 |
| `session stop` | Discard the session entirely |

### Step

| Command | Description |
|---------|-------------|
| `step next [--timeout N]` | Validate current step and advance |
| `step previous` / `prev` | Go back one step |
| `step skip` | Skip without validating |
| `step show` / `current` | Display current step instruction |

### Module

| Command | Description |
|---------|-------------|
| `module list` | Browse available modules |
| `module create [name]` | Scaffold a new module YAML |
| `module search [query]` | Search GitHub for community modules |
| `module install <url\|path>` | Install a module |
| `module update <url\|path>` | Update an installed module |
| `module uninstall <name>` | Remove an installed module |

### Setup

| Command | Description |
|---------|-------------|
| `init` | One-time setup: create `~/.tuto` layout and inject shell prompt |
| `init shell-setup` | Inject shell prompt only (zsh, bash, fish) |

---

## Shell prompt integration

```bash
tuto init           # detects zsh / bash / fish and writes the snippet automatically
source ~/.zshrc     # activate in the current session
```

Your prompt shows `(git-basics: 3/11)` while a session is active, nothing when paused or stopped.

---

## Writing a module

```bash
tuto module create my-tutorial     # creates my-tutorial.yaml
```

```yaml
name: My Tutorial
description: What learners will accomplish

steps:
  - id: first-step
    instruction: |
      ## First step

      Explain what to do. Supports **markdown** and fenced code blocks.

      ```bash
      echo "hello"
      ```

    command_to_run: echo "hello"
    expected_output: '(?i)hello'    # Go regex matched against stdout+stderr
```

Full schema → [docs/modules.md](docs/modules.md)

---

## ⚠️ Security — trust your module sources

`tuto step next` executes the `command_to_run` field from the module YAML verbatim via `sh -c`. This is by design — it is what allows validation against your real environment.

**A malicious module can run arbitrary commands on your machine.**

This is the same trust model as running a shell script you downloaded. Only install modules from sources you trust. The bundled modules under `modules/` are safe.

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

1. Push your YAML to a public GitHub repo
2. Add the topic `tuto-module` to the repo
3. Learners discover and install it:

```bash
tuto module search kubernetes
tuto module install https://raw.githubusercontent.com/you/repo/main/k8s-basics.yaml
```

---

## Repository layout

```
cmd/tuto/            main entry point
internal/cmd/        cobra commands
internal/tutorial/   YAML loading and module discovery
internal/paths/      ~/.tuto layout
internal/state/      session state (JSON)
internal/modstore/   module install / update / remove
internal/runner/     shell execution + regex check + timeout
internal/ui/         lipgloss + glamour rendering
internal/version/    build-time version info
modules/             bundled example modules
docs/                full documentation
```

---

## Documentation

| Doc | Description |
|-----|-------------|
| [docs/commands.md](docs/commands.md) | Full CLI reference |
| [docs/modules.md](docs/modules.md) | YAML schema and authoring guide |
| [docs/installation.md](docs/installation.md) | Build and install options |
| [docs/shell-integration.md](docs/shell-integration.md) | Prompt integration for zsh / bash / fish / Starship |

---

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) — includes the security model, PR checklist, and module authoring guidelines.

---

## License

MIT — see [LICENSE](LICENSE).
