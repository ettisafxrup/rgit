package commands

import (
	"errors"
	"fmt"
	"net/mail"
	"path"
	"strings"

	"github.com/ettisafxrup/rgit/internal/cli"
	"github.com/ettisafxrup/rgit/internal/git"
	"github.com/ettisafxrup/rgit/internal/ui"
)

func cloneCommand() *cli.Command {
	return &cli.Command{
		Name:     "clone",
		Group:    "Setup",
		Summary:  "Clone a repository (accepts GitHub owner/repo shorthand)",
		Usage:    "<repository> [folder]",
		Examples: []string{"clone ettisafxrup/rgit", "clone https://gitlab.com/group/project.git my-folder"},
		MaxArgs:  2,
		Run: func(ctx *cli.Context) error {
			if err := git.Installed(); err != nil {
				return err
			}
			target := ctx.Arg(0)
			if target == "" {
				return &cli.UsageError{Message: "tell rgit which repository to clone"}
			}
			url := git.ExpandCloneURL(target)
			folder := ctx.Arg(1)
			if folder == "" {
				folder = strings.TrimSuffix(path.Base(strings.ReplaceAll(url, ":", "/")), ".git")
			}

			ctx.UI.Step("Cloning %s", ui.Bold(url))
			if err := ctx.Git.Run("clone", url, folder); err != nil {
				return err
			}
			ctx.UI.Success("Cloned into %s", ui.Bold(folder))
			ctx.UI.Hint("cd %s", folder)
			return nil
		},
	}
}

func loginCommand() *cli.Command {
	return &cli.Command{
		Name:    "login",
		Group:   "Setup",
		Summary: "Set the name and email git records in your commits",
		Details: `Saves your identity in the global git configuration. Inside a repository
with a remote, rgit also checks that you can reach that remote.`,
		Examples: []string{"login"},
		MaxArgs:  0,
		Run:      runLogin,
	}
}

func runLogin(ctx *cli.Context) error {
	g, u := ctx.Git, ctx.UI
	if err := git.Installed(); err != nil {
		return err
	}

	u.Title("Git identity")
	name, err := u.AskRequired("Your name:", g.GlobalConfig("user.name"))
	if err != nil {
		return err
	}
	var email string
	for {
		if email, err = u.AskRequired("Your email:", g.GlobalConfig("user.email")); err != nil {
			return err
		}
		if _, parseErr := mail.ParseAddress(email); parseErr == nil {
			break
		}
		u.Warn("%q does not look like an email address.", email)
		if u.AssumeYes {
			return errors.New("invalid email address")
		}
	}

	if err := g.SetGlobalConfig("user.name", name); err != nil {
		return err
	}
	if err := g.SetGlobalConfig("user.email", email); err != nil {
		return err
	}
	u.Success("Commits will be signed as %s <%s>", ui.Bold(name), email)

	if !g.IsRepo() {
		return nil
	}
	remote := g.DefaultRemote()
	if remote == "" {
		return nil
	}
	u.Step("Checking access to %s", remote)
	if _, err := g.Output("ls-remote", "--heads", remote); err != nil {
		u.Warn("Could not reach %s. You may need to sign in or set up a token / SSH key.", remote)
		u.Hint("Try 'git fetch %s' to see the full message.", remote)
		return nil
	}
	u.Success("You can access %s", remote)
	return nil
}

func whoamiCommand() *cli.Command {
	return &cli.Command{
		Name:     "whoami",
		Aliases:  []string{"userinfo", "me"},
		Group:    "Setup",
		Summary:  "Show the git identity used for your commits",
		Examples: []string{"whoami"},
		MaxArgs:  0,
		Run: func(ctx *cli.Context) error {
			g, u := ctx.Git, ctx.UI
			if err := git.Installed(); err != nil {
				return err
			}
			notSet := ui.Red("not set")
			valueOr := func(v string) string {
				if v == "" {
					return notSet
				}
				return v
			}

			u.Title("Git identity")
			u.Field("Name", valueOr(g.Config("user.name")))
			u.Field("Email", valueOr(g.Config("user.email")))
			if g.IsRepo() {
				if local, _ := g.Output("config", "--local", "--get", "user.email"); local != "" {
					u.Field("", ui.Dim("(set for this repository only)"))
				}
			}
			if g.Config("user.name") == "" || g.Config("user.email") == "" {
				u.Println()
				u.Hint("Run 'rgit login' to set your identity.")
			}
			u.Println()
			return nil
		},
	}
}

func openCommand() *cli.Command {
	return &cli.Command{
		Name:     "open",
		Aliases:  []string{"web"},
		Group:    "Setup",
		Summary:  "Open the repository's web page in your browser",
		Usage:    "[remote]",
		Details:  "Works with GitHub, GitLab, Bitbucket and most other hosts. On GitHub,\nthe page of the current branch is opened.",
		Examples: []string{"open", "open upstream"},
		MaxArgs:  1,
		Run: func(ctx *cli.Context) error {
			if err := requireRepo(ctx); err != nil {
				return err
			}
			remote := ctx.Arg(0)
			if remote == "" {
				if remote = ctx.Git.DefaultRemote(); remote == "" {
					return errors.New("this repository has no remote to open")
				}
			}
			remoteURL, err := ctx.Git.RemoteURL(remote)
			if err != nil {
				return fmt.Errorf("no remote named %q", remote)
			}
			page, err := git.WebURL(remoteURL)
			if err != nil {
				return err
			}
			if _, _, isGitHub := git.GitHubRepo(remoteURL); isGitHub {
				if branch, err := ctx.Git.CurrentBranch(); err == nil && branch != ctx.Git.DefaultBranch() {
					page += "/tree/" + branch
				}
			}

			ctx.UI.Step("Opening %s", ui.Bold(page))
			return openBrowser(page)
		},
	}
}
