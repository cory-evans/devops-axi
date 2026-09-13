package main

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

const version = "0.1.0"

func main() {
	os.Exit(execute(os.Args[1:], os.Stdout))
}

func execute(args []string, out io.Writer) int {
	if len(args) == 1 && args[0] == "-V" {
		args[0] = "--version"
	}

	root := newRootCommand(out)
	root.SetArgs(args)
	if err := root.Execute(); err != nil {
		fmt.Fprintf(out, "error: %s\n", err)
		return 2
	}
	return 0
}

func newRootCommand(out io.Writer) *cobra.Command {
	root := &cobra.Command{
		Use:   "devops-axi",
		Short: "A focused Azure DevOps CLI for work items.",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	root.SetOut(out)
	root.SetErr(out)
	root.Version = version
	root.SetVersionTemplate("{{.Version}}\n")
	root.AddCommand(newWorkItemCommand(out))
	return root
}

func newWorkItemCommand(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "work-item",
		Short: "Manage Azure DevOps work items.",
	}
	cmd.SetOut(out)
	cmd.SetErr(out)
	cmd.AddCommand(newListCommand(out))
	return cmd
}
