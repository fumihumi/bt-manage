package tty

import "os"

// IsInteractive returns true when stdout/stderr are terminals.
// This is a conservative check for whether an interactive UI can be shown.
func IsInteractive() bool {
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}
