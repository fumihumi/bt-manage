package cmd

import (
	"errors"
	"fmt"

	"github.com/fumihumi/bt-manage/internal/core"
)

func userFacingError(err error) error {
	if err == nil {
		return nil
	}

	var nf core.ErrNotFound
	if errors.As(err, &nf) {
		if nf.Query == "" {
			return fmt.Errorf("no device selected")
		}
		return fmt.Errorf("no device matched %q", nf.Query)
	}

	var am core.ErrAmbiguous
	if errors.As(err, &am) {
		const hint = "try --exact or use --interactive to choose"
		if am.Query == "" {
			return fmt.Errorf("device selection is ambiguous (%s)", hint)
		}
		if am.Count > 0 {
			return fmt.Errorf("%q matched %d devices (%s)", am.Query, am.Count, hint)
		}
		return fmt.Errorf("%q matched multiple devices (%s)", am.Query, hint)
	}

	var ce core.ErrCanceled
	if errors.As(err, &ce) {
		return fmt.Errorf("canceled")
	}

	var dm core.ErrDependencyMissing
	if errors.As(err, &dm) {
		if dm.Dependency != "" {
			return fmt.Errorf("missing dependency: %s (install via Homebrew: brew install blueutil)", dm.Dependency)
		}
		return fmt.Errorf("missing dependency")
	}

	return err
}
