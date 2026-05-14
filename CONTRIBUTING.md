# Contributing to tuto

Thank you for your interest in contributing!

---

## Getting started

```bash
git clone https://github.com/<you>/tuto
cd tuto
make build
./bin/tuto --help
```

Requirements: **Go 1.22+**

---

## Development workflow

```bash
make build    # compile to bin/tuto
make test     # run all tests
make lint     # run golangci-lint (install separately)
make fmt      # gofmt + goimports
```

---

## Project layout

```
cmd/tuto/           entry point
internal/
  cmd/              cobra commands (session, step, module, init)
  modstore/         module install / update / remove
  paths/            ~/.tuto directory layout
  runner/           shell command execution + regex validation
  state/            session state (JSON)
  tutorial/         YAML module loading + validation
  ui/               terminal rendering (lipgloss, glamour)
modules/            bundled tutorial YAML files
docs/               user-facing documentation
```

---

## ⚠️ Security model — read before contributing

**tuto executes shell commands from module YAML files.**

The `command_to_run` field of every step is passed verbatim to `sh -c` (Unix) or `cmd /C` (Windows). This is intentional — it is what lets tuto validate real work in a real environment.

**Consequences:**

- A malicious module can run arbitrary code on the learner's machine.
- tuto makes **no attempt** to sandbox or restrict what modules can do.
- This is the same trust model as running a shell script you downloaded.

**What this means for contributors:**

- Bundled modules (under `modules/`) must never execute destructive or network-facing commands without clear user instruction.
- The README and `tuto module install` documentation must always remind users to **only install modules from sources they trust**.
- Do not add a "run all commands automatically without validation" mode.
- Do not weaken the one-step-at-a-time validation flow.

If you discover a way a malicious module author could escalate beyond the above (e.g., exploit the runner, tamper with state, path traversal in module install), please report it privately — see **Reporting vulnerabilities** below.

---

## Writing or improving modules

- Keep steps focused: one skill per step.
- Use idempotent `command_to_run` commands where possible.
- Anchor regexes (`^...$`) when the output shape is predictable.
- Test with `tuto step skip` to walk through all steps without executing.
- See [docs/modules.md](docs/modules.md) for the full schema reference.

---

## Pull request checklist

- [ ] `make build` passes
- [ ] `make test` passes (or new tests added for new behaviour)
- [ ] `make lint` passes (no new warnings)
- [ ] Docs updated if commands / YAML schema changed
- [ ] New module YAML validated with `tuto session start <name>`

---

## Reporting vulnerabilities

Please **do not** open a public issue for security vulnerabilities.

Email the maintainers directly or use GitHub's private security advisory:
`https://github.com/<org>/tuto/security/advisories/new`

---

## Code of conduct

Be kind. Assume good intent. Disagreement on technical matters is fine;
personal attacks are not. We follow the [Contributor Covenant](https://www.contributor-covenant.org/).
