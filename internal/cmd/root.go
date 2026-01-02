package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "bt-manage",
		Short:         "Switch Bluetooth device connections on macOS",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			// `bt-manage` 単体実行は `list` と同義。
			return newListCmd().RunE(cmd, args)
		},
	}

	cmd.AddCommand(
		newListCmd(),
		newConnectCmd(),
		newDisconnectCmd(),
	)

	return cmd
}

// Execute is the entrypoint for the CLI.
func Execute() {
	root := newRootCmd()
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, prettyError(err))
		os.Exit(exitCodeFor(err))
	}
}
