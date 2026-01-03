package core

import (
	"context"
	"fmt"
	"io"
)

type Repairer struct {
	Bluetooth      BluetoothPort
	Picker         PickerPort
	ProgressWriter io.Writer // optional; if set, progress will be printed
}

type RepairParams struct {
	Interactive     bool
	IsTTY           bool
	InquiryDuration int // seconds (total window)
	Pin             string
	SkipUnpair      bool
	WaitConnect     int // seconds (0 disables)
	MaxAttempts     int // default 3
}

func (r Repairer) progressf(format string, args ...any) {
	if r.ProgressWriter == nil {
		return
	}
	fmt.Fprintf(r.ProgressWriter, format, args...)
}

// Repair performs: select paired device -> (optional) unpair -> inquiry(loop) -> pick discovered device (streaming) -> pair -> connect.
func (r Repairer) Repair(ctx context.Context, p RepairParams) (from Device, to Device, err error) {
	if !p.Interactive {
		return Device{}, Device{}, fmt.Errorf("repair requires --interactive (TTY only)")
	}
	if !p.IsTTY {
		return Device{}, Device{}, fmt.Errorf("repair requires a TTY")
	}
	if r.Picker == nil {
		return Device{}, Device{}, fmt.Errorf("repair requires a picker")
	}

	paired, err := r.Bluetooth.List(ctx)
	if err != nil {
		return Device{}, Device{}, err
	}
	if len(paired) == 0 {
		return Device{}, Device{}, ErrNotFound{Query: ""}
	}

	from, err = r.Picker.PickDevice(ctx, "Repair: select paired device to remove", paired)
	if err != nil {
		return Device{}, Device{}, err
	}

	if !p.SkipUnpair {
		if err := r.Bluetooth.Unpair(ctx, from.Address); err != nil {
			return from, Device{}, err
		}
	}

	picked, err := pickByInquiryStream(ctx, r.Bluetooth, r.Picker, r.progressf, "Repair: select device to pair", p.InquiryDuration)
	if err != nil {
		return from, Device{}, err
	}

	pairer := Pairer{Bluetooth: r.Bluetooth, Picker: r.Picker, ProgressWriter: r.ProgressWriter}
	to, err = pairer.pairPickedAndConnect(ctx, picked, PairParams{
		Interactive:     p.Interactive,
		IsTTY:           p.IsTTY,
		InquiryDuration: p.InquiryDuration,
		Pin:             p.Pin,
		WaitConnect:     p.WaitConnect,
		MaxAttempts:     p.MaxAttempts,
	})
	if err != nil {
		return from, to, err
	}

	return from, to, nil
}
