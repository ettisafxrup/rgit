package commands

import (
	"fmt"
	"strings"

	"github.com/ettisafxrup/rgit/internal/cli"
	"github.com/ettisafxrup/rgit/internal/ui"
)

func renameCommand() *cli.Command {
	return &cli.Command{
		Name:    "rename",
		Aliases: []string{"mv"},
		Group:   "Branches",
		Summary: "Rename a branch, locally and on the remote",
		Usage:   "<new> | <old>:<new> | <old> <new>",
		Details: `Renames the branch locally. If the old branch exists on the remote, rgit
offers to rename it there as well: the new branch is pushed first and the
old one is deleted only after that succeeded.`,
		Examples: []string{"rename feature/auth", "rename master:main", "rename old-name new-name"},
		MaxArgs:  2,
		Run:      runRename,
	}
}

func runRename(ctx *cli.Context) error {
	g, u := ctx.Git, ctx.UI
	if err := requireRepo(ctx); err != nil {
		return err
	}

	oldName, newName, err := parseRenameArgs(ctx)
	if err != nil {
		return err
	}
	if oldName == "" {
		if oldName, err = g.CurrentBranch(); err != nil {
			return err
		}
	}
	if !g.BranchExists(oldName) {
		return fmt.Errorf("there is no local branch named %q", oldName)
	}
	if err := validateBranchName(ctx, newName); err != nil {
		return err
	}
	if g.BranchExists(newName) {
		return fmt.Errorf("a branch named %q already exists", newName)
	}

	u.Info("Renaming %s → %s", ui.Bold(oldName), ui.Bold(newName))
	if _, err := g.Output("branch", "-m", oldName, newName); err != nil {
		return err
	}
	u.Success("Renamed the local branch")

	remote := g.DefaultRemote()
	if remote == "" {
		return nil
	}
	onRemote, err := g.RemoteBranchExists(remote, oldName)
	if err != nil {
		u.Warn("Could not check the remote: %v", err)
		return nil
	}
	if !onRemote {
		return nil
	}

	ok, err := u.Confirm(fmt.Sprintf("Rename %s on %s as well?", oldName, remote), true)
	if err != nil || !ok {
		return err
	}
	u.Step("Pushing %s", ui.Bold(newName))
	if err := g.Run("push", "--set-upstream", remote, newName); err != nil {
		return err
	}
	u.Step("Deleting %s from %s", ui.Bold(oldName), remote)
	if err := g.Run("push", remote, "--delete", oldName); err != nil {
		u.Hint("If %s is the default branch, change the default in your hosting settings first.", oldName)
		return fmt.Errorf("pushed %q but could not delete %q on the remote: %w", newName, oldName, err)
	}
	u.Success("Renamed %s → %s on %s", ui.Bold(oldName), ui.Bold(newName), remote)
	return nil
}

// parseRenameArgs accepts "<new>", "<old>:<new>" and "<old> <new>".
// An empty old name means the current branch.
func parseRenameArgs(ctx *cli.Context) (oldName, newName string, err error) {
	switch len(ctx.Args) {
	case 2:
		oldName, newName = ctx.Args[0], ctx.Args[1]
	case 1:
		if before, after, found := strings.Cut(ctx.Args[0], ":"); found {
			oldName, newName = before, after
			if oldName == "" {
				return "", "", &cli.UsageError{Message: "missing the old branch name before ':'"}
			}
		} else {
			newName = ctx.Args[0]
		}
	default:
		newName, err = ctx.UI.AskRequired("New branch name:", "")
	}
	if err == nil && strings.TrimSpace(newName) == "" {
		err = &cli.UsageError{Message: "missing the new branch name"}
	}
	return strings.TrimSpace(oldName), strings.TrimSpace(newName), err
}
