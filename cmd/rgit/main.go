// Command rgit is a friendly companion for everyday git work.
package main

import (
	"errors"
	"os"

	"github.com/ettisafxrup/rgit/internal/cli"
	"github.com/ettisafxrup/rgit/internal/commands"
	"github.com/ettisafxrup/rgit/internal/git"
	"github.com/ettisafxrup/rgit/internal/ui"
)

// Exit codes.
const (
	exitOK    = 0
	exitError = 1
	exitUsage = 2
)

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	ui.PrepareConsole()
	ui.SetColor(ui.DetectColor())

	u := ui.New(os.Stdin, os.Stdout, os.Stderr)
	app := cli.NewApp(u, git.New(""), commands.All()...)

	err := app.Run(args)
	if err == nil {
		return exitOK
	}

	var usageErr *cli.UsageError
	switch {
	case errors.Is(err, cli.ErrAborted):
		u.Warn("Aborted. Nothing else was changed.")
	case errors.As(err, &usageErr):
		u.Error("%s", usageErr.Message)
		if usageErr.Command != nil {
			u.Hint("See 'rgit help %s' for usage.", usageErr.Command.Name)
		}
		return exitUsage
	default:
		u.Error("%s", err)
	}
	return exitError
}
