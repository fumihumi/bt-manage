package core

import (
	"context"
	"fmt"
)

type Connector struct {
	Bluetooth BluetoothPort
	Picker    PickerPort
}

type ConnectParams struct {
	Name        string
	Exact       bool
	Interactive bool
	IsTTY       bool
	DryRun      bool
}

func (c Connector) ConnectByNameOrInteractive(ctx context.Context, p ConnectParams) (Device, error) {
	devices, err := c.Bluetooth.List(ctx)
	if err != nil {
		return Device{}, err
	}

	if p.Name == "" {
		if c.Picker == nil {
			return Device{}, ErrNotFound{Query: ""}
		}
		selected, err := c.Picker.PickDevice(ctx, "Connect", devices)
		if err != nil {
			return Device{}, err
		}
		if p.DryRun {
			return selected, nil
		}
		if err := c.Bluetooth.Connect(ctx, selected.Address); err != nil {
			return Device{}, err
		}
		return selected, nil
	}

	matches := findByName(devices, p.Name, p.Exact)
	switch len(matches) {
	case 0:
		return Device{}, ErrNotFound{Query: p.Name}
	case 1:
		selected := matches[0]
		if p.Interactive {
			if !p.IsTTY {
				return Device{}, fmt.Errorf("interactive mode requires a TTY")
			}
			if c.Picker == nil {
				return Device{}, ErrNotFound{Query: p.Name}
			}
			selected2, err := c.Picker.PickDevice(ctx, "Connect", matches)
			if err != nil {
				return Device{}, err
			}
			selected = selected2
		}
		if p.DryRun {
			return selected, nil
		}
		if err := c.Bluetooth.Connect(ctx, selected.Address); err != nil {
			return Device{}, err
		}
		return selected, nil
	default:
		if !p.IsTTY {
			return Device{}, ErrAmbiguous{Query: p.Name, Count: len(matches)}
		}
		if c.Picker == nil {
			return Device{}, ErrAmbiguous{Query: p.Name, Count: len(matches)}
		}
		selected, err := c.Picker.PickDevice(ctx, "Connect", matches)
		if err != nil {
			return Device{}, err
		}
		if p.DryRun {
			return selected, nil
		}
		if err := c.Bluetooth.Connect(ctx, selected.Address); err != nil {
			return Device{}, err
		}
		return selected, nil
	}
}
