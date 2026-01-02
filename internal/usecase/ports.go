package usecase

import (
	"context"

	"github.com/fumihumi/bt-manage/internal/domain"
)

type BluetoothPort interface {
	List(ctx context.Context) ([]domain.Device, error)
	Connect(ctx context.Context, address string) error
	Disconnect(ctx context.Context, address string) error
}

type PickerPort interface {
	PickDevice(ctx context.Context, title string, devices []domain.Device) (domain.Device, error)
	PickDevices(ctx context.Context, title string, devices []domain.Device) ([]domain.Device, error)
}
