package core

import "context"

type BluetoothPort interface {
	List(ctx context.Context) ([]Device, error)
	Connect(ctx context.Context, address string) error
	Disconnect(ctx context.Context, address string) error

	// Pair pairs with a device.
	// This is required for "repair" (unpair -> inquiry -> pair -> connect).
	Pair(ctx context.Context, address string, pin string) error
	// Unpair removes pairing information for a device.
	Unpair(ctx context.Context, address string) error
	// Inquiry scans nearby devices (classic inquiry) for the specified duration.
	Inquiry(ctx context.Context, durationSeconds int) ([]Device, error)

	// WaitConnect waits until the device becomes connected (or times out).
	WaitConnect(ctx context.Context, address string, timeoutSeconds int) error
	// IsConnected returns whether the device is currently connected.
	IsConnected(ctx context.Context, address string) (bool, error)
	// ConnectedDevices lists currently connected devices.
	ConnectedDevices(ctx context.Context) ([]Device, error)
}

type PickerPort interface {
	PickDevice(ctx context.Context, title string, devices []Device) (Device, error)
	PickDevices(ctx context.Context, title string, devices []Device) ([]Device, error)

	// PickDeviceStream opens UI and updates device list as updates are received.
	PickDeviceStream(ctx context.Context, title string, updates <-chan []Device) (Device, error)
}
