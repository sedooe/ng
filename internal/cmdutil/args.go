package cmdutil

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func ExactArgs(n int) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) != n {
			return errors.Errorf(
				"%q requires %d argument",
				cmd.CommandPath(),
				n,
			)
		}
		return nil
	}
}
