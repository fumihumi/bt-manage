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

			format, err := output.ParseFormat(formatStr)
			if err != nil {
				return err
			}

			l := core.Lister{Bluetooth: e.bluetooth}
			devices, err := l.ListDevices(context.Background())
			if err != nil {
				return err
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

	cmd.Flags().String("format", "tsv", "Output format (tsv|json)")
	cmd.Flags().Bool("no-header", false, "Do not print header (tsv only)")

	return cmd
}
