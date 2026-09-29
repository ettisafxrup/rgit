package commands

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ettisafxrup/rgit/internal/cli"
	"github.com/ettisafxrup/rgit/internal/git"
	"github.com/ettisafxrup/rgit/internal/ui"
)

func pushCommand() *cli.Command {
	return &cli.Command{
		Name:    "push",
		Aliases: []string{"p"},
		Group:   "Everyday",
		Summary: "Stage, commit and push all changes in one step",
		Usage:   "[message]",
		Details: `Does everything needed to get your work onto the remote:

  1. initializes a repository and adds a remote, if they are missing
  2. shows your changes and stages them all (after asking)
  3. commits them with your message (quotes are optional)
  4. rebases onto the latest remote changes, then pushes

Without a message you are asked for one; pressing Enter uses "Update <date>".`,
		Examples: []string{`push "Fix the login redirect"`, "push Add dark mode", "push -y"},
		MaxArgs:  -1,
		Run:      runPush,
	}
}

func runPush(ctx *cli.Context) error {
	g := ctx.Git

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
	branch, err := g.CurrentBranch()
	if err != nil {
		return err
	}

	if err := stageChanges(ctx); err != nil {
		return err
	}
	if err := commitStaged(ctx, strings.Join(ctx.Args, " ")); err != nil {
		return err
	}
	if !g.HasCommits() {
		return errors.New("there is nothing to push yet: the repository has no commits")
	}
	return syncWithRemote(ctx, remote, branch)
}

// stageChanges shows the working tree changes and stages them all.
func stageChanges(ctx *cli.Context) error {
	g, u := ctx.Git, ctx.UI
	status, err := g.Status()
	if err != nil {
		return err
	}
	if status.Conflicts > 0 {
		return errors.New("some files have merge conflicts; resolve them before pushing")
	}
	if status.Clean() {
		return nil
	}

	u.Info("%d changed %s:", len(status.Changes), plural(len(status.Changes), "file", "files"))
	printChanges(u, status.Changes, 15)

	ok, err := u.Confirm("Stage all of these changes?", true)
	if err != nil {
		return err
	}
	if !ok {
		if g.HasStagedChanges() {
			u.Info("Keeping only what is already staged.")
			return nil
		}
		return cli.ErrAborted
	}
	_, err = g.Output("add", "--all")
	return err
}

// commitStaged commits the staged changes, asking for a message if needed.
func commitStaged(ctx *cli.Context, message string) error {
	g, u := ctx.Git, ctx.UI
	if !g.HasStagedChanges() {
		u.Info("Nothing new to commit.")
		return nil
	}

	message = strings.TrimSpace(message)
	if message == "" {
		var err error
		if message, err = u.Ask("Commit message:", "Update "+today()); err != nil {
			return err
		}
	}
	if _, err := g.Output("commit", "--quiet", "-m", message); err != nil {
		return commitError(err)
	}
	last, err := g.LastCommit()
	if err != nil {
		return err
	}
	u.Success("Committed %s %s", ui.Yellow(last.Hash), last.Subject)
	return nil
}

// commitError turns git's "please tell me who you are" into clear advice.
func commitError(err error) error {
	var gitErr *git.Error
	if errors.As(err, &gitErr) && strings.Contains(gitErr.Stderr, "tell me who you are") {
		return errors.New("git does not know who you are yet; run 'rgit login' to set your name and email")
	}
	return err
}

// printChanges lists up to limit changes with a colored label.
func printChanges(u *ui.UI, changes []git.Change, limit int) {
	for i, c := range changes {
		if i == limit {
			u.Printf("   %s\n", ui.Dim(fmt.Sprintf("… and %d more", len(changes)-limit)))
			return
		}
		label := fmt.Sprintf("%-10s", c.Label())
		switch c.Label() {
		case "new", "added":
			label = ui.Green(label)
		case "deleted", "conflict":
			label = ui.Red(label)
		default:
			label = ui.Yellow(label)
		}
		u.Printf("   %s %s\n", label, c.Path)
	}
}

func plural(n int, one, many string) string {
	if n == 1 {
		return one
	}
	return many
}
