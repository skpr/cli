package deploy

import (
	"fmt"

	"github.com/skpr/api/pb"

	"github.com/skpr/cli/internal/client"
)

// GetNames returns a list of environment names.
func GetNames(opts *client.Options) []string {
	client, ctx, err := opts.New()
	if err != nil {
		return []string{}
	}

	resp, err := client.Environment().List(ctx, &pb.EnvironmentListRequest{})
	if err != nil {
		fmt.Println("failed to list environments: %w", err)
	}

	var envs []string
	for _, env := range resp.Environments {
		envs = append(envs, env.Name)
	}

	return envs
}

// Contains checks if an environment already exists.
func Contains(name string, list []*pb.Environment) bool {
	for _, item := range list {
		if item.Name == name {
			return true
		}
	}

	return false
}
