package core

import (
	"context"
	"errors"
	"testing"
)

type fakeBluetooth struct {
	devices []Device

	connected    []string
	disconnected []string

	paired   []string
	unpaired []string
	inquiry  []Device

	waitConnected []string

	listErr       error
	connectErr    error
	disconnectErr error
	pairErr       error
	unpairErr     error
	inquiryErr    error
	waitErr       error

	isConnected bool
	connectedList []Device
	isConnErr   error
	connListErr error
}

func (f *fakeBluetooth) List(ctx context.Context) ([]Device, error) {
	if f.listErr != nil {
		return nil, f.listErr
	}
	return append([]Device(nil), f.devices...), nil
}

func (f *fakeBluetooth) Connect(ctx context.Context, address string) error {
	if f.connectErr != nil {
		return f.connectErr
	}
	f.connected = append(f.connected, address)
	return nil
}

func (f *fakeBluetooth) Disconnect(ctx context.Context, address string) error {
	if f.disconnectErr != nil {
		return f.disconnectErr
	}
	f.disconnected = append(f.disconnected, address)
	return nil
}

func (f *fakeBluetooth) Pair(ctx context.Context, address string, pin string) error {
	if f.pairErr != nil {
		return f.pairErr
	}
	f.paired = append(f.paired, address)
	return nil
}

func (f *fakeBluetooth) Unpair(ctx context.Context, address string) error {
	if f.unpairErr != nil {
		return f.unpairErr
	}
	f.unpaired = append(f.unpaired, address)
	return nil
}

func (f *fakeBluetooth) Inquiry(ctx context.Context, durationSeconds int) ([]Device, error) {
	if f.inquiryErr != nil {
		return nil, f.inquiryErr
	}
	if f.inquiry != nil {
		return append([]Device(nil), f.inquiry...), nil
	}
	return nil, nil
}

func (f *fakeBluetooth) WaitConnect(ctx context.Context, address string, timeoutSeconds int) error {
	if f.waitErr != nil {
		return f.waitErr
	}
	f.waitConnected = append(f.waitConnected, address)
	return nil
}

func (f *fakeBluetooth) IsConnected(ctx context.Context, address string) (bool, error) {
	if f.isConnErr != nil {
		return false, f.isConnErr
	}
	return f.isConnected, nil
}

func (f *fakeBluetooth) ConnectedDevices(ctx context.Context) ([]Device, error) {
	if f.connListErr != nil {
		return nil, f.connListErr
	}
	return append([]Device(nil), f.connectedList...), nil
}

type fakePicker struct {
	picked Device
	err    error
	calls  int
}

func (p *fakePicker) PickDevice(ctx context.Context, title string, devices []Device) (Device, error) {
	p.calls++
	if p.err != nil {
		return Device{}, p.err
	}
	return p.picked, nil
}

func (p *fakePicker) PickDevices(ctx context.Context, title string, devices []Device) ([]Device, error) {
	p.calls++
	if p.err != nil {
		return nil, p.err
	}
	return []Device{p.picked}, nil
}

func (p *fakePicker) PickDeviceStream(ctx context.Context, title string, updates <-chan []Device) (Device, error) {
	p.calls++
	if p.err != nil {
		return Device{}, p.err
	}
	// Consume at least one update if present.
	select {
	case <-ctx.Done():
		return Device{}, ctx.Err()
	case <-updates:
	default:
	}
	return p.picked, nil
}

func TestConnectByNameOrInteractive(t *testing.T) {
	ctx := context.Background()
	devices := []Device{
		{Name: "MX Keys", Address: "AA"},
		{Name: "MX Master", Address: "BB"},
		{Name: "AirPods", Address: "CC"},
	}

	cases := []struct {
		name string
		p    ConnectParams

		wantErr       any
		wantConnected []string
		wantPickCalls int
	}{
		{
			name:    "prefix not found",
			p:       ConnectParams{Name: "ZZ"},
			wantErr: ErrNotFound{},
		},
		{
			name:          "prefix one match connects",
			p:             ConnectParams{Name: "Air"},
			wantConnected: []string{"CC"},
		},
		{
			name:    "prefix multiple non-tty is ambiguous",
			p:       ConnectParams{Name: "MX", IsTTY: false},
			wantErr: ErrAmbiguous{},
		},
		{
			name:          "prefix multiple tty uses picker",
			p:             ConnectParams{Name: "MX", IsTTY: true},
			wantConnected: []string{"BB"},
			wantPickCalls: 1,
		},
		{
			name:          "no name uses picker",
			p:             ConnectParams{Name: ""},
			wantConnected: []string{"BB"},
			wantPickCalls: 1,
		},
		{
			name:    "dry-run does not connect",
			p:       ConnectParams{Name: "Air", DryRun: true},
			wantConnected: nil,
		},
		{
			name:    "picker cancel propagates",
			p:       ConnectParams{Name: "MX", IsTTY: true},
			wantErr: ErrCanceled{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			bt := &fakeBluetooth{devices: devices}
			pk := &fakePicker{picked: Device{Name: "MX Master", Address: "BB"}}
			c := Connector{Bluetooth: bt, Picker: pk}

			if tc.name == "picker cancel propagates" {
				pk.err = ErrCanceled{}
			}

			_, err := c.ConnectByNameOrInteractive(ctx, tc.p)
			if tc.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error")
				}
				switch tc.wantErr.(type) {
				case ErrNotFound:
					var nf ErrNotFound
					if !errors.As(err, &nf) {
						t.Fatalf("expected ErrNotFound, got %T: %v", err, err)
					}
				case ErrAmbiguous:
					var am ErrAmbiguous
					if !errors.As(err, &am) {
						t.Fatalf("expected ErrAmbiguous, got %T: %v", err, err)
					}
				case ErrCanceled:
					var ce ErrCanceled
					if !errors.As(err, &ce) {
						t.Fatalf("expected ErrCanceled, got %T: %v", err, err)
					}
				default:
					t.Fatalf("unsupported wantErr type: %T", tc.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(bt.connected) != len(tc.wantConnected) {
				t.Fatalf("connected=%v, want %v", bt.connected, tc.wantConnected)
			}
			for i := range tc.wantConnected {
				if bt.connected[i] != tc.wantConnected[i] {
					t.Fatalf("connected=%v, want %v", bt.connected, tc.wantConnected)
				}
			}
			if pk.calls != tc.wantPickCalls {
				t.Fatalf("picker calls=%d, want %d", pk.calls, tc.wantPickCalls)
			}
		})
	}
}

func TestDisconnectByNameOrInteractive(t *testing.T) {
	ctx := context.Background()
	devices := []Device{
		{Name: "MX Keys", Address: "AA"},
		{Name: "AirPods", Address: "CC"},
	}

	bt := &fakeBluetooth{devices: devices}
	pk := &fakePicker{picked: Device{Name: "AirPods", Address: "CC"}}
	d := Disconnector{Bluetooth: bt, Picker: pk}

	_, err := d.DisconnectByNameOrInteractive(ctx, DisconnectParams{Name: "", IsTTY: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(bt.disconnected) != 1 || bt.disconnected[0] != "CC" {
		t.Fatalf("disconnected=%v, want [CC]", bt.disconnected)
	}
}
