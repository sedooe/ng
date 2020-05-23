package main

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

const label = "plugins/ng"

func newRootCmd(io genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ng",
		Short: "namespace-group abstraction for Kubernetes.",
		//Long:                   globalUsage,
		SilenceUsage: true,
	}

	cf := genericclioptions.NewConfigFlags(true)

	// Add subcommands
	cmd.AddCommand(
		newAddCmd(io, cf),
		newRemoveCmd(io, cf),
		newShowCmd(io, cf),
	)

	cf.AddFlags(cmd.Flags())

	return cmd
}
