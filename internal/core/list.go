package core

import (
	"context"
	"sort"
)

type Lister struct {
	Bluetooth BluetoothPort
}

func (l Lister) ListDevices(ctx context.Context) ([]Device, error) {
	devices, err := l.Bluetooth.List(ctx)
	if err != nil {
		return nil, err
	}

	sort.SliceStable(devices, func(i, j int) bool {
		a := devices[i].LastConnectedAt
		b := devices[j].LastConnectedAt
		if a == nil && b == nil {
			return devices[i].Name < devices[j].Name
		}
		if a == nil {
			return false
		}
		if b == nil {
			return true
		}
		return a.After(*b)
	})

	return devices, nil
}
