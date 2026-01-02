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
			formatStr, _ := cmd.Flags().GetString("format")
			noHeader, _ := cmd.Flags().GetBool("no-header")
			onlyConnected, _ := cmd.Flags().GetBool("connected")

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

	return cmd
}
