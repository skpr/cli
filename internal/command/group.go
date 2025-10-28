package command

import (
	"github.com/spf13/cobra"
)

const (
	// GroupAlias for grouping commands related to alias management
	GroupAlias = "alias"
	// GroupDisasterRecovery for grouping our backup/restore commands
	GroupDisasterRecovery = "dr"
	// GroupSecureShell for grouping our ssh related commands eg. exec, shell, rsync
	GroupSecureShell = "secure-shell"
	// GroupBackground for grouping our background task commands eg. daemons, cron
	GroupBackground = "background"
	// GroupDataStorage for grouping our data storage commands eg. mysql, filesystem
	GroupDataStorage = "data-storage"
	// GroupAuthentication for grouping our authentication commands eg. login, logout
	GroupAuthentication = "authentication"
	// GroupDebug for grouping our debugging commands eg. trace, logs
	GroupDebug = "debug"
	// GroupCDN for grouping our CDN commands eg. purge
	GroupCDN = "cdn"
	// GroupLifecycle for grouping our lifecycle commands eg. create, delete, deploy, config
	GroupLifecycle = "lifecycle"
	// GroupOther for grouping all other commands.
	GroupOther = "other"
)

// AddGroupsToCommand adds command groups to the provided cobra.Command
func AddGroupsToCommand(cmd *cobra.Command) {
	cmd.AddGroup(&cobra.Group{
		ID:    GroupAuthentication,
		Title: "Authentication",
	})

	cmd.AddGroup(&cobra.Group{
		ID:    GroupLifecycle,
		Title: "Deployment Lifecycle",
	})

	cmd.AddGroup(&cobra.Group{
		ID:    GroupSecureShell,
		Title: "Secure Shell",
	})

	cmd.AddGroup(&cobra.Group{
		ID:    GroupDisasterRecovery,
		Title: "Disaster Recovery",
	})

	cmd.AddGroup(&cobra.Group{
		ID:    GroupDataStorage,
		Title: "Data Storage",
	})

	cmd.AddGroup(&cobra.Group{
		ID:    GroupBackground,
		Title: "Background Tasks",
	})

	cmd.AddGroup(&cobra.Group{
		ID:    GroupDebug,
		Title: "Debug",
	})

	cmd.AddGroup(&cobra.Group{
		ID:    GroupCDN,
		Title: "CDN",
	})

	cmd.AddGroup(&cobra.Group{
		ID:    GroupAlias,
		Title: "Alias",
	})

	cmd.AddGroup(&cobra.Group{
		ID:    GroupOther,
		Title: "Other",
	})
}
