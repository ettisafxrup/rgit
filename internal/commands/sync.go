package commands

import (
	"errors"

	"github.com/ettisafxrup/rgit/internal/cli"
)

func syncCommand() *cli.Command {
	return &cli.Command{
		Name:    "sync",
		Group:   "Everyday",
		Summary: "Pull remote changes, then push your commits",
		Details: `Brings the current branch fully in line with the remote without making
a new commit: rebases onto the remote branch, then pushes local commits.
Uncommitted changes are stashed and restored automatically.`,
		Examples: []string{"sync"},
		MaxArgs:  0,
		Run: func(ctx *cli.Context) error {
			if err := requireRepo(ctx); err != nil {
				return err
			}
			if !ctx.Git.HasCommits() {
				return errors.New("this repository has no commits yet; use 'rgit push' to make the first one")
			}
			remote, err := ensureRemote(ctx)
			if err != nil {
				return err
			}
			branch, err := ctx.Git.CurrentBranch()
			if err != nil {
				return err
			}
			return syncWithRemote(ctx, remote, branch)
		},
	}
}
