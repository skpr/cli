package aliases

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/skpr/cli/internal/client/config/command"
)

// Expand expands aliases in the list of args.
func Expand(args []string, aliases command.Aliases) (bool, []string, error) {
	found := false
	if len(args) == 0 {
		return found, args, nil
	}
	expansion, found := aliases[args[0]]
	if !found {
		// No alias found.
		return found, args, nil
	}

	var extraArgs []string
	for i, a := range args[1:] {
		if !strings.Contains(expansion, "$") {
			extraArgs = append(extraArgs, a)
		} else {
			expansion = strings.ReplaceAll(expansion, fmt.Sprintf("$%d", i+1), a)
		}
	}
	lingeringRE := regexp.MustCompile(`\$\d`)
	if lingeringRE.MatchString(expansion) {
		err := fmt.Errorf("not enough arguments for alias: %s", expansion)
		return found, []string{}, err
	}

	var newArgs = strings.Split(expansion, " ")

	expanded := append(newArgs, extraArgs...)
	return found, expanded, nil
}
