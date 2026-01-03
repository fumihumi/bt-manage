package cmd

import (
	"context"
	"fmt"

	"github.com/fumihumi/bt-manage/internal/core"
	"github.com/fumihumi/bt-manage/internal/output"
	"github.com/spf13/cobra"
)

func newListCmd(e env) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List Bluetooth devices",
		RunE: func(cmd *cobra.Command, args []string) error {
			namesOnly, _ := cmd.Flags().GetBool("names-only")
			formatStr, _ := cmd.Flags().GetString("format")
			noHeader, _ := cmd.Flags().GetBool("no-header")
			onlyConnected, _ := cmd.Flags().GetBool("connected")
			onlyDisconnected, _ := cmd.Flags().GetBool("disconnected")
			// Currently `list` always lists paired devices. `--paired` is a compatibility/explicitness flag.
			_, _ = cmd.Flags().GetBool("paired")

			if onlyConnected && onlyDisconnected {
				return fmt.Errorf("--connected and --disconnected are mutually exclusive")
			}

			if namesOnly {
				// Keep behaviour simple and predictable.
				if cmd.Flags().Changed("format") {
					return fmt.Errorf("--names-only cannot be used with --format")
				}
				if cmd.Flags().Changed("no-header") {
					return fmt.Errorf("--names-only cannot be used with --no-header")
				}
			}

			format, err := output.ParseFormat(formatStr)
			if err != nil {
				return err
			}

			l := core.Lister{Bluetooth: e.bluetooth}
			devices, err := l.ListDevices(context.Background())
			if err != nil {
				return err
			}

			if onlyConnected {
				filtered := make([]core.Device, 0, len(devices))
				for _, d := range devices {
					if d.Connected {
						filtered = append(filtered, d)
					}
				}
				devices = filtered
			}

			if onlyDisconnected {
				filtered := make([]core.Device, 0, len(devices))
				for _, d := range devices {
					if !d.Connected {
						filtered = append(filtered, d)
					}
				}
				devices = filtered
			}

			if namesOnly {
				for _, d := range devices {
					if d.Name == "" {
						fmt.Fprintln(cmd.OutOrStdout(), "(unknown)")
						continue
					}
					fmt.Fprintln(cmd.OutOrStdout(), d.Name)
				}
				return nil
			}

			switch format {
			case output.FormatTSV:
				return output.WriteTSV(cmd.OutOrStdout(), devices, !noHeader)
			case output.FormatJSON:
				return output.WriteJSON(cmd.OutOrStdout(), devices)
			default:
				return fmt.Errorf("unsupported format")
			}
		},
	}

	cmd.Flags().StringP("format", "f", "tsv", "Output format (tsv|json)")
	cmd.Flags().BoolP("no-header", "H", false, "Do not print header (tsv only)")
	cmd.Flags().BoolP("connected", "c", false, "Show connected devices only")
	cmd.Flags().BoolP("disconnected", "d", false, "Show disconnected devices only")
	cmd.Flags().BoolP("names-only", "N", false, "Print device names only (one per line)")
	cmd.Flags().Bool("paired", true, "List paired devices (default)")

	return cmd
}
