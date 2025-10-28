package command

import (
	"github.com/spf13/cobra"
)

const (
	GroupAlias            = "alias"
	GroupDisasterRecovery = "dr"
	GroupSecureShell      = "secure-shell"
	GroupBackground       = "background"
	GroupDataStorage      = "data-storage"
	GroupAuthentication   = "authentication"
	GroupDebug            = "debug"
	GroupCDN              = "cdn"
	GroupLifecycle        = "lifecycle"
)

// AddGroupsToCommand adds command groups to the provided cobra.Command
func AddGroupsToCommand(cmd *cobra.Command) {
	cmd.AddGroup(&cobra.Group{
		ID:    GroupAuthentication,
		Title: "Authentication Commands",
	})

	cmd.AddGroup(&cobra.Group{
		ID:    GroupLifecycle,
		Title: "Lifecycle Commands",
	})

	cmd.AddGroup(&cobra.Group{
		ID:    GroupSecureShell,
		Title: "Secure Shell Commands",
	})

	cmd.AddGroup(&cobra.Group{
		ID:    GroupDisasterRecovery,
		Title: "Disaster Recovery Commands",
	})

	cmd.AddGroup(&cobra.Group{
		ID:    GroupDataStorage,
		Title: "Data Storage Commands",
	})

	cmd.AddGroup(&cobra.Group{
		ID:    GroupBackground,
		Title: "Background Task Commands",
	})

	cmd.AddGroup(&cobra.Group{
		ID:    GroupDebug,
		Title: "Debug Commands",
	})

	cmd.AddGroup(&cobra.Group{
		ID:    GroupCDN,
		Title: "CDN Commands",
	})

	cmd.AddGroup(&cobra.Group{
		ID:    GroupAlias,
		Title: "Alias Commands",
	})
}
