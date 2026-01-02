package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/fumihumi/bt-manage/internal/core"
	"github.com/fumihumi/bt-manage/internal/output"
	"github.com/spf13/cobra"
)

func newDisconnectCmd(e env) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disconnect [<Name>]",
		Short: "Disconnect a Bluetooth device",
		Long:  "Disconnect a Bluetooth device. Output is a single device in the selected format (json is a 1-element array, consistent with 'list').",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			name := ""
			if len(args) == 1 {
				name = args[0]
			}

			exact, _ := cmd.Flags().GetBool("exact")
			interactive, _ := cmd.Flags().GetBool("interactive")
			multi, _ := cmd.Flags().GetBool("multi")
			dryRun, _ := cmd.Flags().GetBool("dry-run")
			formatStr, _ := cmd.Flags().GetString("format")
			noHeader, _ := cmd.Flags().GetBool("no-header")

			if multi {
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

			if multi {
				if pk == nil {
					return fmt.Errorf("--multi requires a TTY")
				}

				devices, err := e.bluetooth.List(ctx)
				if err != nil {
					return err
				}

				selected, err := pk.PickDevices(ctx, "Disconnect", devices)
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

				fmt.Fprintln(cmd.ErrOrStderr(), "Disconnecting...")
				for _, d := range selected {
					if ctx.Err() != nil {
						return fmt.Errorf("disconnect timed out")
					}
					fmt.Fprintf(cmd.ErrOrStderr(), "- %s (%s)\n", d.Name, d.Address)
					if err := e.bluetooth.Disconnect(ctx, d.Address); err != nil {
						return err
					}
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

			d := core.Disconnector{Bluetooth: e.bluetooth, Picker: pk}
			selected, err := d.DisconnectByNameOrInteractive(ctx, core.DisconnectParams{
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
				fmt.Fprintln(cmd.ErrOrStderr(), "Disconnecting...")
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
	cmd.Flags().BoolP("dry-run", "n", false, "Do not disconnect; only resolve and print the target device")
	cmd.Flags().StringP("format", "f", "tsv", "Output format (tsv|json)")
	cmd.Flags().BoolP("no-header", "H", false, "Do not print header (tsv only)")

	return cmd
}
