package core

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"
)

type Pairer struct {
	Bluetooth      BluetoothPort
	Picker         PickerPort
	ProgressWriter io.Writer
}

type PairParams struct {
	Interactive     bool
	IsTTY           bool
	InquiryDuration int // seconds
	Pin             string
	WaitConnect     int // seconds
	MaxAttempts     int // default 3
}

func (p Pairer) progressf(format string, args ...any) {
	if p.ProgressWriter == nil {
		return
	}
	fmt.Fprintf(p.ProgressWriter, format, args...)
}

func (p Pairer) ensureInteractivePairing(params PairParams) error {
	if !params.Interactive {
		return fmt.Errorf("pair requires --interactive (TTY only)")
	}
	if !params.IsTTY {
		return fmt.Errorf("pair requires a TTY")
	}
	if p.Picker == nil {
		return fmt.Errorf("pair requires a picker")
	}
	return nil
}

func normalizeInquiryTotalSeconds(n int) int {
	if n <= 0 {
		return 60
	}
	return n
}

// pickByInquiryStream performs inquiry(loop) and opens a streaming picker once at least one device is discovered.
// If the scan window ends after UI starts, we keep the UI open and simply stop updating (no error).
func pickByInquiryStream(
	ctx context.Context,
	bluetooth BluetoothPort,
	picker PickerPort,
	progressf func(string, ...any),
	title string,
	totalSeconds int,
) (Device, error) {
	// Wrap progressf so that we can silence output once the TUI starts.
	uiActive := false
	pf := func(format string, args ...any) {
		if progressf == nil {
			return
		}
		if uiActive {
			// Bubble Tea uses the terminal; writing logs while it's running corrupts the UI.
			return
		}
		progressf(format, args...)
	}

	total := normalizeInquiryTotalSeconds(totalSeconds)
	chunk := 3
	deadline := time.Now().Add(time.Duration(total) * time.Second)

	pf("Searching nearby devices (up to %ds)...\n", total)

	updates := make(chan []Device, 8)
	uiStarted := make(chan struct{})

	scanCtx, scanCancel := context.WithDeadline(ctx, deadline)
	defer scanCancel()

	go func() {
		defer close(updates)
		seen := map[string]Device{}
		tick := 0
		for {
			if scanCtx.Err() != nil {
				return
			}
			tick++
			pf("  inquiry tick %d (chunk=%ds)\n", tick, chunk)

			found, err := bluetooth.Inquiry(scanCtx, chunk)
			if err != nil {
				pf("  inquiry error: %v\n", err)
				// If UI hasn't started yet, stop so caller can surface the error.
				select {
				case <-uiStarted:
					return
				default:
					return
				}
			}

			changed := false
			for _, d := range found {
				if strings.TrimSpace(d.Name) == "" || strings.TrimSpace(d.Address) == "" {
					continue
				}
				seen[d.Address] = d
				changed = true
			}

			if changed && len(seen) > 0 {
				snapshot := make([]Device, 0, len(seen))
				for _, d := range seen {
					snapshot = append(snapshot, d)
				}
				pf("  found %d device(s)\n", len(snapshot))
				select {
				case updates <- snapshot:
				default:
				}
			} else {
				pf("  found 0 device(s)\n")
			}
		}
	}()

	// Wait for first discovery until scan deadline.
	var first []Device
	select {
	case first = <-updates:
		if len(first) == 0 {
			return Device{}, ErrNotFound{Query: "no devices found"}
		}
	case <-scanCtx.Done():
		return Device{}, ErrNotFound{Query: "no devices found"}
	}

	close(uiStarted)
	uiActive = true

	uiUpdates := make(chan []Device, 16)
	uiUpdates <- first

	uiCtx, uiCancel := context.WithCancel(ctx)
	defer uiCancel()
	go func() {
		defer close(uiUpdates)
		for {
			select {
			case <-uiCtx.Done():
				return
			case ds, ok := <-updates:
				if !ok {
					return
				}
				select {
				case uiUpdates <- ds:
				default:
				}
			}
		}
	}()

	picked, err := picker.PickDeviceStream(ctx, title, uiUpdates)
	uiCancel()
	uiActive = false
	if err != nil {
		return Device{}, err
	}
	if strings.TrimSpace(picked.Address) == "" {
		return Device{}, fmt.Errorf("selected device has empty address")
	}
	return picked, nil
}

func connectWithRetryVerify(
	ctx context.Context,
	bluetooth BluetoothPort,
	progressf func(string, ...any),
	address string,
	waitConnectSeconds int,
	maxAttempts int,
) error {
	const waitConnectChunkSeconds = 5

	attempts := maxAttempts
	if attempts <= 0 {
		attempts = 3
	}
	var lastErr error
	remainingWaitSeconds := waitConnectSeconds

	for i := 1; i <= attempts; i++ {
		if progressf != nil {
			progressf("Connecting (attempt %d/%d)...\n", i, attempts)
		}
		if err := bluetooth.Connect(ctx, address); err != nil {
			lastErr = err
			if progressf != nil {
				progressf("  connect failed: %v\n", err)
			}
			continue
		}

		if remainingWaitSeconds > 0 {
			attemptWaitSeconds := remainingWaitSeconds
			if attemptWaitSeconds > waitConnectChunkSeconds {
				attemptWaitSeconds = waitConnectChunkSeconds
			}
			if progressf != nil {
				progressf(
					"  waiting for connection (up to %ds now; remaining budget %ds)...\n",
					attemptWaitSeconds,
					remainingWaitSeconds,
				)
			}
			start := time.Now()
			if err := bluetooth.WaitConnect(ctx, address, attemptWaitSeconds); err != nil {
				lastErr = err
				if progressf != nil {
					progressf("  wait-connect failed: %v\n", err)
					if ok, e := bluetooth.IsConnected(ctx, address); e == nil {
						progressf("  is-connected=%v\n", ok)
					}
					if cds, e := bluetooth.ConnectedDevices(ctx); e == nil {
						progressf("  connected devices: %d\n", len(cds))
					}
				}
				elapsed := int(time.Since(start).Seconds())
				remainingWaitSeconds -= elapsed
				if remainingWaitSeconds < 0 {
					remainingWaitSeconds = 0
				}
				continue
			}
			elapsed := int(time.Since(start).Seconds())
			remainingWaitSeconds -= elapsed
			if remainingWaitSeconds < 0 {
				remainingWaitSeconds = 0
			}
		}

		ok, err := bluetooth.IsConnected(ctx, address)
		if err == nil && ok {
			if progressf != nil {
				progressf("  connected confirmed\n")
			}
			return nil
		}
		if err != nil {
			lastErr = err
			if progressf != nil {
				progressf("  connect verification failed: %v\n", err)
			}
		} else {
			lastErr = fmt.Errorf("device is not connected")
			if progressf != nil {
				progressf("  connect verification failed: device is not connected\n")
			}
		}
	}

	if lastErr == nil {
		lastErr = fmt.Errorf("failed to connect")
	}
	return lastErr
}

func (p Pairer) pairPickedAndConnect(ctx context.Context, picked Device, params PairParams) (Device, error) {
	if strings.TrimSpace(picked.Address) == "" {
		return Device{}, fmt.Errorf("selected device has empty address")
	}
	if err := p.Bluetooth.Pair(ctx, picked.Address, params.Pin); err != nil {
		return Device{}, err
	}
	if err := connectWithRetryVerify(ctx, p.Bluetooth, p.progressf, picked.Address, params.WaitConnect, params.MaxAttempts); err != nil {
		return Device{}, err
	}
	return picked, nil
}

// Pair performs: inquiry(loop) -> pick discovered device (streaming) -> pair -> connect(wait/retry).
// It is intended for the situation where a device was already unpaired but connection is not yet established.
func (p Pairer) Pair(ctx context.Context, params PairParams) (Device, error) {
	if err := p.ensureInteractivePairing(params); err != nil {
		return Device{}, err
	}

	picked, err := pickByInquiryStream(ctx, p.Bluetooth, p.Picker, p.progressf, "Pair: select device", params.InquiryDuration)
	if err != nil {
		return Device{}, err
	}

	return p.pairPickedAndConnect(ctx, picked, params)
}
