package blueutil

import (
	"bytes"
	"context"
	"os/exec"
	"runtime"
	"syscall"
)

type ExecPort interface {
	Run(ctx context.Context, name string, args ...string) (stdout []byte, stderr []byte, err error)
}

type OSExec struct{}

func (OSExec) Run(ctx context.Context, name string, args ...string) ([]byte, []byte, error) {
	cmd := exec.CommandContext(ctx, name, args...)

	// On macOS, isolate the child process into its own process group.
	// This reduces the chance that signals or process-group management
	// affect the parent process (bt-manage) when the child is terminated.
	if runtime.GOOS == "darwin" {
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	}

	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	err := cmd.Run()
	return outBuf.Bytes(), errBuf.Bytes(), err
}
