package main

import (
	"errors"
	"fmt"

	"github.com/sedooe/ng/internal/cmdutil"
	"github.com/sedooe/ng/internal/kubeclient"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

type removeOptions struct {
	configFlags *genericclioptions.ConfigFlags
	io          genericclioptions.IOStreams
	cli         kubeclient.KubernetesClient
}

func newRemoveCmd(io genericclioptions.IOStreams, cf *genericclioptions.ConfigFlags) *cobra.Command {
	o := &removeOptions{io: io, configFlags: cf}

	cmd := &cobra.Command{
		Use:          "remove",
		Short:        "Removes namespace from namespace-group with configured RBAC rules.",
		SilenceUsage: true,
		Args:         cmdutil.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			o.cli = kubeclient.NewKubeClient()
			return o.run()
		},
	}

	o.configFlags.AddFlags(cmd.Flags())

	return cmd
}

func (o *removeOptions) run() error {
	ns, err := o.cli.GetOperationNamespace(o.configFlags)
	if err != nil {
		return err
	}

	removed, err := o.removeNg(ns)
	if err != nil {
		return err
	}

	fmt.Fprintf(o.io.Out, "Namespace \"%s\" successfully removed from the group \"%s\".\n", ns.Name, removed)

	return nil
}

func (o *removeOptions) removeNg(ns *v1.Namespace) (string, error) {
	labels := ns.Labels
	if labels == nil || labels[label] == "" {
		return "", errors.New("this namespace does not belong to any group.")
	}

	removed := labels[label]

	delete(labels, label)

	ns.SetLabels(labels)

	err := o.cli.UpdateNamespace(ns)
	if err != nil {
		return "", err
	}

	return removed, nil
}
