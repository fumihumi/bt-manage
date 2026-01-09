package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/fumihumi/bt-manage/internal/core"
	"github.com/spf13/cobra"
)

func newPairCmd(e env) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pair",
		Short: "Pair and connect to a Bluetooth device",
		Long:  "Pair performs: inquiry -> pair -> connect. This is useful after unpairing when connection is not yet established.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			interactive, _ := cmd.Flags().GetBool("interactive")
			inquiry, _ := cmd.Flags().GetDuration("inquiry")
			pin, _ := cmd.Flags().GetString("pin")
			waitConnect, _ := cmd.Flags().GetDuration("wait-connect")
			maxAttempts, _ := cmd.Flags().GetInt("max-attempts")

			isTTY := e.isTTY()
			if interactive && !isTTY {
				return fmt.Errorf("--interactive requires a TTY")
			}

			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
			defer cancel()

			p := core.Pairer{Bluetooth: e.bluetooth, Picker: e.picker, ProgressWriter: cmd.ErrOrStderr()}
			dev, err := p.Pair(ctx, core.PairParams{
				Interactive:     interactive,
				IsTTY:           isTTY,
				InquiryDuration: int(inquiry.Truncate(time.Second).Seconds()),
				Pin:             pin,
				WaitConnect:     int(waitConnect.Truncate(time.Second).Seconds()),
				MaxAttempts:     maxAttempts,
			})
			if err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "paired: %s (%s)\n", dev.Name, dev.Address)
			return nil
		},
	}

	cmd.Flags().BoolP("interactive", "i", true, "Use interactive picker (TTY required)")
	cmd.Flags().Duration("inquiry", 60*time.Second, "Inquiry duration (e.g. 60s)")
	cmd.Flags().String("pin", "", "Optional PIN (if required by pairing)")
	cmd.Flags().Duration("wait-connect", 10*time.Second, "Total time budget to wait for the device to become connected across retries")
	cmd.Flags().Int("max-attempts", 6, "Connect retry count")

	return cmd
}
