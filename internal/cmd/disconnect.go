package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/fumihumi/bt-manage/internal/core"
	"github.com/fumihumi/bt-manage/internal/platform/macos/blueutil"
	"github.com/fumihumi/bt-manage/internal/platform/ui/nopicker"
	"github.com/spf13/cobra"
)

func newDisconnectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disconnect [<Name>]",
		Short: "Disconnect a Bluetooth device",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := ""
			if len(args) == 1 {
				name = args[0]
			}

			exact, _ := cmd.Flags().GetBool("exact")
			interactive, _ := cmd.Flags().GetBool("interactive")
			dryRun, _ := cmd.Flags().GetBool("dry-run")

			d := core.Disconnector{Bluetooth: blueutil.Client{}, Picker: nopicker.Picker{}}
			selected, err := d.DisconnectByNameOrInteractive(context.Background(), core.DisconnectParams{
				Name:        name,
				Exact:       exact,
				Interactive: interactive,
				IsTTY:       cmd.OutOrStdout() == os.Stdout,
				DryRun:      dryRun,
			})
			if err != nil {
				return err
			}

			if dryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "DRY-RUN: disconnect %s (%s)\n", selected.Name, selected.Address)
				return nil
			}
			fmt.Fprintf(cmd.OutOrStdout(), "disconnected: %s (%s)\n", selected.Name, selected.Address)
			return nil
		},
	}

	cmd.Flags().Bool("exact", false, "Match device name exactly")
	cmd.Flags().Bool("interactive", false, "Always use interactive picker (TTY required)")
	cmd.Flags().Bool("dry-run", false, "Print what would be executed without disconnecting")

	return cmd
}
