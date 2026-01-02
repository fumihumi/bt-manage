package cmd

import (
	"bytes"
	"context"
	"testing"

	"github.com/fumihumi/bt-manage/internal/core"
)

type fakeBluetooth2 struct {
	devices []core.Device
	err     error
}

func (f fakeBluetooth2) List(ctx context.Context) ([]core.Device, error) {
	return append([]core.Device(nil), f.devices...), f.err
}

func (f fakeBluetooth2) Connect(ctx context.Context, address string) error { return nil }
func (f fakeBluetooth2) Disconnect(ctx context.Context, address string) error { return nil }

func TestListNamesOnlyPrintsOneNamePerLine(t *testing.T) {
	e := env{
		bluetooth: fakeBluetooth2{devices: []core.Device{
			{Name: "Keyboard", Address: "AA", Connected: true},
			{Name: "Mouse", Address: "BB", Connected: false},
		}},
		isTTY: func() bool { return false },
	}

	cmd := newListCmd(e)
	cmd.SetArgs([]string{"--names-only"})

	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error: %v", err)
	}

	got := out.String()
	want := "Keyboard\nMouse\n"
	if got != want {
		t.Fatalf("unexpected output\nwant=%q\n got=%q", want, got)
	}
}

func TestListNamesOnlyRejectsFormatAndNoHeader(t *testing.T) {
	e := env{
		bluetooth: fakeBluetooth2{devices: []core.Device{{Name: "Keyboard", Address: "AA"}}},
		isTTY:     func() bool { return false },
	}

	for _, tc := range []struct {
		name string
		args []string
	}{
		{name: "with-format", args: []string{"--names-only", "--format", "json"}},
		{name: "with-no-header", args: []string{"--names-only", "--no-header"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cmd := newListCmd(e)
			cmd.SetArgs(tc.args)
			cmd.SetOut(&bytes.Buffer{})
			cmd.SetErr(&bytes.Buffer{})
			if err := cmd.Execute(); err == nil {
				t.Fatalf("expected error, got nil")
			}
		})
	}
}
