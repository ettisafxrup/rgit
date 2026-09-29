package commands

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ettisafxrup/rgit/internal/cli"
	"github.com/ettisafxrup/rgit/internal/tree"
	"github.com/ettisafxrup/rgit/internal/ui"
)

func archiCommand() *cli.Command {
	var output string
	var depth int
	var print bool
	return &cli.Command{
		Name:    "archi",
		Aliases: []string{"tree"},
		Group:   "Project files",
		Summary: "Write the project's folder structure to ARCHITECTURE",
		Details: `Inside a git repository the tree contains exactly the files git sees,
so everything in .gitignore is left out automatically. Elsewhere,
common clutter (node_modules, build output, hidden files, logs) is skipped.`,
		Examples: []string{"archi", "archi --print", "archi --depth 2 -o STRUCTURE.md"},
		MaxArgs:  0,
		Flags: func(fs *flag.FlagSet) {
			fs.StringVar(&output, "o", "ARCHITECTURE", "write the tree to `file`")
			fs.IntVar(&depth, "depth", 0, "limit the tree to `levels` deep (0 = no limit)")
			fs.BoolVar(&print, "print", false, "print the tree instead of writing a file")
		},
		Run: func(ctx *cli.Context) error {
			structure, err := projectTree(ctx, depth, output)
			if err != nil {
				return err
			}
			if print {
				ctx.UI.Printf("%s", structure)
				return nil
			}
			content := fmt.Sprintf("Project structure of %s\n\n%s", projectName(), structure)
			if err := os.WriteFile(output, []byte(content), 0o644); err != nil {
				return err
			}
			ctx.UI.Success("Saved the project structure to %s", ui.Bold(output))
			return nil
		},
	}
}

// projectTree renders the files of the current directory as a tree,
// leaving out the file named exclude.
func projectTree(ctx *cli.Context, depth int, exclude string) (string, error) {
	var paths []string
	var err error
	if git := ctx.Git; git.IsRepo() {
		// Tracked files plus new files that are not ignored, relative to here.
		var out string
		out, err = git.RawOutput("ls-files", "--cached", "--others", "--exclude-standard", "-z")
		for _, p := range strings.Split(out, "\x00") {
			if p != "" {
				paths = append(paths, p)
			}
		}
	} else {
		paths, err = tree.Walk(".", tree.IsNoise)
	}
	if err != nil {
		return "", err
	}

	exclude = filepath.ToSlash(filepath.Clean(exclude))
	kept := paths[:0]
	for _, p := range paths {
		if p != exclude {
			kept = append(kept, p)
		}
	}
	return tree.Render(tree.Build(kept), depth), nil
}

// projectName is the name of the current directory.
func projectName() string {
	dir, err := os.Getwd()
	if err != nil {
		return "project"
	}
	return filepath.Base(dir)
}
