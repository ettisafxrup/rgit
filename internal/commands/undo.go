package commands

import (
	"errors"
	"fmt"

	"github.com/ettisafxrup/rgit/internal/cli"
	"github.com/ettisafxrup/rgit/internal/ui"
)

func undoCommand() *cli.Command {
	return &cli.Command{
		Name:    "undo",
		Group:   "Everyday",
		Summary: "Undo the last commit but keep its changes",
		Details: `Removes the most recent commit while keeping every change it contained
staged, so you can fix the message, add a forgotten file or split it up.
rgit warns you first if that commit has already been pushed.`,
		Examples: []string{"undo"},
		MaxArgs:  0,
		Run:      runUndo,
	}
}

func runUndo(ctx *cli.Context) error {
	g, u := ctx.Git, ctx.UI
	if err := requireRepo(ctx); err != nil {
		return err
	}
	if !g.HasCommits() {
		return errors.New("there is no commit to undo")
	}
	if g.Succeeds("rev-parse", "--verify", "--quiet", "HEAD^2") {
		return errors.New("the last commit is a merge; undo it with 'git reset --merge HEAD~1' instead")
	}

	last, err := g.LastCommit()
	if err != nil {
		return err
	}
	u.Info("Last commit: %s %s %s", ui.Yellow(last.Hash), last.Subject, ui.Dim("("+last.When+")"))

	pushed, _ := g.Lines("branch", "--remotes", "--contains", "HEAD")
	defaultAnswer := true
	if len(pushed) > 0 {
		u.Warn("This commit is already on the remote. Undoing it means you will need to force-push later.")
		defaultAnswer = false
	}
	ok, err := u.Confirm("Undo this commit? Its changes will stay staged.", defaultAnswer)
	if err != nil {
		return err
	}
	if !ok {
		return cli.ErrAborted
	}

	if g.Succeeds("rev-parse", "--verify", "--quiet", "HEAD~1") {
		_, err = g.Output("reset", "--soft", "HEAD~1")
	} else {
		// The very first commit has no parent: drop the branch ref instead.
		branch, branchErr := g.CurrentBranch()
		if branchErr != nil {
			return branchErr
		}
		_, err = g.Output("update-ref", "-d", "refs/heads/"+branch)
	}
	if err != nil {
		return fmt.Errorf("could not undo the commit: %w", err)
	}
	u.Success("Undid %s — your changes are still staged.", ui.Yellow(last.Hash))
	return nil
}
