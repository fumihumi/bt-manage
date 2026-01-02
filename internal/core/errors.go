package core

import "fmt"

type ErrNotFound struct {
	Query string
}

func (e ErrNotFound) Error() string {
	if e.Query == "" {
		return "device not found"
	}
	return fmt.Sprintf("device not found: %s", e.Query)
}

type ErrAmbiguous struct {
	Query string
	Count int
}

func (e ErrAmbiguous) Error() string {
	if e.Query == "" {
		return "device selection is ambiguous"
	}
	return fmt.Sprintf("device selection is ambiguous: %s (%d matches)", e.Query, e.Count)
}

type ErrCanceled struct{}

func (e ErrCanceled) Error() string { return "canceled" }

type ErrDependencyMissing struct {
	Dependency string
}

func (e ErrDependencyMissing) Error() string {
	if e.Dependency == "" {
		return "dependency missing"
	}
	return fmt.Sprintf("dependency missing: %s", e.Dependency)
}
