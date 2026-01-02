# bt-manage

`bt-manage` is a small CLI tool for managing Bluetooth device connections on **macOS**.

It can:

- list paired devices
- connect / disconnect by name (or prefix)
- optionally prompt you with a built-in TUI picker when input is omitted/ambiguous

## Requirements

- macOS
- [`blueutil`](https://github.com/toy/blueutil) available in your `PATH`

## Installation

### Install `blueutil`

If you use Homebrew:

```bash
brew install blueutil
```

### Build `bt-manage`

```bash
go build ./cmd/bt-manage
```

Move the resulting `bt-manage` binary to a directory in your `PATH`.

## Usage

### List devices

```bash
bt-manage list
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

### Connect

```bash
bt-manage connect <name-or-prefix>
```

- If `<name-or-prefix>` is omitted and stdin is a TTY, a TUI picker is shown.
- If multiple devices match the prefix and stdin is a TTY, the picker is shown.

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

Output formats:

```bash
bt-manage disconnect <name-or-prefix> --format tsv
bt-manage disconnect <name-or-prefix> --format json
bt-manage disconnect <name-or-prefix> --no-header
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
```

- Prints invoked `blueutil` commands to stderr.

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
