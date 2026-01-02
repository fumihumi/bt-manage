package nopicker

import (
	"context"

	"github.com/fumihumi/bt-manage/internal/core"
)

type Picker struct{}

func (p Picker) PickDevice(ctx context.Context, title string, devices []core.Device) (core.Device, error) {
	return core.Device{}, core.ErrCanceled{}
}
