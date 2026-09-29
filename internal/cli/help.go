package cli

import (
	"flag"
	"fmt"
	"strings"

	"github.com/ettisafxrup/rgit/internal/ui"
	"github.com/ettisafxrup/rgit/internal/version"
)

// Tagline describes rgit in a few words.
const Tagline = "git, without the ceremony"

// banner is the rgit wordmark drawn with box characters.
var banner = []string{
	"┬─┐┌─┐┬┌┬┐",
	"├┬┘│ ┬│ │ ",
	"┴└─└─┘┴ ┴ ",
}

// groupOrder is the order command groups appear in the help overview.
var groupOrder = []string{"Everyday", "Branches", "Project files", "Setup"}

// PrintHelp shows the overview of all commands.
func (a *App) PrintHelp() {
	u := a.UI
	u.Println()
	for i, line := range banner {
		switch i {
		case 0:
			u.Printf("  %s   %s\n", ui.Yellow(line), ui.Bold("rgit "+version.Version))
		case 1:
			u.Printf("  %s   %s\n", ui.Yellow(line), ui.Dim(Tagline))
		default:
			u.Printf("  %s\n", ui.Yellow(line))
		}
	}

	u.Printf("\n %s  rgit <command> [arguments] [flags]\n", ui.Bold("Usage:"))

	for _, group := range groupOrder {
		u.Printf("\n %s\n", ui.Bold(ui.Yellow(group)))
		for _, c := range a.commands {
			if c.Group == group {
				u.Printf("   %s %s\n", ui.Green(fmt.Sprintf("%-9s", c.Name)), c.Summary)
			}
		}
	}

	u.Printf("\n %s\n", ui.Bold(ui.Yellow("Global flags")))
	u.Printf("   %s answer yes to every question and accept defaults\n", ui.Green(fmt.Sprintf("%-15s", "-y, --yes")))
	u.Printf("   %s disable colored output (NO_COLOR is honored too)\n", ui.Green(fmt.Sprintf("%-15s", "--no-color")))
	u.Printf("   %s show help for a command\n", ui.Green(fmt.Sprintf("%-15s", "-h, --help")))
	u.Printf("   %s print the version\n", ui.Green(fmt.Sprintf("%-15s", "-v, --version")))

	u.Printf("\n Run %s for details and examples.\n\n", ui.Bold("rgit help <command>"))
}

// PrintCommandHelp shows the detailed help of one command.
func (a *App) PrintCommandHelp(c *Command) {
	u := a.UI
	u.Printf("\n %s — %s\n", ui.Bold(ui.Green("rgit "+c.Name)), c.Summary)

	usage := "rgit " + c.Name
	if c.Usage != "" {
		usage += " " + c.Usage
	}
	u.Printf("\n %s  %s [flags]\n", ui.Bold("Usage:"), usage)
	if len(c.Aliases) > 0 {
		u.Printf(" %s %s\n", ui.Bold("Aliases:"), strings.Join(c.Aliases, ", "))
	}

	if c.Details != "" {
		u.Printf("\n%s\n", ui.Indent(c.Details, 1))
	}

	if c.Flags != nil {
		fs := newFlagSet(c.Name)
		c.Flags(fs)
		u.Printf("\n %s\n", ui.Bold("Flags:"))
		fs.VisitAll(func(f *flag.Flag) {
			synopsis, usage := flagSynopsis(f)
			u.Printf("   %s %s\n", ui.Green(fmt.Sprintf("%-18s", synopsis)), usage)
		})
	}

	if len(c.Examples) > 0 {
		u.Printf("\n %s\n", ui.Bold("Examples:"))
		for _, ex := range c.Examples {
			u.Printf("   %s\n", ui.Dim("rgit "+ex))
		}
	}
	u.Println()
}

// flagSynopsis renders a flag as "-o <file>" plus its usage text. The
// placeholder comes from a `back-quoted` word in the usage, as in the flag package.
func flagSynopsis(f *flag.Flag) (string, string) {
	placeholder, usage := flag.UnquoteUsage(f)
	synopsis := "-" + f.Name
	if len(f.Name) > 1 {
		synopsis = "--" + f.Name
	}
	if placeholder != "" {
		synopsis += " <" + placeholder + ">"
	}
	return synopsis, usage
}

// suggest returns the command name closest to a mistyped one, or "".
func (a *App) suggest(input string) string {
	best, bestDistance := "", 3 // suggest only close matches
	for _, c := range a.commands {
		for _, name := range append([]string{c.Name}, c.Aliases...) {
			if d := distance(input, name); d < bestDistance {
				best, bestDistance = c.Name, d
			}
		}
	}
	return best
}

// distance computes the Levenshtein edit distance between two words.
func distance(a, b string) int {
	prev := make([]int, len(b)+1)
	for j := range prev {
		prev[j] = j
	}
	for i := 1; i <= len(a); i++ {
		cur := make([]int, len(b)+1)
		cur[0] = i
		for j := 1; j <= len(b); j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			cur[j] = min(prev[j]+1, cur[j-1]+1, prev[j-1]+cost)
		}
		prev = cur
	}
	return prev[len(b)]
}
