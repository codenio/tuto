# Installation

## Requirements

- **Go 1.22 or later** — [download](https://go.dev/dl/)
- **macOS or Linux** — the check runner uses `sh -c`. Windows is supported via `cmd /C`.

---

## go install (recommended)

Requires Go 1.22+. No clone needed.

```bash
go install github.com/codenio/tuto/cmd/tuto@latest
```

The binary lands in `$GOPATH/bin` (usually `~/go/bin`). Add it to your `PATH` if not already there:

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

---

## From source

```bash
git clone https://github.com/codenio/tuto
cd tuto
make build

./bin/tuto session start git-basics
```

---

## Build options

### `make build`

Compiles the binary to **`bin/tuto`**:

```bash
make build
./bin/tuto --help
```

### `make install`

Runs `go install ./cmd/tuto`. The binary lands in:

- `$GOBIN` if set, otherwise
- `$(go env GOPATH)/bin` (usually `~/go/bin`)

### Manual copy

After `make build`, copy the binary anywhere on your `PATH`:

```bash
sudo cp bin/tuto /usr/local/bin/
```

### Clean build artifacts

```bash
make clean    # removes bin/
```

---

## First run

On first launch, tuto creates **`~/.tuto/modules/`** automatically. No further setup is needed.

```bash
tuto module list          # see bundled modules
tuto session start git-basics
```

---

## Module locations

tuto looks for tutorial YAML files in two places (user directory first):

| Location | Purpose |
|----------|---------|
| `~/.tuto/modules/` | Your personal installed modules (`session install`, hand-copied) |
| `--modules <dir>` | Project or local modules (default: `./modules`) |

Install a remote module:

```bash
tuto module install https://raw.githubusercontent.com/you/repo/main/k8s-basics.yaml
```

---

## Session state

Progress is stored in **`~/.tuto/state.json`**. If an older `~/.tuto-state.json` exists it is migrated automatically on first read.

---

## Shell prompt integration (optional)

Show the active session in your prompt — run once after installing:

```bash
tuto session shell-setup
```

See [Shell Integration](shell-integration.md) for manual setup and Starship config.
