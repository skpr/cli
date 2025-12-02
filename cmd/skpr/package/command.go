package pkg

import (
	"github.com/spf13/cobra"

	skprcommand "github.com/skpr/cli/internal/command"
	v1package "github.com/skpr/cli/internal/command/package"
)

var (
	cmdLong = `Package a release that will be deployed to environments.`

	cmdExample = `
  # Package release 1.0.0 for deployment
  skpr package 1.0.0

  # Test the packaging process
  skpr package 1.0.0 --no-push

  # Package release 1.0.0 for deployment and print out a JSON manifest which
  # can be used for for other automation eg. Code scanning.
  skpr package 1.0.0 --print-manifest`
)

// NewCommand creates a new cobra.Command for 'package' sub command
func NewCommand() *cobra.Command {
	command := v1package.Command{}

	cmd := &cobra.Command{
		Use:                   "package <version>",
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"version"},
		Short:                 "Package a release for deployment",
		Long:                  cmdLong,
		Example:               cmdExample,
		GroupID:               skprcommand.GroupLifecycle,
		RunE: func(cmd *cobra.Command, args []string) error {
			command.Params.Version = args[0]
			return command.Run(cmd.Context())
		},
	}

	cmd.Flags().StringVar(&command.Region, "region", "ap-southeast-2", "Region which the AWS ECR registry resides.")
	cmd.Flags().StringVar(&command.Params.Context, "context", ".", "The context to use for the package command.")
	cmd.Flags().StringVar(&command.Params.IgnoreFile, "ignore-file", ".dockerignore", "A file containing patterns to exclude from the build context.")
	cmd.Flags().BoolVar(&command.Params.NoPush, "no-push", command.Params.NoPush, "Do not push the image to the registry.")
	cmd.Flags().BoolVar(&command.PrintManifest, "print-manifest", command.PrintManifest, "Print the manifest to stdout.")
	cmd.Flags().StringVar(&command.PackageDir, "dir", ".skpr/package", "The location of the package directory.")
	cmd.Flags().StringSliceVar(&command.BuildArgs, "build-arg", []string{}, "Additional build arguments.")
	cmd.Flags().BoolVar(&command.Debug, "debug", command.Debug, "Enable debug output.")

	return cmd
}
