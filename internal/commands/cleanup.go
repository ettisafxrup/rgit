package commands

import (
	"flag"
	"fmt"

	"github.com/ettisafxrup/rgit/internal/cli"
	"github.com/ettisafxrup/rgit/internal/ui"
)

// protectedBranches are never deleted by cleanup.
var protectedBranches = map[string]bool{"main": true, "master": true, "develop": true}

func cleanupCommand() *cli.Command {
	var prune bool
	return &cli.Command{
		Name:    "cleanup",
		Aliases: []string{"clean"},
		Group:   "Branches",
		Summary: "Delete local branches that are already merged",
		Usage:   "[base-branch]",
		Details: `Finds local branches whose work is fully merged into the base branch
(the repository's default branch unless you name one) and deletes them
after showing you the list. main, master, develop, the base branch and
the branch you are on are always kept.`,
		Examples: []string{"cleanup", "cleanup develop", "cleanup --prune"},
		MaxArgs:  1,
		Flags: func(fs *flag.FlagSet) {
			fs.BoolVar(&prune, "prune", false, "also forget remote branches that were deleted on the server")
		},
		Run: func(ctx *cli.Context) error {
			return runCleanup(ctx, prune)
		},
	}
}

func runCleanup(ctx *cli.Context, prune bool) error {
	g, u := ctx.Git, ctx.UI
	if err := requireRepo(ctx); err != nil {
		return err
	}

	if prune {
		if remote := g.DefaultRemote(); remote != "" {
			u.Step("Pruning deleted branches of %s", remote)
			if err := g.Run("fetch", "--prune", remote); err != nil {
				return err
			}
		}
	}

	base := ctx.Arg(0)
	if base == "" {
		base = g.DefaultBranch()
	}
	if !g.Succeeds("rev-parse", "--verify", "--quiet", base) {
		return fmt.Errorf("base branch %q does not exist", base)
	}
	current, _ := g.CurrentBranch()

	merged, err := g.MergedBranches(base)
	if err != nil {
		return err
	}
	var stale []string
	for _, b := range merged {
		if b != base && b != current && !protectedBranches[b] {
			stale = append(stale, b)
		}
	}
	if len(stale) == 0 {
		u.Success("Nothing to clean up: no merged branches besides %s.", ui.Bold(base))
		return nil
	}

	u.Info("Branches fully merged into %s:", ui.Bold(base))
	for _, b := range stale {
		u.Printf("   %s %s\n", ui.Red("-"), b)
	}
	ok, err := u.Confirm(fmt.Sprintf("Delete these %d %s?", len(stale), plural(len(stale), "branch", "branches")), true)
	if err != nil {
		return err
	}
	if !ok {
		return cli.ErrAborted
	}

	for _, b := range stale {
		// -D is safe here: every branch was verified as merged into base.
		// (-d would compare against HEAD instead and could refuse.)
		if _, err := g.Output("branch", "-D", b); err != nil {
			u.Warn("Could not delete %s: %v", b, err)
			continue
		}
		u.Success("Deleted %s", b)
	}
	return nil
}
