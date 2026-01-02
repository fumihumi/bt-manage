package core

import (
	"context"
	"fmt"
)

type Disconnector struct {
	Bluetooth BluetoothPort
	Picker    PickerPort
}

type DisconnectParams struct {
	Name        string
	Exact       bool
	Interactive bool
	IsTTY       bool
	DryRun      bool
}

func (d Disconnector) DisconnectByNameOrInteractive(ctx context.Context, p DisconnectParams) (Device, error) {
	devices, err := d.Bluetooth.List(ctx)
	if err != nil {
		return Device{}, err
	}

	if p.Name == "" {
		if d.Picker == nil {
			return Device{}, ErrNotFound{Query: ""}
		}
		selected, err := d.Picker.PickDevice(ctx, "Disconnect", devices)
		if err != nil {
			return Device{}, err
		}
		if p.DryRun {
			return selected, nil
		}
		if err := d.Bluetooth.Disconnect(ctx, selected.Address); err != nil {
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
			if d.Picker == nil {
				return Device{}, ErrNotFound{Query: p.Name}
			}
			selected2, err := d.Picker.PickDevice(ctx, "Disconnect", matches)
			if err != nil {
				return Device{}, err
			}
			selected = selected2
		}
		if p.DryRun {
			return selected, nil
		}
		if err := d.Bluetooth.Disconnect(ctx, selected.Address); err != nil {
			return Device{}, err
		}
		return selected, nil
	default:
		if !p.IsTTY {
			return Device{}, ErrAmbiguous{Query: p.Name, Count: len(matches)}
		}
		if d.Picker == nil {
			return Device{}, ErrAmbiguous{Query: p.Name, Count: len(matches)}
		}
		selected, err := d.Picker.PickDevice(ctx, "Disconnect", matches)
		if err != nil {
			return Device{}, err
		}
		if p.DryRun {
			return selected, nil
		}
		if err := d.Bluetooth.Disconnect(ctx, selected.Address); err != nil {
			return Device{}, err
		}
		return selected, nil
	}
}
