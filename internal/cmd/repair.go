package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/fumihumi/bt-manage/internal/core"
	"github.com/spf13/cobra"
)

func newRepairCmd(e env) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "repair",
		Short: "Unpair and re-pair a Bluetooth device",
		Long:  "Repair performs: select a paired device -> unpair -> inquiry -> pair -> connect. This is useful when a device is visible but cannot connect.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			interactive, _ := cmd.Flags().GetBool("interactive")
			inquiry, _ := cmd.Flags().GetDuration("inquiry")
			pin, _ := cmd.Flags().GetString("pin")
			skipUnpair, _ := cmd.Flags().GetBool("skip-unpair")
			waitConnect, _ := cmd.Flags().GetDuration("wait-connect")
			maxAttempts, _ := cmd.Flags().GetInt("max-attempts")

			isTTY := e.isTTY()
			if interactive && !isTTY {
				return fmt.Errorf("--interactive requires a TTY")
			}

			// Inquiry/pair/connect may take time.
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
			defer cancel()

			r := core.Repairer{Bluetooth: e.bluetooth, Picker: e.picker, ProgressWriter: cmd.ErrOrStderr()}
			from, to, err := r.Repair(ctx, core.RepairParams{
				Interactive:     interactive,
				IsTTY:           isTTY,
				InquiryDuration: int(inquiry.Truncate(time.Second).Seconds()),
				Pin:             pin,
				SkipUnpair:      skipUnpair,
				WaitConnect:     int(waitConnect.Truncate(time.Second).Seconds()),
				MaxAttempts:     maxAttempts,
			})
			if err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "repaired: %s (%s) -> %s (%s)\n", from.Name, from.Address, to.Name, to.Address)
			return nil
		},
	}

	cmd.Flags().BoolP("interactive", "i", true, "Use interactive picker (TTY required)")
	cmd.Flags().Duration("inquiry", 60*time.Second, "Inquiry duration (e.g. 60s)")
	cmd.Flags().String("pin", "", "Optional PIN (if required by pairing)")
	cmd.Flags().Bool("skip-unpair", false, "Skip unpair step")
	cmd.Flags().Duration("wait-connect", 10*time.Second, "Total time budget to wait for the device to become connected across retries")
	cmd.Flags().Int("max-attempts", 6, "Connect retry count")

	return cmd
}
