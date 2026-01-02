package blueutil

import (
	"context"
	"errors"
	"os/exec"
	"testing"

	"github.com/fumihumi/bt-manage/internal/core"
)

type fakeExec struct {
	stdout []byte
	err    error
	calls  []struct {
		name string
		args []string
	}
}

func (f *fakeExec) Run(ctx context.Context, name string, args ...string) ([]byte, []byte, error) {
	f.calls = append(f.calls, struct {
		name string
		args []string
	}{name: name, args: append([]string(nil), args...)})
	return f.stdout, nil, f.err
}

func TestClient_List(t *testing.T) {
	fx := &fakeExec{stdout: []byte(`[{"address":"aa-bb-cc-dd-ee-ff","name":"X","connected":false,"paired":true}]`)}
	c := Client{Exec: fx, Bin: "blueutil"}

	devices, err := c.List(context.Background())
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(devices) != 1 {
		t.Fatalf("len=%d", len(devices))
	}
	if devices[0].Address != "aa:bb:cc:dd:ee:ff" {
		t.Fatalf("addr=%q", devices[0].Address)
	}
	if len(fx.calls) != 1 {
		t.Fatalf("calls=%d", len(fx.calls))
	}
	if fx.calls[0].args[0] != "--paired" {
		t.Fatalf("args=%v", fx.calls[0].args)
	}
}

func TestClient_DependencyMissing(t *testing.T) {
	fx := &fakeExec{err: &exec.Error{Name: "blueutil", Err: exec.ErrNotFound}}
	c := Client{Exec: fx, Bin: "blueutil"}

	err := c.Connect(context.Background(), "aa:bb:cc:dd:ee:ff")
	var dm core.ErrDependencyMissing
	if !errors.As(err, &dm) {
		t.Fatalf("expected ErrDependencyMissing, got %T: %v", err, err)
	}
}
