package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/fang"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/spf13/cobra"

	"github.com/skpr/cli/cmd/skpr/alias"
	"github.com/skpr/cli/cmd/skpr/backup"
	"github.com/skpr/cli/cmd/skpr/config"
	"github.com/skpr/cli/cmd/skpr/create"
	"github.com/skpr/cli/cmd/skpr/cron"
	"github.com/skpr/cli/cmd/skpr/daemon"
	"github.com/skpr/cli/cmd/skpr/dashboard"
	deletecmd "github.com/skpr/cli/cmd/skpr/delete"
	"github.com/skpr/cli/cmd/skpr/deploy"
	execcmd "github.com/skpr/cli/cmd/skpr/exec"
	"github.com/skpr/cli/cmd/skpr/filesystem"
	"github.com/skpr/cli/cmd/skpr/info"
	"github.com/skpr/cli/cmd/skpr/list"
	"github.com/skpr/cli/cmd/skpr/login"
	"github.com/skpr/cli/cmd/skpr/logout"
	"github.com/skpr/cli/cmd/skpr/logs"
	"github.com/skpr/cli/cmd/skpr/mysql"
	pkg "github.com/skpr/cli/cmd/skpr/package"
	"github.com/skpr/cli/cmd/skpr/purge"
	"github.com/skpr/cli/cmd/skpr/release"
	"github.com/skpr/cli/cmd/skpr/restore"
	"github.com/skpr/cli/cmd/skpr/rsync"
	"github.com/skpr/cli/cmd/skpr/shell"
	"github.com/skpr/cli/cmd/skpr/trace"
	"github.com/skpr/cli/cmd/skpr/validate"
	"github.com/skpr/cli/cmd/skpr/version"
	aliastemplate "github.com/skpr/cli/internal/alias"
	"github.com/skpr/cli/internal/client/config/user"
	"github.com/skpr/cli/internal/color"
	skprcommand "github.com/skpr/cli/internal/command"
)

const cmdExample = `
    # Package application into container images.
    skpr package 0.0.1

    # Deploy packaged application to production.
    skpr deploy prod 0.0.1

    # Execute command on production environment.
    skpr exec prod -- echo "Skpr Rocks"

    # Configure secret for external "myapi" service.
    skpr config set prod myapi.key xxxyyyzzz
`

var cmd = &cobra.Command{
	Use:     "skpr",
	Short:   "Hugo is a very fast static site generator",
	Example: cmdExample,
	Long: `░██████╗██╗░░██╗██████╗░██████╗░  ░█████╗░██╗░░░░░██╗
██╔════╝██║░██╔╝██╔══██╗██╔══██╗  ██╔══██╗██║░░░░░██║
╚█████╗░█████═╝░██████╔╝██████╔╝  ██║░░╚═╝██║░░░░░██║
░╚═══██╗██╔═██╗░██╔═══╝░██╔══██╗  ██║░░██╗██║░░░░░██║
██████╔╝██║░╚██╗██║░░░░░██║░░██║  ╚█████╔╝███████╗██║
╚═════╝░╚═╝░░╚═╝╚═╝░░░░░╚═╝░░╚═╝  ░╚════╝░╚══════╝╚═╝
	
Develop with Skpr’s secure, dedicated hosting platform and discover 24/7 peace of mind.`,
}

func main() {
	// Load our configuration which contains aliases and feature flags.
	userConfig, err := user.NewClient()
	if err != nil {
		fmt.Println("Failed to load user config file:", err)
		os.Exit(1)
	}

	skprcommand.AddGroupsToCommand(cmd)

	cmd.AddCommand(alias.NewCommand())
	cmd.AddCommand(backup.NewCommand())
	cmd.AddCommand(config.NewCommand())
	cmd.AddCommand(create.NewCommand())
	cmd.AddCommand(cron.NewCommand())
	cmd.AddCommand(daemon.NewCommand())
	cmd.AddCommand(dashboard.NewCommand())
	cmd.AddCommand(deploy.NewCommand())
	cmd.AddCommand(deletecmd.NewCommand())
	cmd.AddCommand(execcmd.NewCommand())
	cmd.AddCommand(filesystem.NewCommand())
	cmd.AddCommand(info.NewCommand())
	cmd.AddCommand(list.NewCommand())
	cmd.AddCommand(login.NewCommand())
	cmd.AddCommand(logout.NewCommand())
	cmd.AddCommand(logs.NewCommand())
	cmd.AddCommand(mysql.NewCommand())
	cmd.AddCommand(pkg.NewCommand())
	cmd.AddCommand(purge.NewCommand())
	cmd.AddCommand(release.NewCommand())
	cmd.AddCommand(restore.NewCommand())
	cmd.AddCommand(rsync.NewCommand())
	cmd.AddCommand(shell.NewCommand())
	cmd.AddCommand(validate.NewCommand())
	cmd.AddCommand(version.NewCommand())

	// Hide the help command.
	cmd.SetHelpCommand(&cobra.Command{Hidden: true})

	// Hide the completions command.
	cmd.CompletionOptions = cobra.CompletionOptions{
		HiddenDefaultCmd: true,
	}

	// Experimental commands.
	featureFlags, err := userConfig.LoadFeatureFlags()
	if err != nil {
		fmt.Println("Failed to load feature flags:", err)
		os.Exit(1)
	}

	if featureFlags.Trace {
		cmd.AddCommand(trace.NewCommand())
	}

	// Alias commands.
	aliases, err := userConfig.ListAliases()
	if err != nil {
		fmt.Println("Failed to add alias commands:", err)
		os.Exit(1)
	}

	binPath, err := os.Executable()
	if err != nil {
		fmt.Println("Failed to get skpr executable path:", err)
		os.Exit(1)
	}

	for k, v := range aliases {
		cmd.AddCommand(&cobra.Command{
			Use:                   k,
			Args:                  cobra.ExactArgs(aliastemplate.CountTemplateArgs(v)),
			Short:                 fmt.Sprintf("Command: %s", v),
			DisableFlagsInUseLine: true,
			GroupID:               skprcommand.GroupAlias,
			RunE: func(cmd *cobra.Command, args []string) error {
				template := aliastemplate.ExpandTemplate(v, args)
				e := exec.Command(binPath, strings.Split(template, " ")...)
				e.Stdin = os.Stdin
				e.Stdout = os.Stdout
				e.Stderr = os.Stderr
				e.Env = os.Environ()
				return e.Run()
			},
		})
	}

	if err := fang.Execute(context.Background(), cmd, fang.WithColorSchemeFunc(MyColorScheme)); err != nil {
		os.Exit(1)
	}
}

// MyColorScheme customizes the default fang color scheme
func MyColorScheme(ld lipgloss.LightDarkFunc) fang.ColorScheme {
	// start from the defaults
	s := fang.DefaultColorScheme(ld)

	primary := ld(
		lipgloss.Color(color.HexOrange), // light mode
		lipgloss.Color(color.HexOrange), // dark mode
	)

	secondary := ld(
		lipgloss.Color(color.HexWhite), // light mode
		lipgloss.Color(color.HexWhite), // dark mode
	)

	s.Title = primary
	s.Command = secondary
	s.Flag = secondary

	s.Program = secondary

	return s
}
