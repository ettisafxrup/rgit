package commands

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ettisafxrup/rgit/internal/cli"
	"github.com/ettisafxrup/rgit/internal/templates"
	"github.com/ettisafxrup/rgit/internal/ui"
)

const licenseFile = "LICENSE"

func licenseCommand() *cli.Command {
	var opts licenseOptions
	return &cli.Command{
		Name:    "license",
		Group:   "Project files",
		Summary: "Add a LICENSE file (MIT, Apache, GPL and more)",
		Usage:   "[license]",
		Details: "Available licenses: " + licenseIDs() + `.

The copyright year defaults to this year and the holder to your git
user.name. Pass a license name to skip the menu.`,
		Examples: []string{"license", "license mit", `license apache-2.0 --name "Acme Inc."`},
		MaxArgs:  1,
		Flags: func(fs *flag.FlagSet) {
			fs.StringVar(&opts.holder, "name", "", "copyright `holder` (person or organization)")
			fs.StringVar(&opts.year, "year", "", "copyright `year`")
		},
		Run: func(ctx *cli.Context) error {
			opts.id = ctx.Arg(0)
			_, err := createLicense(ctx, opts)
			return err
		},
	}
}

type licenseOptions struct {
	id, holder, year string
}

// createLicense writes a LICENSE file and returns the chosen license.
func createLicense(ctx *cli.Context, opts licenseOptions) (templates.License, error) {
	u := ctx.UI
	if err := confirmOverwrite(ctx, licenseFile); err != nil {
		return templates.License{}, err
	}

	license, err := chooseLicense(ctx, opts.id)
	if err != nil {
		return license, err
	}

	year := opts.year
	if year == "" {
		if year, err = u.Ask("Copyright year:", strconv.Itoa(time.Now().Year())); err != nil {
			return license, err
		}
	}
	holder := opts.holder
	if holder == "" {
		if holder, err = u.AskRequired("Copyright holder (your name or organization):", ctx.Git.Config("user.name")); err != nil {
			return license, err
		}
	}

	text, err := license.Render(year, holder)
	if err != nil {
		return license, err
	}
	if err := os.WriteFile(licenseFile, []byte(text), 0o644); err != nil {
		return license, err
	}
	u.Success("Created %s with the %s", ui.Bold(licenseFile), license.Name)
	return license, nil
}

func chooseLicense(ctx *cli.Context, id string) (templates.License, error) {
	if id != "" {
		license, ok := templates.FindLicense(id)
		if !ok {
			return license, fmt.Errorf("unknown license %q (available: %s)", id, licenseIDs())
		}
		return license, nil
	}
	names := make([]string, len(templates.Licenses))
	for i, l := range templates.Licenses {
		names[i] = fmt.Sprintf("%-14s %s", l.ID, ui.Dim(l.Name))
	}
	i, err := ctx.UI.Choose("Which license do you want?", names)
	if err != nil {
		return templates.License{}, err
	}
	return templates.Licenses[i], nil
}

func licenseIDs() string {
	ids := make([]string, len(templates.Licenses))
	for i, l := range templates.Licenses {
		ids[i] = l.ID
	}
	return strings.Join(ids, ", ")
}

// detectLicense returns a readable name for an existing LICENSE file, or "".
func detectLicense() string {
	data, err := os.ReadFile(licenseFile)
	if err != nil {
		return ""
	}
	if license, ok := templates.DetectLicense(string(data)); ok {
		return license.Name
	}
	// An unknown license usually names itself on its first line.
	return templates.FirstLine(string(data))
}
