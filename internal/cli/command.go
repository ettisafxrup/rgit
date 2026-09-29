// Package cli contains the command model of rgit: how commands are
// described, how arguments are parsed and how help is shown.
package cli

import (
	"errors"
	"flag"

	"github.com/ettisafxrup/rgit/internal/git"
	"github.com/ettisafxrup/rgit/internal/ui"
)

// ErrAborted is returned when the user declines to continue.
var ErrAborted = errors.New("aborted")

// Command is a single rgit sub-command, such as "push" or "license".
type Command struct {
	Name     string
	Aliases  []string
	Group    string   // heading the command is listed under in help
	Summary  string   // one line shown in the command overview
	Usage    string   // argument synopsis, e.g. "[message]"
	Details  string   // longer explanation shown by "rgit help <command>"
	Examples []string // example invocations, without the leading "rgit "

	// Flags registers command-specific flags. It may be nil.
	Flags func(fs *flag.FlagSet)

	// MaxArgs limits the number of positional arguments; -1 means unlimited.
	MaxArgs int

	Run func(ctx *Context) error
}

// Context is everything a command needs while it runs.
type Context struct {
	UI   *ui.UI
	Git  *git.Client
	Args []string // positional arguments, flags removed
}

// Arg returns the i-th positional argument, or "" if it was not given.
func (c *Context) Arg(i int) string {
	if i < len(c.Args) {
		return c.Args[i]
	}
	return ""
}

// UsageError reports that a command was invoked incorrectly.
type UsageError struct {
	Command *Command
	Message string
}

func (e *UsageError) Error() string { return e.Message }
