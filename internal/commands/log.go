package commands

import (
	"errors"
	"flag"
	"strconv"

	"github.com/ettisafxrup/rgit/internal/cli"
	"github.com/ettisafxrup/rgit/internal/ui"
)

func logCommand() *cli.Command {
	var count int
	var all bool
	return &cli.Command{
		Name:     "log",
		Aliases:  []string{"l"},
		Group:    "Everyday",
		Summary:  "Show a compact, colorful commit graph",
		Examples: []string{"log", "log -n 30", "log --all"},
		MaxArgs:  0,
		Flags: func(fs *flag.FlagSet) {
			fs.IntVar(&count, "n", 15, "show at most `count` commits")
			fs.BoolVar(&all, "all", false, "include every branch, not just the current one")
		},
		Run: func(ctx *cli.Context) error {
			if err := requireRepo(ctx); err != nil {
				return err
			}
			if !ctx.Git.HasCommits() {
				return errors.New("this repository has no commits yet")
			}
			color := "--color=never"
			if ui.ColorEnabled() {
				color = "--color=always"
			}
			args := []string{
				"--no-pager", "log", "--graph", color, "-n", strconv.Itoa(max(count, 1)),
				"--format=%C(yellow)%h%C(reset) %s %C(dim)— %an, %cr%C(reset)%C(auto)%d",
			}
			if all {
				args = append(args, "--all")
			}
			return ctx.Git.Run(args...)
		},
	}
}
