package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/ettisafxrup/rgit/internal/git"
	"github.com/ettisafxrup/rgit/internal/ui"
	"github.com/ettisafxrup/rgit/internal/version"
)

// App dispatches command lines to commands.
type App struct {
	UI       *ui.UI
	Git      *git.Client
	commands []*Command
}

// NewApp creates an application with the given commands plus "help".
func NewApp(u *ui.UI, g *git.Client, commands ...*Command) *App {
	a := &App{UI: u, Git: g}
	a.commands = append(commands, a.helpCommand())
	return a
}

func (a *App) helpCommand() *Command {
	return &Command{
		Name:     "help",
		Group:    "Setup",
		Summary:  "Show help for rgit or one of its commands",
		Usage:    "[command]",
		Examples: []string{"help", "help push"},
		MaxArgs:  1,
		Run: func(ctx *Context) error {
			name := ctx.Arg(0)
			if name == "" {
				a.PrintHelp()
				return nil
			}
			cmd := a.Find(name)
			if cmd == nil {
				return fmt.Errorf("unknown command %q (see 'rgit help')", name)
			}
			a.PrintCommandHelp(cmd)
			return nil
		},
	}
}

// globalOptions are accepted by every command.
type globalOptions struct {
	yes     bool
	noColor bool
	help    bool
	version bool
}

func (o *globalOptions) register(fs *flag.FlagSet) {
	fs.BoolVar(&o.yes, "y", false, "answer yes to all questions and accept defaults")
	fs.BoolVar(&o.yes, "yes", false, "answer yes to all questions and accept defaults")
	fs.BoolVar(&o.noColor, "no-color", false, "disable colored output")
	fs.BoolVar(&o.help, "h", false, "show help")
	fs.BoolVar(&o.help, "help", false, "show help")
	fs.BoolVar(&o.version, "v", false, "print the version")
	fs.BoolVar(&o.version, "version", false, "print the version")
}

// Run executes the command line (without the program name).
func (a *App) Run(args []string) error {
	name, rest := splitCommand(args)
	if name == "" {
		return a.runWithoutCommand(rest)
	}

	cmd := a.Find(name)
	if cmd == nil {
		msg := fmt.Sprintf("unknown command %q", name)
		if s := a.suggest(name); s != "" {
			msg += fmt.Sprintf(" — did you mean %q?", s)
		}
		return errors.New(msg + " (see 'rgit help')")
	}

	var opts globalOptions
	fs := newFlagSet(cmd.Name)
	opts.register(fs)
	if cmd.Flags != nil {
		cmd.Flags(fs)
	}
	positional, err := parseInterspersed(fs, rest)
	if err != nil {
		return &UsageError{Command: cmd, Message: err.Error()}
	}
	a.apply(opts)

	switch {
	case opts.help:
		a.PrintCommandHelp(cmd)
		return nil
	case opts.version:
		a.printVersion()
		return nil
	case cmd.MaxArgs >= 0 && len(positional) > cmd.MaxArgs:
		return &UsageError{Command: cmd, Message: fmt.Sprintf("too many arguments for 'rgit %s'", cmd.Name)}
	}

	err = cmd.Run(&Context{UI: a.UI, Git: a.Git, Args: positional})
	var usageErr *UsageError
	if errors.As(err, &usageErr) && usageErr.Command == nil {
		usageErr.Command = cmd
	}
	return err
}

// runWithoutCommand handles "rgit", "rgit --version" and "rgit --help".
func (a *App) runWithoutCommand(args []string) error {
	var opts globalOptions
	fs := newFlagSet("rgit")
	opts.register(fs)
	if err := fs.Parse(args); err != nil {
		return err
	}
	a.apply(opts)
	if opts.version {
		a.printVersion()
		return nil
	}
	a.PrintHelp()
	return nil
}

func (a *App) apply(opts globalOptions) {
	if opts.noColor {
		ui.SetColor(false)
	}
	a.UI.AssumeYes = opts.yes
}

func (a *App) printVersion() {
	a.UI.Printf("rgit %s\n", version.Version)
}

// Commands returns every registered command.
func (a *App) Commands() []*Command { return a.commands }

// Find looks a command up by name or alias.
func (a *App) Find(name string) *Command {
	name = strings.ToLower(name)
	for _, c := range a.commands {
		if c.Name == name {
			return c
		}
		for _, alias := range c.Aliases {
			if alias == name {
				return c
			}
		}
	}
	return nil
}

// splitCommand separates leading global flags from the command name and
// its arguments: "rgit --no-color push -y" -> ("push", ["--no-color", "-y"]).
func splitCommand(args []string) (string, []string) {
	for i, arg := range args {
		if !strings.HasPrefix(arg, "-") {
			rest := append(append([]string{}, args[:i]...), args[i+1:]...)
			return arg, rest
		}
	}
	return "", args
}

func newFlagSet(name string) *flag.FlagSet {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(io.Discard) // errors are reported by rgit itself
	return fs
}

// parseInterspersed parses flags that may appear before, between or after
// positional arguments, so both "rgit push -y msg" and "rgit push msg -y" work.
// Everything after "--" is treated as positional.
func parseInterspersed(fs *flag.FlagSet, args []string) ([]string, error) {
	var tail []string
	for i, arg := range args {
		if arg == "--" {
			args, tail = args[:i], args[i+1:]
			break
		}
	}

	var positional []string
	for {
		if err := fs.Parse(args); err != nil {
			return nil, err
		}
		if fs.NArg() == 0 {
			return append(positional, tail...), nil
		}
		// flag stops at the first positional argument; keep it and go on.
		positional = append(positional, fs.Arg(0))
		args = fs.Args()[1:]
	}
}
