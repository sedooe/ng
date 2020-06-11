package main

import (
	"fmt"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/sedooe/ng/internal/cmdutil"
	"github.com/sedooe/ng/internal/kubeclient"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

type showOptions struct {
	configFlags *genericclioptions.ConfigFlags
	io          genericclioptions.IOStreams
	cli         kubeclient.KubernetesClient

	name string
}

func newShowCmd(io genericclioptions.IOStreams, cf *genericclioptions.ConfigFlags) *cobra.Command {
	o := &showOptions{io: io, configFlags: cf}

	cmd := &cobra.Command{
		Use:          "show",
		Short:        "Show namespace group with its associated namespaces.",
		SilenceUsage: true,
		Args:         cmdutil.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			o.cli = kubeclient.NewKubeClient()
			if len(args) == 1 {
				o.name = args[0]
			}

			return o.run()
		},
	}

	o.configFlags.AddFlags(cmd.Flags())

	return cmd
}

func (o *showOptions) run() error {
	result, err := o.show()
	if err != nil {
		return errors.Wrap(err, "could not show namespace groups")
	}

	if len(result) == 0 {
		if o.name == "" {
			fmt.Fprintln(o.io.Out, "No namespace group found.")
		} else {
			fmt.Fprintf(o.io.Out, "Namespace group \"%s\" not found.\n", o.name)
		}
		return nil
	}

	table := tablewriter.NewWriter(o.io.Out)
	table.SetHeader([]string{"Namespace Group", "Namespaces"})

	for k, v := range result {
		row := []string{k, strings.Join(v, ", ")}
		table.Append(row)
	}
	table.Render()

	return nil
}

func (o *showOptions) show() (map[string][]string, error) {
	var selector string
	if o.name == "" {
		selector = label
	} else {
		selector = label + "=" + o.name
	}

	list, err := o.cli.GetNamespaces(selector)
	if err != nil {
		return nil, err
	}

	result := make(map[string][]string)

	for _, n := range list {
		ng := n.Labels[label]
		result[ng] = append(result[ng], n.Name)
	}

	return result, nil
}
