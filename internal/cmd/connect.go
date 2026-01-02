package cmd

import (
	"context"
	"fmt"

	"github.com/fumihumi/bt-manage/internal/core"
	"github.com/fumihumi/bt-manage/internal/output"
	"github.com/fumihumi/bt-manage/internal/platform/macos/blueutil"
	"github.com/fumihumi/bt-manage/internal/platform/tty"
	"github.com/fumihumi/bt-manage/internal/tui/picker"
	"github.com/spf13/cobra"
)

func newConnectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "connect [<Name>]",
		Short: "Connect to a Bluetooth device",
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

			format, err := output.ParseFormat(formatStr)
			if err != nil {
				return err
			}

			isTTY := tty.IsInteractive()
			if interactive && !isTTY {
				return fmt.Errorf("--interactive requires a TTY")
			}

			var pk core.PickerPort
			if interactive && isTTY {
				pk = picker.Picker{}
			}

			c := core.Connector{Bluetooth: blueutil.Client{}, Picker: pk}
			selected, err := c.ConnectByNameOrInteractive(context.Background(), core.ConnectParams{
				Name:        name,
				Exact:       exact,
				Interactive: interactive,
				IsTTY:       isTTY,
				DryRun:      dryRun,
			})
			if err != nil {
				return err
			}

			// For structured formats, print the selected device.
			switch format {
			case output.FormatTSV:
				return output.WriteTSV(cmd.OutOrStdout(), []core.Device{selected}, true)
			case output.FormatJSON:
				return output.WriteJSON(cmd.OutOrStdout(), []core.Device{selected})
			default:
				return fmt.Errorf("unsupported format")
			}
		},
	}

	cmd.Flags().Bool("exact", false, "Match device name exactly")
	cmd.Flags().Bool("interactive", false, "Always use interactive picker (TTY required)")
	cmd.Flags().Bool("dry-run", false, "Print what would be executed without connecting")
	cmd.Flags().String("format", "tsv", "Output format (tsv|json)")

	return cmd
}
