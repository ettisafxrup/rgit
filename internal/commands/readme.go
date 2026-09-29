package commands

import (
	"os"
	"strings"

	"github.com/ettisafxrup/rgit/internal/cli"
	"github.com/ettisafxrup/rgit/internal/git"
	"github.com/ettisafxrup/rgit/internal/templates"
	"github.com/ettisafxrup/rgit/internal/ui"
)

const readmeFile = "README.md"

func readmeCommand() *cli.Command {
	return &cli.Command{
		Name:    "readme",
		Group:   "Project files",
		Summary: "Generate a README.md for your project",
		Details: `Asks a few questions and writes a README with a project description,
setup instructions, the project structure, contribution guidelines and
license information. GitHub details are taken from your remote when
possible, and you are offered a license if the project has none.`,
		Examples: []string{"readme", "readme -y"},
		MaxArgs:  0,
		Run:      runReadme,
	}
}

func runReadme(ctx *cli.Context) error {
	u := ctx.UI
	if err := confirmOverwrite(ctx, readmeFile); err != nil {
		return err
	}

	name := projectName()
	var data templates.Readme
	var err error
	if data.Title, err = u.Ask("Project title:", name); err != nil {
		return err
	}
	if data.Description, err = u.Ask("Short description:", "A short description of "+data.Title+"."); err != nil {
		return err
	}
	if data.RunCommand, err = u.Ask("Command to run it (e.g. npm run dev, optional):", ""); err != nil {
		return err
	}

	if err := fillRepository(ctx, &data, name); err != nil {
		return err
	}
	if data.License, err = readmeLicense(ctx); err != nil {
		return err
	}
	if data.Structure, err = projectTree(ctx, 3, readmeFile); err != nil {
		return err
	}

	content, err := templates.RenderReadme(data)
	if err != nil {
		return err
	}
	if err := os.WriteFile(readmeFile, []byte(content), 0o644); err != nil {
		return err
	}
	u.Success("Created %s", ui.Bold(readmeFile))
	return nil
}

// fillRepository sets the GitHub owner and repository, taken from the
// remote when it points at GitHub, otherwise asked for.
func fillRepository(ctx *cli.Context, data *templates.Readme, fallbackName string) error {
	if git.Installed() == nil && ctx.Git.IsRepo() {
		if remote := ctx.Git.DefaultRemote(); remote != "" {
			url, _ := ctx.Git.RemoteURL(remote)
			if owner, repo, ok := git.GitHubRepo(url); ok {
				setRepository(data, owner, repo)
				return nil
			}
		}
	}

	answer, err := ctx.UI.Ask("GitHub repository as owner/repo (optional):", "")
	if err != nil || answer == "" {
		return err
	}
	owner, repo, _ := strings.Cut(strings.Trim(answer, "/ "), "/")
	if owner == "" {
		return nil
	}
	if repo == "" {
		repo = fallbackName
	}
	setRepository(data, owner, repo)
	return nil
}

func setRepository(data *templates.Readme, owner, repo string) {
	data.Owner = owner
	data.RepoName = repo
	data.RepoURL = "https://github.com/" + owner + "/" + repo
}

// readmeLicense names the project's license, offering to create one.
func readmeLicense(ctx *cli.Context) (string, error) {
	if name := detectLicense(); name != "" {
		ctx.UI.Info("Found license: %s", name)
		return name, nil
	}
	ok, err := ctx.UI.Confirm("The project has no LICENSE yet. Add one now?", true)
	if err != nil || !ok {
		return "", err
	}
	license, err := createLicense(ctx, licenseOptions{})
	if err != nil {
		return "", err
	}
	return license.Name, nil
}
