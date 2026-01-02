package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// These variables are intended to be set at build time via -ldflags.
var (
	version = "dev"
	commit  = ""
	date    = ""
)

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Minimal, stable format.
			if commit != "" || date != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "%s (%s %s)\n", version, commit, date)
				return nil
			}
			fmt.Fprintln(cmd.OutOrStdout(), version)
			return nil
		},
	}
}
