package cmd

import (
	"errors"
	"fmt"
	"runtime"

	"github.com/fumihumi/bt-manage/internal/core"
)

const (
	exitOK               = 0
	exitGeneric          = 1
	exitUsage            = 2
	exitDependencyMissing = 3
	exitUnsupported       = 4
)

func exitCodeFor(err error) int {
	if err == nil {
		return exitOK
	}

	// platform
	if runtime.GOOS != "darwin" {
		return exitUnsupported
	}

	var dm core.ErrDependencyMissing
	if errors.As(err, &dm) {
		return exitDependencyMissing
	}

	var nf core.ErrNotFound
	if errors.As(err, &nf) {
		return exitUsage
	}
	var am core.ErrAmbiguous
	if errors.As(err, &am) {
		return exitUsage
	}
	var ce core.ErrCanceled
	if errors.As(err, &ce) {
		return exitUsage
	}

	// Local CLI-level errors like "--interactive requires a TTY".
	return exitUsage
}

func prettyError(err error) string {
	var dm core.ErrDependencyMissing
	if errors.As(err, &dm) {
		if dm.Dependency != "" {
			return fmt.Sprintf("missing dependency: %s", dm.Dependency)
		}
		return "missing dependency"
	}
	return err.Error()
}
