package commands

import (
	"fmt"
	"strings"

	"github.com/ettisafxrup/rgit/internal/cli"
	"github.com/ettisafxrup/rgit/internal/ui"
)

func statusCommand() *cli.Command {
	return &cli.Command{
		Name:     "status",
		Aliases:  []string{"st"},
		Group:    "Everyday",
		Summary:  "Show a clear overview of the repository",
		Details:  "Shows the branch, how it compares to the remote, the last commit,\nyour uncommitted changes and stashed work, all on one screen.",
		Examples: []string{"status", "st"},
		MaxArgs:  0,
		Run:      runStatus,
	}
}

func runStatus(ctx *cli.Context) error {
	g, u := ctx.Git, ctx.UI
	if err := requireRepo(ctx); err != nil {
		return err
	}

	u.Println()
	branch, err := g.CurrentBranch()
	if err != nil {
		branch = ui.Red("detached HEAD")
	}
	u.Field("Branch", ui.Bold(branch)+trackingInfo(ctx))

	if last, err := g.LastCommit(); err == nil {
		u.Field("Last commit", fmt.Sprintf("%s %s %s", ui.Yellow(last.Hash), last.Subject, ui.Dim("("+last.When+")")))
	} else {
		u.Field("Last commit", ui.Dim("none yet"))
	}

	if remote := g.DefaultRemote(); remote != "" {
		url, _ := g.RemoteURL(remote)
		u.Field("Remote", fmt.Sprintf("%s %s", remote, ui.Dim(url)))
	} else {
		u.Field("Remote", ui.Dim("none — 'rgit push' will help you add one"))
	}

	status, err := g.Status()
	if err != nil {
		return err
	}
	if status.Clean() {
		u.Field("Changes", ui.Green("working tree clean"))
	} else {
		u.Field("Changes", changeSummary(status.Staged, status.Modified, status.Untracked, status.Conflicts))
		printChanges(u, status.Changes, 20)
	}

	if n := g.StashCount(); n > 0 {
		u.Field("Stash", fmt.Sprintf("%d %s", n, plural(n, "entry", "entries")))
	}
	u.Println()
	return nil
}

// trackingInfo describes the upstream and how far the branch has diverged.
func trackingInfo(ctx *cli.Context) string {
	upstream, ok := ctx.Git.Upstream()
	if !ok {
		return ui.Dim("  (not pushed yet)")
	}
	ahead, behind, err := ctx.Git.AheadBehind()
	if err != nil {
		return ui.Dim("  → " + upstream)
	}
	var state string
	switch {
	case ahead == 0 && behind == 0:
		state = ui.Green("up to date")
	case behind == 0:
		state = ui.Yellow(fmt.Sprintf("%d to push", ahead))
	case ahead == 0:
		state = ui.Cyan(fmt.Sprintf("%d to pull", behind))
	default:
		state = ui.Red(fmt.Sprintf("diverged: %d to push, %d to pull", ahead, behind))
	}
	return fmt.Sprintf("  → %s  %s", upstream, state)
}

func changeSummary(staged, modified, untracked, conflicts int) string {
	var parts []string
	if conflicts > 0 {
		parts = append(parts, ui.Red(fmt.Sprintf("%d conflicted", conflicts)))
	}
	if staged > 0 {
		parts = append(parts, ui.Green(fmt.Sprintf("%d staged", staged)))
	}
	if modified > 0 {
		parts = append(parts, ui.Yellow(fmt.Sprintf("%d modified", modified)))
	}
	if untracked > 0 {
		parts = append(parts, ui.Cyan(fmt.Sprintf("%d new", untracked)))
	}
	return strings.Join(parts, ", ")
}
