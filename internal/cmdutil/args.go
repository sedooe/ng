package cmdutil

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// thanks: https://github.com/helm/helm/blob/master/cmd/helm/require/args.go

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

func MaximumNArgs(n int) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) > n {
			return errors.Errorf(
				"%q accepts at most %d argument",
				cmd.CommandPath(),
				n,
			)
		}
		return nil
	}
}
