package core

import "context"

type BluetoothPort interface {
	List(ctx context.Context) ([]Device, error)
	Connect(ctx context.Context, address string) error
	Disconnect(ctx context.Context, address string) error
}

type PickerPort interface {
	PickDevice(ctx context.Context, title string, devices []Device) (Device, error)
	PickDevices(ctx context.Context, title string, devices []Device) ([]Device, error)
}
