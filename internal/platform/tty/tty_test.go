package tty

import (
	"os"
	"testing"
)

func TestIsInteractive_WithPipes_ReturnsFalse(t *testing.T) {
	origIn := os.Stdin
	origOut := os.Stdout
	t.Cleanup(func() {
		os.Stdin = origIn
		os.Stdout = origOut
	})

	rin, win, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	defer rin.Close()
	defer win.Close()

	rout, wout, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	defer rout.Close()
	defer wout.Close()

	os.Stdin = rin
	os.Stdout = wout

	if IsInteractive() {
		t.Fatalf("expected false when stdin/stdout are pipes")
	}
}
