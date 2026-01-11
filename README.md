# portik

portik is a developer-friendly CLI to inspect, explain, and manage port ownership.

It helps answer:

- Who is using a port (e.g., 5432)
- Why a port is stuck or bind fails (TIME_WAIT, zombie processes, permission issues, IPv4/IPv6 bind confusion)
- Whether a port is mapped from a container and which container/service
- Whether it's safe to restart the owner (kill + restart last command)
- Who started a process (best-effort process tree and system hints)

Note: resolving sockets → PIDs can require elevated privileges on some systems (especially macOS). If you see missing PID/cmdline, try `sudo`.

---

## Install

### Build from source

```bash
git clone https://github.com/<you>/portik
cd portik
go build ./cmd/portik
./portik --help
```

## Requirements

- Go 1.22+ (use an up-to-date toolchain on macOS Apple Silicon)
- Linux: `ss` and `ps` in `PATH`
- macOS: `lsof` and `ps` in `PATH`
- Optional: `docker` in `PATH` for `--docker` features

## Quickstart

```bash
# who owns the port?
portik who 5432

# explain why it may be stuck
portik explain 5432

# include Docker port mapping (optional)
portik who 5432 --docker
portik explain 5432 --docker

# terminate owner (safe by default)
portik kill 5432

# smart restart: stop + rerun last command (safe by default)
portik restart 5432

# record & watch ownership changes
portik watch 5432 --interval 10s
portik history 5432 --since 7d
portik history 5432 --since 30d --detect-patterns

# daemon (foreground; use nohup/systemd if desired)
portik daemon --ports 5432,6379 --interval 30s --docker

# blame / process tree ("who started this?")
portik blame 5432 --docker


## Scan / Free / Reserve

### Scan ports (range/table)

Scan a list or range of ports and print a quick table of what’s in use.

```bash
# scan a range
portik scan --ports 3000-3010

# scan a mixed spec
portik scan --ports 22,80,443,5432,6379,8080

# scan with Docker mapping
portik scan --ports 3000-3010 --docker

# UDP scan
portik scan --ports 53,123 --proto udp

# JSON output (for scripts)
portik scan --ports 3000-3010 --json

# tune concurrency (default: CPU count, max 32)
portik scan --ports 3000-3999 --concurrency 16

# ask OS for an ephemeral free port (default)
portik free
# -> 51243

# find a free port within a range
portik free --ports 30000-40000

# UDP free port
portik free --proto udp

# JSON output
portik free --ports 30000-40000 --json

portik reserve --for 30s

# reserve a specific port for 2 minutes
portik reserve 5432 --for 2m

# reserve UDP port
portik reserve --proto udp --for 45s

# reserve on a specific bind address
portik reserve --bind 0.0.0.0 --for 30s

# JSON output
portik reserve --for 30s --json

## `use` — run a command on a free port automatically

`portik use` picks a free port (optionally from a range) and runs your command with:

- `PORT=<chosen>` set in the environment
- optional `{PORT}` template replacement in args (no shell)
- optional `--shell` mode to allow `$PORT` expansion

### Basic usage

# pick an ephemeral free port and run a command
portik use -- python -m http.server

# pick a free port from the range 3000-3999
portik use --ports 3000-3999 -- npm run dev

#If your command accepts a port argument, you can replace {PORT} in args:
portik use --ports 3000-3999 --template -- python -m http.server {PORT}

#Use this if your command string relies on $PORT expansion (or shell operators):
portik use --ports 3000-3999 --shell -- sh -lc 'echo "PORT=$PORT"; python -m http.server $PORT'

#UDP mode
portik use --proto udp --ports 40000-40100 --print
```

History is stored at: `~/.portik/history.json`

## Commands

- `portik who <port>` — show listeners for a port.
	- Flags: `--proto tcp|udp` (default `tcp`), `--docker`, `--json`

- `portik explain <port>` — adds diagnostics: port in use, IPv6-only hint, TIME_WAIT sockets, zombie hints, privileged port hints, docker mapping hints.

- `portik kill <port>` — graceful terminate then force kill after timeout.
	- Flags: `--timeout`, `--force`, `--yes`, `--proto`, `--docker`

- `portik restart <port>` — smart restart (captures cmdline, terminates owner, restarts detached).
	- Flags: `--timeout`, `--force`, `--yes`, `--docker`, `--container`, `--proto`

- `portik watch <port>` — poll periodically and record ownership changes to history.
	- Flags: `--interval`, `--proto`, `--docker`, `--json`

- `portik daemon` — monitor multiple ports and record history (foreground).
	- Flags: `--ports` (required), `--interval`, `--proto`, `--docker`, `--quiet`, `--json`

- `portik history <port>` — view history in a time window.
	- Flags: `--since`, `--detect-patterns`, `--json`

- `portik blame <port>` — best-effort process tree and "who started this" hints.
	- Flags: `--depth`, `--proto`, `--docker`, `--json`

## TUI (optional)

portik includes an optional interactive TUI (like `htop`, but for ports). It's not included in the default build to keep the CLI lightweight.

Build / Run (with TUI):

```bash
# run directly
go run -tags tui ./cmd/portik tui --ports 5432,6379 --interval 2s --docker

# or build a binary with TUI enabled
go build -tags tui -o portik ./cmd/portik
./portik tui --ports 5432,6379 --interval 2s --docker
```

TUI flags: `--ports`, `--interval`, `--proto`, `--docker`, `--actions`, `--force`

Keybindings (TUI):

- ↑/↓ or j/k — move selection
- Tab — toggle Who / Explain
- w — Who view
- e — Explain view
- r — refresh now
- / — filter/search
- Esc — clear filter / cancel
- q — quit

Actions (only if started with `--actions`):

- K — kill selected port owner (confirm)
- R — restart selected port owner (confirm)

Each port shows a sparkline representing ownership-change activity over the last 24 hours (best-effort).

## Examples

Who (text):

```
Port 5432/tcp
LISTEN 127.0.0.1:5432  pid=8123  user=me  postgres  /usr/lib/postgresql/...
```

Explain (high-level):

```
Summary
- [INFO] Port is in use
- [WARN] Only IPv6 listener detected (IPv4 bind confusion)
- [INFO] TIME_WAIT sockets present
```

## Platform support

- Linux: uses `ss` and `ps`
- macOS: uses `lsof` and `ps`
- Windows: builds, but kill/restart/socket inspection are not implemented yet

## Safety notes

Destructive actions (kill, restart, TUI actions) are conservative by default:

- confirmation prompts (unless `--yes`)
- refuses to act on processes not owned by your user (unless `--force`)

Use `sudo` when needed for PID/cmdline visibility.

## Contributing

See CONTRIBUTING.md. Please keep behavior conservative and backwards compatible for JSON output.

## License

MIT (see LICENSE).