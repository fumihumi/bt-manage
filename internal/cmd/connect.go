package cmd

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/fumihumi/bt-manage/internal/core"
	"github.com/fumihumi/bt-manage/internal/output"
	"github.com/spf13/cobra"
)

func newConnectCmd(e env) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "connect [<Name>]",
		Short: "Connect to a Bluetooth device",
		Long:  "Connect to a Bluetooth device. Output is a single device in the selected format (json is a 1-element array, consistent with 'list').",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Do NOT apply timeout to interactive (TUI) selection.
			baseCtx := context.Background()

			name := ""
			if len(args) == 1 {
				name = args[0]
			}

			exact, _ := cmd.Flags().GetBool("exact")
			interactiveFlagSet := cmd.Flags().Changed("interactive")
			interactive, _ := cmd.Flags().GetBool("interactive")
			multi, _ := cmd.Flags().GetBool("multi")
			dryRun, _ := cmd.Flags().GetBool("dry-run")
			formatStr, _ := cmd.Flags().GetString("format")
			noHeader, _ := cmd.Flags().GetBool("no-header")

			// Default behaviour: interactive picker is enabled by default when Name is omitted.
			// If Name is provided, keep the fast non-interactive behaviour unless user explicitly requested --interactive.
			if !interactiveFlagSet {
				if name == "" {
					interactive = true
				}
			}

			if multi {
				// Multi-select is only supported in interactive mode.
				interactive = true
				if name != "" {
					return fmt.Errorf("--multi cannot be used with a name argument")
				}
			}

			format, err := output.ParseFormat(formatStr)
			if err != nil {
				return err
			}

			isTTY := e.isTTY()
			if interactive && !isTTY {
				return fmt.Errorf("--interactive requires a TTY")
			}

			var pk core.PickerPort
			if interactive && isTTY {
				pk = e.picker
			}

			// Multi-select mode.
			if multi {
				if pk == nil {
					return fmt.Errorf("--multi requires a TTY")
				}

				// List/Pick can take time (user interaction); no timeout.
				devices, err := e.bluetooth.List(baseCtx)
				if err != nil {
					return err
				}

				selected, err := pk.PickDevices(baseCtx, "Connect", devices)
				if err != nil {
					return err
				}

				if dryRun {
					switch format {
					case output.FormatTSV:
						return output.WriteTSV(cmd.OutOrStdout(), selected, !noHeader)
					case output.FormatJSON:
						return output.WriteJSON(cmd.OutOrStdout(), selected)
					default:
						return fmt.Errorf("unsupported format")
					}
				}

				fmt.Fprintln(cmd.ErrOrStderr(), "Connecting...")

				type result struct {
					d   core.Device
					err error
				}
				results := make([]result, 0, len(selected))
				var mu sync.Mutex

				var wg sync.WaitGroup
				wg.Add(len(selected))
				for _, dev := range selected {
					dev := dev
					go func() {
						defer wg.Done()

						fmt.Fprintf(cmd.ErrOrStderr(), "- %s (%s)\n", dev.Name, dev.Address)

						// Per-device timeout (independent).
						dctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
						defer cancel()
						err := e.bluetooth.Connect(dctx, dev.Address)
						if err == nil {
							fmt.Fprintf(cmd.ErrOrStderr(), "  ok: connected %s (%s)\n", dev.Name, dev.Address)
						}

						mu.Lock()
						results = append(results, result{d: dev, err: err})
						mu.Unlock()
					}()
				}
				wg.Wait()

				var failed []string
				for _, r := range results {
					if r.err != nil {
						failed = append(failed, fmt.Sprintf("%s (%s): %v", r.d.Name, r.d.Address, r.err))
					}
				}
				if len(failed) > 0 {
					return errors.New("some connects failed: " + strings.Join(failed, "; "))
				}

				switch format {
				case output.FormatTSV:
					return output.WriteTSV(cmd.OutOrStdout(), selected, !noHeader)
				case output.FormatJSON:
					return output.WriteJSON(cmd.OutOrStdout(), selected)
				default:
					return fmt.Errorf("unsupported format")
				}
			}

			// Single-select.
			var ctx context.Context
			var cancel func()
			ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			c := core.Connector{Bluetooth: e.bluetooth, Picker: pk}
			selected, err := c.ConnectByNameOrInteractive(ctx, core.ConnectParams{
				Name:        name,
				Exact:       exact,
				Interactive: interactive,
				IsTTY:       isTTY,
				DryRun:      dryRun,
			})
			if err != nil {
				return err
			}

			if !dryRun {
				fmt.Fprintln(cmd.ErrOrStderr(), "Connecting...")
				fmt.Fprintf(cmd.ErrOrStderr(), "- %s (%s)\n", selected.Name, selected.Address)
			}

			switch format {
			case output.FormatTSV:
				return output.WriteTSV(cmd.OutOrStdout(), []core.Device{selected}, !noHeader)
			case output.FormatJSON:
				return output.WriteJSON(cmd.OutOrStdout(), []core.Device{selected})
			default:
				return fmt.Errorf("unsupported format")
			}
		},
	}

	cmd.Flags().BoolP("exact", "e", false, "Match device name exactly")
	cmd.Flags().BoolP("interactive", "i", false, "Always use interactive picker (TTY required)")
	cmd.Flags().BoolP("multi", "m", false, "Select multiple devices in the picker (implies --interactive; TTY only)")
	cmd.Flags().BoolP("dry-run", "n", false, "Do not connect; only resolve and print the target device")
	cmd.Flags().StringP("format", "f", "tsv", "Output format (tsv|json)")
	cmd.Flags().BoolP("no-header", "H", false, "Do not print header (tsv only)")

	return cmd
}
