# bt-manage

`bt-manage` is a small CLI tool for managing Bluetooth device connections on **macOS**.

It can:

- list paired devices
- connect / disconnect by name (or prefix)
- optionally prompt you with a built-in TUI picker when input is omitted/ambiguous
- **pair / repair** devices interactively (useful for flaky devices like Magic Trackpad)

## Requirements

- macOS
- [`blueutil`](https://github.com/toy/blueutil) available in your `PATH`

## Installation

### Install `blueutil`

Install [`blueutil`](https://github.com/toy/blueutil) and make sure it is available in your `PATH`.

If you use Homebrew:

```bash
brew install blueutil
```

### Install `bt-manage`

Install via `go install`:

```bash
go install github.com/fumihumi/bt-manage/cmd/bt-manage@latest
```

### Build from source

```bash
go build ./cmd/bt-manage
```

Move the resulting `bt-manage` binary to a directory in your `PATH`.

## Usage

### List devices

```bash
bt-manage list
```

List paired devices explicitly (default):

```bash
bt-manage list --paired
```

Show connected devices only:

```bash
bt-manage list --connected
bt-manage list -c
```

Show disconnected devices only:

```bash
bt-manage list --disconnected
bt-manage list -d
```

Print names only (one per line):

```bash
bt-manage list --names-only
bt-manage list -N
```

Output formats:

```bash
bt-manage list --format tsv
bt-manage list --format json
bt-manage list --format tsv --no-header
```

You can also omit `list` (fallback to list):

```bash
bt-manage
bt-manage -c
bt-manage -N
bt-manage --format json
```

### Connect

```bash
bt-manage connect <name-or-prefix>
```

- If `<name-or-prefix>` is omitted and stdin is a TTY, a TUI picker is shown.
- If multiple devices match the prefix and stdin is a TTY, the picker is shown.

Multi-select (interactive, space to toggle):

```bash
bt-manage connect --multi
bt-manage connect -m
```

Output formats:

```bash
bt-manage connect <name-or-prefix> --format tsv
bt-manage connect <name-or-prefix> --format json
bt-manage connect <name-or-prefix> --no-header
```

### Disconnect

```bash
bt-manage disconnect <name-or-prefix>
```

Multi-select (interactive, space to toggle):

```bash
bt-manage disconnect --multi
bt-manage disconnect -m
```

Output formats:

```bash
bt-manage disconnect <name-or-prefix> --format tsv
bt-manage disconnect <name-or-prefix> --format json
bt-manage disconnect <name-or-prefix> --no-header
```

### Pair (interactive)

Use this when you already unpaired the device (manually or via other tooling) and want to re-pair + connect.

```bash
bt-manage pair --interactive
```

Notes:

- Always interactive (TTY required). A streaming picker will show nearby discovered devices.
- Internally does: `inquiry` → `pair` → `connect` (with retries + connection verification).

Common options:

- `--inquiry-duration <sec>`: total scan window (default: 60)
- `--wait-connect <sec>`: wait for connection after connect (recommended: 10)
- `--max-attempts <n>`: connect retry count (default: 3)
- `--pin <pin>`: pass PIN if needed

### Repair (interactive)

Use this when the device is paired but becomes flaky (e.g. Magic Trackpad). This performs unpair + re-pair + connect.

```bash
bt-manage repair --interactive
```

Internally does: pick a *paired* device → `unpair` → `inquiry` → `pair` → `connect` (with retries + verification).

You can skip unpairing (e.g. if you already removed it elsewhere):

```bash
bt-manage repair --interactive --skip-unpair
```

### Force interactive mode

```bash
bt-manage connect --interactive
bt-manage disconnect --interactive
```

- Using `--interactive` in a non-TTY environment results in an error.

### Dry run

```bash
bt-manage connect <name-or-prefix> --dry-run
bt-manage disconnect <name-or-prefix> --dry-run
```

- Does not execute the action; prints **machine-readable output only**.

### Verbose

```bash
bt-manage --verbose list
bt-manage --verbose connect <name-or-prefix>
bt-manage --verbose pair --interactive
bt-manage --verbose repair --interactive
```

- Prints invoked `blueutil` commands to stderr.
- TUI picker runs in an alternate screen to reduce UI corruption when verbose logs are printed.

### Version

```bash
bt-manage version
```

## Development

```bash
go test ./...
```

## Known limitations

- macOS only.
- Behaviour depends on `blueutil` output (it may vary across environments).
