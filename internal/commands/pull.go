package commands

import (
	"fmt"

	"github.com/ettisafxrup/rgit/internal/cli"
	"github.com/ettisafxrup/rgit/internal/ui"
)

func pullCommand() *cli.Command {
	return &cli.Command{
		Name:    "pull",
		Group:   "Everyday",
		Summary: "Get the latest remote changes (rebase, keeps local work safe)",
		Usage:   "[branch]",
		Details: `Pulls with --rebase so your history stays linear, and --autostash so
uncommitted work is set aside and restored automatically.

In a folder that is not a repository yet, rgit offers to initialize it and
connect it to a remote, then downloads the remote's default branch.`,
		Examples: []string{"pull", "pull develop"},
		MaxArgs:  1,
		Run:      runPull,
	}
}

func runPull(ctx *cli.Context) error {
	g, u := ctx.Git, ctx.UI

	if err := ensureRepo(ctx); err != nil {
		return err
	}
	remote, err := ensureRemote(ctx)
	if err != nil {
		return err
	}
	if err := alignUnbornBranch(ctx, remote); err != nil {
		return err
	}

	branch := ctx.Arg(0)
	if upstream, ok := g.Upstream(); ok && branch == "" {
		u.Step("Pulling latest changes from %s", ui.Bold(upstream))
		if err := g.Run("pull", "--rebase", "--autostash"); err != nil {
			return pullFailed(ctx, err)
		}
		u.Success("Up to date with %s", ui.Bold(upstream))
		return nil
	}

	current, err := g.CurrentBranch()
	if err != nil {
		return err
	}
	if branch == "" {
		branch = current
	}
	exists, err := g.RemoteBranchExists(remote, branch)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the remote %q has no branch named %q yet", remote, branch)
	}

	u.Step("Pulling latest changes from %s", ui.Bold(remote+"/"+branch))
	if err := g.Run("pull", "--rebase", "--autostash", remote, branch); err != nil {
		return pullFailed(ctx, err)
	}
	if branch == current {
		// Remember where this branch lives so plain "git pull" works too.
		g.Output("branch", "--set-upstream-to="+remote+"/"+branch)
	}
	u.Success("Up to date with %s", ui.Bold(remote+"/"+branch))
	return nil
}
