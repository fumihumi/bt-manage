package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List Bluetooth devices",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintln(cmd.OutOrStdout(), "TODO: list")
			return nil
		},
	}

	cmd.Flags().String("format", "tsv", "Output format (tsv|json)")
	cmd.Flags().Bool("no-header", false, "Do not print header (tsv only)")

	return cmd
}
