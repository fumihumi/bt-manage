package cmd

import (
	"bytes"
	"context"
	"testing"

	"github.com/fumihumi/bt-manage/internal/core"
)

type fakeBluetooth struct {
	devices []core.Device
	err     error
}

func (f fakeBluetooth) List(ctx context.Context) ([]core.Device, error) {
	return append([]core.Device(nil), f.devices...), f.err
}

func (f fakeBluetooth) Connect(ctx context.Context, address string) error { return nil }
func (f fakeBluetooth) Disconnect(ctx context.Context, address string) error { return nil }
func (f fakeBluetooth) Pair(ctx context.Context, address string, pin string) error { return nil }
func (f fakeBluetooth) Unpair(ctx context.Context, address string) error { return nil }
func (f fakeBluetooth) Inquiry(ctx context.Context, durationSeconds int) ([]core.Device, error) {
	return nil, nil
}
func (f fakeBluetooth) WaitConnect(ctx context.Context, address string, timeoutSeconds int) error { return nil }
func (f fakeBluetooth) IsConnected(ctx context.Context, address string) (bool, error) { return false, nil }
func (f fakeBluetooth) ConnectedDevices(ctx context.Context) ([]core.Device, error) { return nil, nil }

func TestListConnectedFlagFiltersDevices(t *testing.T) {
	e := env{
		bluetooth: fakeBluetooth{devices: []core.Device{
			{Name: "A", Address: "AA", Connected: true},
			{Name: "B", Address: "BB", Connected: false},
			{Name: "C", Address: "CC", Connected: true},
		}},
		isTTY: func() bool { return false },
	}

	cmd := newListCmd(e)
	cmd.SetArgs([]string{"--format", "json", "--connected"})

	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error: %v", err)
	}

	got := out.String()
	// JSON is an array; ensure only connected devices appear.
	if bytes.Contains([]byte(got), []byte("\"name\": \"B\"")) {
		t.Fatalf("output should not include disconnected device B; got=%s", got)
	}
	if !bytes.Contains([]byte(got), []byte("\"name\": \"A\"")) || !bytes.Contains([]byte(got), []byte("\"name\": \"C\"")) {
		t.Fatalf("output should include connected devices A and C; got=%s", got)
	}
}

func TestListConnectedAndDisconnectedAreMutuallyExclusive(t *testing.T) {
	e := env{
		bluetooth: fakeBluetooth{devices: []core.Device{{Name: "A", Address: "AA", Connected: true}}},
		isTTY:     func() bool { return false },
	}

	cmd := newListCmd(e)
	cmd.SetArgs([]string{"--connected", "--disconnected"})
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})

	if err := cmd.Execute(); err == nil {
		t.Fatalf("expected error, got nil")
	}
}
