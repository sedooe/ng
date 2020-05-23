package main

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sedooe/ng/internal/cmdutil"
	"github.com/sedooe/ng/internal/kubeclient"
	"github.com/spf13/cobra"
	"k8s.io/api/core/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

type addOptions struct {
	configFlags *genericclioptions.ConfigFlags
	io          genericclioptions.IOStreams
	cli         kubeclient.KubernetesClient

	name string
}

func newAddCmd(io genericclioptions.IOStreams, cf *genericclioptions.ConfigFlags) *cobra.Command {
	o := &addOptions{io: io, configFlags: cf}

	cmd := &cobra.Command{
		Use:          "add",
		Short:        "Adds namespace-group to a namespace.",
		SilenceUsage: true,
		Args:         cmdutil.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			o.name = args[0]
			o.cli = kubeclient.NewKubeClient()
			return o.run()
		},
	}

	o.configFlags.AddFlags(cmd.Flags())

	return cmd
}

func (o *addOptions) run() error {
	ns, err := o.cli.GetOperationNamespace(o.configFlags)
	if err != nil {
		return err
	}

	err = o.addNg(ns)
	if err != nil {
		return err
	}

	fmt.Fprintf(o.io.Out, "Namespace \"%s\" successfully added to the group \"%s\".\n", ns.Name, o.name)

	return nil
}

func (o *addOptions) addNg(ns *v1.Namespace) error {
	labels := ns.Labels
	if labels == nil {
		labels = map[string]string{label: o.name}
	} else {
		if labels[label] != "" {
			return fmt.Errorf("this namespace already belongs to a group: %s", labels[label])
		}
		labels[label] = o.name
	}

	ns.SetLabels(labels)

	err := o.cli.UpdateNamespace(ns)
	if err != nil {
		return errors.Wrap(err, "namespace update failed")
	}

	return nil
}
