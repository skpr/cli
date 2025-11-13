package alias

import (
	"fmt"
	"regexp"
	"slices"
)

// ExpandTemplate replaces $N in template with the corresponding args.
func ExpandTemplate(template string, args []string) string {
	replacer := func(match string) string {
		// Handle $1, $2, etc.
		re := regexp.MustCompile(`\$(\d+)`)
		m := re.FindStringSubmatch(match)
		if len(m) == 2 {
			idx := m[1]
			var i int
			fmt.Sscanf(idx, "%d", &i)
			if i > 0 && i <= len(args) {
				return args[i-1]
			}
		}
		return match
	}

	// Replace all $N or $@ matches
	re := regexp.MustCompile(`\$\d+`)
	return re.ReplaceAllStringFunc(template, replacer)
}

func CountTemplateArgs(template string) int {
	re := regexp.MustCompile(`\$\d+`)

	args := re.FindAllString(template, -1)
	slices.Sort(args)
	newArgs := slices.Compact(args)
	return len(newArgs)
}
