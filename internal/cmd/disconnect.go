package cmd

import (
	"fmt"

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
			fmt.Fprintf(cmd.OutOrStdout(), "TODO: disconnect %s\n", name)
			return nil
		},
	}

	cmd.Flags().Bool("exact", false, "Match device name exactly")
	cmd.Flags().Bool("interactive", false, "Always use interactive picker (TTY required)")
	cmd.Flags().Bool("dry-run", false, "Print what would be executed without disconnecting")

	return cmd
}
