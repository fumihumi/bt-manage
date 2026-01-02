package tty

import "os"

// IsInteractive returns true when both stdin and stdout are terminals.
// This is a conservative check for whether an interactive UI can be shown.
func IsInteractive() bool {
	in, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	out, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (in.Mode()&os.ModeCharDevice) != 0 && (out.Mode()&os.ModeCharDevice) != 0
}
