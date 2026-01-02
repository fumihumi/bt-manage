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
}

func defaultEnv() env {
	e := env{
		bluetooth: blueutil.Client{},
		isTTY:     tty.IsInteractive,
	}
	// picker は "TTYかつ--interactive" のときだけ使うが、生成は安価なのでここで固定。
	e.picker = picker.Picker{}
	return e
}

func newRootCmd() *cobra.Command {
	e := defaultEnv()

	cmd := &cobra.Command{
		Use:           "bt-manage",
		Short:         "Switch Bluetooth device connections on macOS",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			// `bt-manage` 単体実行は `list` と同義。
			return newListCmd(e).RunE(cmd, args)
		},
	}

	cmd.AddCommand(
		newListCmd(e),
		newConnectCmd(e),
		newDisconnectCmd(e),
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
