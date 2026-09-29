package commands

import (
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/ettisafxrup/rgit/internal/cli"
	"github.com/ettisafxrup/rgit/internal/git"
	"github.com/ettisafxrup/rgit/internal/ui"
	"github.com/ettisafxrup/rgit/internal/version"
)

func doctorCommand() *cli.Command {
	return &cli.Command{
		Name:     "doctor",
		Group:    "Setup",
		Summary:  "Check that git and rgit are set up correctly",
		Details:  "Checks the git installation, your identity, credential storage and,\ninside a repository, the branch and remote. Every problem comes with a fix.",
		Examples: []string{"doctor"},
		MaxArgs:  0,
		Run:      runDoctor,
	}
}

func runDoctor(ctx *cli.Context) error {
	g, u := ctx.Git, ctx.UI
	u.Title("rgit doctor")
	problems := 0
	check := func(ok bool, good, bad, fix string) {
		if ok {
			u.Success("%s", good)
			return
		}
		problems++
		u.Error("%s", bad)
		if fix != "" {
			u.Hint("%s", fix)
		}
	}

	gitErr := git.Installed()
	gitVersion, _ := g.Output("--version")
	check(gitErr == nil, strings.TrimPrefix(gitVersion, "git version ")+" installed",
		"git is not installed", "Install it from https://git-scm.com and reopen your terminal.")
	if gitErr != nil {
		return errors.New("git is required")
	}

	name, email := g.Config("user.name"), g.Config("user.email")
	check(name != "" && email != "", fmt.Sprintf("Identity: %s <%s>", name, email),
		"Your git identity is incomplete", "Run 'rgit login'.")

	helper := g.Config("credential.helper")
	check(helper != "", "Credential helper: "+helper,
		"No credential helper: git will ask for your password every time",
		"Git for Windows includes one: git config --global credential.helper manager")

	_, lookErr := exec.LookPath("rgit")
	check(lookErr == nil, "rgit is on your PATH", "rgit is not on your PATH",
		"Add the folder containing rgit to PATH, or reinstall with the installer.")

	if g.IsRepo() {
		branch, err := g.CurrentBranch()
		check(err == nil, "On branch "+branch, "You are in 'detached HEAD' state", "Run 'rgit switch' to pick a branch.")

		remote := g.DefaultRemote()
		remoteURL, _ := g.RemoteURL(remote)
		check(remote != "", fmt.Sprintf("Remote %s → %s", remote, remoteURL),
			"This repository has no remote", "Run 'rgit push' and rgit will help you add one.")

		status, _ := g.Status()
		check(status.Conflicts == 0, "No merge conflicts",
			fmt.Sprintf("%d files have merge conflicts", status.Conflicts), "Fix them, then 'git add' the files.")
	} else {
		u.Info("Not inside a repository; skipped repository checks.")
	}

	u.Printf("\n   %s\n\n", ui.Dim(fmt.Sprintf("rgit %s · %s/%s", version.Version, runtime.GOOS, runtime.GOARCH)))
	if problems > 0 {
		return fmt.Errorf("found %d %s", problems, plural(problems, "problem", "problems"))
	}
	u.Success("Everything looks good!")
	return nil
}

func versionCommand() *cli.Command {
	return &cli.Command{
		Name:    "version",
		Group:   "Setup",
		Summary: "Print the rgit version",
		MaxArgs: 0,
		Run: func(ctx *cli.Context) error {
			ctx.UI.Printf("rgit %s (%s/%s, %s)\n", version.Version, runtime.GOOS, runtime.GOARCH, runtime.Version())
			return nil
		},
	}
}
