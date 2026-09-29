// Package commands implements every rgit sub-command.
package commands

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/ettisafxrup/rgit/internal/cli"
	"github.com/ettisafxrup/rgit/internal/git"
	"github.com/ettisafxrup/rgit/internal/ui"
)

// All returns every rgit command in the order they appear in help.
func All() []*cli.Command {
	return []*cli.Command{
		// Everyday
		pushCommand(), pullCommand(), syncCommand(), statusCommand(), logCommand(), undoCommand(),
		// Branches
		newCommand(), switchCommand(), renameCommand(), cleanupCommand(),
		// Project files
		ignoreCommand(), licenseCommand(), readmeCommand(), archiCommand(),
		// Setup
		cloneCommand(), loginCommand(), whoamiCommand(), openCommand(), doctorCommand(), versionCommand(),
	}
}

// defaultBranch is the branch name rgit uses for brand-new repositories.
const defaultBranch = "main"

var errNotRepo = errors.New("this folder is not a git repository (run 'rgit push' to create one, or 'rgit clone' to download one)")

// requireRepo fails unless git is installed and we are inside a repository.
func requireRepo(ctx *cli.Context) error {
	if err := git.Installed(); err != nil {
		return err
	}
	if !ctx.Git.IsRepo() {
		return errNotRepo
	}
	return nil
}

// ensureRepo makes sure we are inside a repository, offering to create one.
func ensureRepo(ctx *cli.Context) error {
	if err := git.Installed(); err != nil {
		return err
	}
	if ctx.Git.IsRepo() {
		return nil
	}
	ctx.UI.Warn("This folder is not a git repository yet.")
	ok, err := ctx.UI.Confirm("Initialize a new repository here?", true)
	if err != nil {
		return err
	}
	if !ok {
		return cli.ErrAborted
	}
	if err := ctx.Git.Init(defaultBranch); err != nil {
		return err
	}
	ctx.UI.Success("Initialized an empty repository on branch %s", ui.Bold(defaultBranch))
	return nil
}

// ensureRemote returns the remote to work with, asking the user to add one
// when the repository has none.
func ensureRemote(ctx *cli.Context) (string, error) {
	if remote := ctx.Git.DefaultRemote(); remote != "" {
		return remote, nil
	}
	ctx.UI.Warn("This repository has no remote yet.")
	name, err := ctx.UI.Ask("Remote name:", "origin")
	if err != nil {
		return "", err
	}
	rawURL, err := ctx.UI.AskRequired("Repository URL (or GitHub owner/repo):", "")
	if err != nil {
		return "", err
	}
	url := git.ExpandCloneURL(rawURL)
	if err := ctx.Git.AddRemote(name, url); err != nil {
		return "", err
	}
	ctx.UI.Success("Added remote %s → %s", ui.Bold(name), url)
	return name, nil
}

// alignUnbornBranch renames the branch of a repository without commits to
// match the remote's default branch, so the first push/pull lines up with
// what already exists on the server (main vs. master).
func alignUnbornBranch(ctx *cli.Context, remote string) error {
	if ctx.Git.HasCommits() {
		return nil
	}
	remoteDefault, err := ctx.Git.RemoteDefaultBranch(remote)
	if err != nil || remoteDefault == "" {
		return err
	}
	current, _ := ctx.Git.CurrentBranch()
	if current == remoteDefault {
		return nil
	}
	if _, err := ctx.Git.Output("symbolic-ref", "HEAD", "refs/heads/"+remoteDefault); err != nil {
		return err
	}
	ctx.UI.Info("Using branch %s to match the remote.", ui.Bold(remoteDefault))
	return nil
}

// syncWithRemote rebases the current branch onto its remote counterpart and
// pushes it, setting the upstream on the first push.
func syncWithRemote(ctx *cli.Context, remote, branch string) error {
	g, u := ctx.Git, ctx.UI

	if upstream, ok := g.Upstream(); ok {
		u.Step("Pulling latest changes from %s", ui.Bold(upstream))
		if err := g.Run("pull", "--rebase", "--autostash"); err != nil {
			return pullFailed(ctx, err)
		}
		u.Step("Pushing %s", ui.Bold(branch))
		if err := g.Run("push"); err != nil {
			return err
		}
		u.Success("Pushed %s to %s", ui.Bold(branch), ui.Bold(upstream))
		return nil
	}

	exists, err := g.RemoteBranchExists(remote, branch)
	if err != nil {
		return err
	}
	if exists {
		u.Step("Pulling latest changes from %s", ui.Bold(remote+"/"+branch))
		if err := g.Run("pull", "--rebase", "--autostash", remote, branch); err != nil {
			return pullFailed(ctx, err)
		}
	}
	u.Step("Pushing %s to %s for the first time", ui.Bold(branch), ui.Bold(remote))
	if err := g.Run("push", "--set-upstream", remote, branch); err != nil {
		return err
	}
	u.Success("Pushed %s to %s", ui.Bold(branch), ui.Bold(remote+"/"+branch))
	return nil
}

// pullFailed explains how to recover from a failed rebase.
func pullFailed(ctx *cli.Context, err error) error {
	if ctx.Git.Succeeds("rev-parse", "--verify", "--quiet", "REBASE_HEAD") {
		ctx.UI.Warn("The rebase stopped because of conflicts.")
		ctx.UI.Hint("Fix the conflicted files, then run:  git add <file>  and  git rebase --continue")
		ctx.UI.Hint("Or cancel everything with:           git rebase --abort")
		return errors.New("pull stopped with conflicts")
	}
	return fmt.Errorf("pull failed: %w", err)
}

// today returns the date in a friendly format, e.g. "29 September 2026".
func today() string { return time.Now().Format("2 January 2006") }

// confirmOverwrite asks before replacing an existing file.
func confirmOverwrite(ctx *cli.Context, path string) error {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return nil
	}
	ok, err := ctx.UI.Confirm(fmt.Sprintf("%s already exists. Overwrite it?", path), false)
	if err != nil {
		return err
	}
	if !ok {
		return cli.ErrAborted
	}
	return nil
}

// validateBranchName checks a branch name with git's own rules.
func validateBranchName(ctx *cli.Context, name string) error {
	if !ctx.Git.Succeeds("check-ref-format", "--branch", name) {
		return fmt.Errorf("%q is not a valid branch name", name)
	}
	return nil
}
