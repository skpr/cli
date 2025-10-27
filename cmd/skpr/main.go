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
	deletecmd "github.com/skpr/cli/cmd/skpr/delete"
	"github.com/skpr/cli/cmd/skpr/deploy"
	execcmd "github.com/skpr/cli/cmd/skpr/exec"
	"github.com/skpr/cli/cmd/skpr/info"
	"github.com/skpr/cli/cmd/skpr/list"
	"github.com/skpr/cli/cmd/skpr/login"
	"github.com/skpr/cli/cmd/skpr/logout"
	"github.com/skpr/cli/cmd/skpr/mysql"
	pkg "github.com/skpr/cli/cmd/skpr/package"
	"github.com/skpr/cli/cmd/skpr/project"
	"github.com/skpr/cli/cmd/skpr/purge"
	"github.com/skpr/cli/cmd/skpr/release"
	"github.com/skpr/cli/cmd/skpr/restore"
	"github.com/skpr/cli/cmd/skpr/rsync"
	"github.com/skpr/cli/cmd/skpr/shell"
	"github.com/skpr/cli/cmd/skpr/validate"
	"github.com/skpr/cli/cmd/skpr/version"
	"github.com/skpr/cli/internal/client/config/user"
	"github.com/skpr/cli/internal/color"
)

const (
	// GroupAliases is the ID for the alias command group.
	GroupAliases = "aliases"
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
	cmd.AddCommand(alias.NewCommand())
	cmd.AddCommand(backup.NewCommand())
	cmd.AddCommand(config.NewCommand())
	cmd.AddCommand(create.NewCommand())
	cmd.AddCommand(cron.NewCommand())
	cmd.AddCommand(daemon.NewCommand())
	cmd.AddCommand(deletecmd.NewCommand())
	cmd.AddCommand(deploy.NewCommand())
	cmd.AddCommand(execcmd.NewCommand())
	cmd.AddCommand(info.NewCommand())
	cmd.AddCommand(list.NewCommand())
	cmd.AddCommand(login.NewCommand())
	cmd.AddCommand(logout.NewCommand())
	cmd.AddCommand(mysql.NewCommand())
	cmd.AddCommand(pkg.NewCommand())
	cmd.AddCommand(project.NewCommand())
	cmd.AddCommand(purge.NewCommand())
	cmd.AddCommand(restore.NewCommand())
	cmd.AddCommand(rsync.NewCommand())
	cmd.AddCommand(shell.NewCommand())
	cmd.AddCommand(version.NewCommand())
	cmd.AddCommand(release.NewCommand())
	cmd.AddCommand(validate.NewCommand())

	// Add user set aliases to the root command.
	err := addAliases(cmd)
	if err != nil {
		fmt.Println("Failed to add alias commands:", err)
		os.Exit(1)
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

// Adds user defined aliases to the root command.
func addAliases(cmd *cobra.Command) error {
	configFile, err := user.NewConfigFile()

	// Only add aliases if we got an error.
	if err == nil {
		// Load the aliases.
		aliases, err := configFile.ReadAliases()
		if err != nil {
			return fmt.Errorf("failed to read aliases: %w", err)
		}

		binPath, err := os.Executable()
		if err != nil {
			return fmt.Errorf("failed to get skpr executable path: %w", err)
		}

		cmd.AddGroup(&cobra.Group{
			ID:    GroupAliases,
			Title: "Alias Commands",
		})

		for k, v := range aliases {
			cmd.AddCommand(&cobra.Command{
				Use:                   k,
				Short:                 fmt.Sprintf("Command: %s", v),
				DisableFlagsInUseLine: true,
				GroupID:               GroupAliases,
				RunE: func(cmd *cobra.Command, args []string) error {
					e := exec.Command(binPath, strings.Split(v, " ")...)
					e.Stdin = os.Stdin
					e.Stdout = os.Stdout
					e.Stderr = os.Stderr
					e.Env = os.Environ()
					return e.Run()
				},
			})
		}
	}

	return nil
}
