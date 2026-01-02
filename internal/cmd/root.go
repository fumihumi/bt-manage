package cmd

import (
	"fmt"
	"os"

	"github.com/fumihumi/bt-manage/internal/core"
	"github.com/fumihumi/bt-manage/internal/platform/macos/blueutil"
	"github.com/fumihumi/bt-manage/internal/platform/tty"
	"github.com/fumihumi/bt-manage/internal/tui/picker"
	"github.com/spf13/cobra"
)

type env struct {
	bluetooth core.BluetoothPort
	picker    core.PickerPort
	isTTY     func() bool
	verbose   bool
}

func defaultEnv(verbose bool) env {
	client := blueutil.Client{Verbose: verbose, Logger: os.Stderr}
	return env{
		bluetooth: client,
		picker:    picker.Picker{},
		isTTY:     tty.IsInteractive,
		verbose:   verbose,
	}
}

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "bt-manage",
		Short:         "Switch Bluetooth device connections on macOS",
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose logging to stderr")

	cmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		// no-op: env is built per-command in RunE below
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		verbose, _ := cmd.Flags().GetBool("verbose")
		e := defaultEnv(verbose)
		// `bt-manage` 単体実行は `list` と同義。
		return newListCmd(e).RunE(cmd, args)
	}

	cmd.AddCommand(
		newListCmd(defaultEnv(false)),
		newConnectCmd(defaultEnv(false)),
		newDisconnectCmd(defaultEnv(false)),
		newVersionCmd(),
	)

	// 子コマンド実行時の env を verbose に追従させるため、各コマンドの PreRun で env を差し替える。
	for _, c := range cmd.Commands() {
		origRunE := c.RunE
		if origRunE == nil {
			continue
		}
		c.RunE = func(cmd2 *cobra.Command, args2 []string) error {
			verbose, _ := cmd2.Flags().GetBool("verbose")
			e := defaultEnv(verbose)
			// コマンド生成時の env を反映するため、ここでは再生成して実行する。
			switch cmd2.Name() {
			case "list":
				return newListCmd(e).RunE(cmd2, args2)
			case "connect":
				return newConnectCmd(e).RunE(cmd2, args2)
			case "disconnect":
				return newDisconnectCmd(e).RunE(cmd2, args2)
			default:
				return origRunE(cmd2, args2)
			}
		}
	}

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
