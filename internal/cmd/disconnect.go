package cmd

import (
	"context"
	"fmt"

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
			name := ""
			if len(args) == 1 {
				name = args[0]
			}

			exact, _ := cmd.Flags().GetBool("exact")
			interactive, _ := cmd.Flags().GetBool("interactive")
			dryRun, _ := cmd.Flags().GetBool("dry-run")
			formatStr, _ := cmd.Flags().GetString("format")
			noHeader, _ := cmd.Flags().GetBool("no-header")

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

			d := core.Disconnector{Bluetooth: e.bluetooth, Picker: pk}
			selected, err := d.DisconnectByNameOrInteractive(context.Background(), core.DisconnectParams{
				Name:        name,
				Exact:       exact,
				Interactive: interactive,
				IsTTY:       isTTY,
				DryRun:      dryRun,
			})
			if err != nil {
				return err
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

	cmd.Flags().Bool("exact", false, "Match device name exactly")
	cmd.Flags().Bool("interactive", false, "Always use interactive picker (TTY required)")
	cmd.Flags().Bool("dry-run", false, "Do not disconnect; only resolve and print the target device")
	cmd.Flags().String("format", "tsv", "Output format (tsv|json)")
	cmd.Flags().Bool("no-header", false, "Do not print header (tsv only)")

	return cmd
}
