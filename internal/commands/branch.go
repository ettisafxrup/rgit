package commands

import (
	"errors"
	"fmt"

	"github.com/ettisafxrup/rgit/internal/cli"
	"github.com/ettisafxrup/rgit/internal/ui"
)

func newCommand() *cli.Command {
	return &cli.Command{
		Name:     "new",
		Aliases:  []string{"branch"},
		Group:    "Branches",
		Summary:  "Create a new branch and switch to it",
		Usage:    "<branch>",
		Details:  "Your uncommitted changes come along to the new branch.",
		Examples: []string{"new feature/login"},
		MaxArgs:  1,
		Run: func(ctx *cli.Context) error {
			if err := requireRepo(ctx); err != nil {
				return err
			}
			name := ctx.Arg(0)
			if name == "" {
				var err error
				if name, err = ctx.UI.AskRequired("New branch name:", ""); err != nil {
					return err
				}
			}
			if err := validateBranchName(ctx, name); err != nil {
				return err
			}
			if ctx.Git.BranchExists(name) {
				return fmt.Errorf("branch %q already exists; use 'rgit switch %s'", name, name)
			}
			if _, err := ctx.Git.Output("checkout", "-b", name); err != nil {
				return err
			}
			ctx.UI.Success("Created and switched to %s", ui.Bold(name))
			return nil
		},
	}
}

func switchCommand() *cli.Command {
	return &cli.Command{
		Name:     "switch",
		Aliases:  []string{"sw", "checkout"},
		Group:    "Branches",
		Summary:  "Switch branches (pick from a list if no name is given)",
		Usage:    "[branch]",
		Examples: []string{"switch", "switch main"},
		MaxArgs:  1,
		Run:      runSwitch,
	}
}

func runSwitch(ctx *cli.Context) error {
	g, u := ctx.Git, ctx.UI
	if err := requireRepo(ctx); err != nil {
		return err
	}

	target := ctx.Arg(0)
	if target == "" {
		branches, err := g.Branches()
		if err != nil {
			return err
		}
		current, _ := g.CurrentBranch()
		var others []string
		for _, b := range branches {
			if b != current {
				others = append(others, b)
			}
		}
		if len(others) == 0 {
			return errors.New("there are no other branches; create one with 'rgit new <branch>'")
		}
		i, err := u.Choose(fmt.Sprintf("Switch from %s to:", ui.Bold(current)), others)
		if err != nil {
			return err
		}
		target = others[i]
	}

	// "git checkout <name>" also creates a local branch tracking the remote
	// one when only the remote has it.
	if _, err := g.Output("checkout", target); err != nil {
		return fmt.Errorf("could not switch to %q: %w", target, err)
	}
	u.Success("Switched to %s", ui.Bold(target))
	return nil
}
