package commands

import (
	"errors"
	"flag"
	"os"
	"strings"

	"github.com/ettisafxrup/rgit/internal/cli"
	"github.com/ettisafxrup/rgit/internal/templates"
	"github.com/ettisafxrup/rgit/internal/ui"
)

const gitignoreFile = ".gitignore"

func ignoreCommand() *cli.Command {
	var list bool
	return &cli.Command{
		Name:    "ignore",
		Group:   "Project files",
		Summary: "Build a .gitignore from presets, patterns or a file picker",
		Usage:   "[preset | pattern]...",
		Details: `Each argument is either a preset name (go, node, python, …) or a
pattern such as "*.log" or "secrets/". Without arguments, rgit lets you
pick files and folders from the current directory and then offers presets.

Entries that are already in .gitignore are never added twice.`,
		Examples: []string{"ignore", "ignore go vscode windows", `ignore node "*.pem" tmp/`, "ignore --list"},
		MaxArgs:  -1,
		Flags: func(fs *flag.FlagSet) {
			fs.BoolVar(&list, "list", false, "list the available presets")
		},
		Run: func(ctx *cli.Context) error {
			if list {
				ctx.UI.Info("Available presets: %s", strings.Join(templates.IgnorePresets(), ", "))
				return nil
			}
			return runIgnore(ctx)
		},
	}
}

func runIgnore(ctx *cli.Context) error {
	existing, err := os.ReadFile(gitignoreFile)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	var blocks [][]string
	if len(ctx.Args) > 0 {
		blocks = blocksFromArgs(ctx.Args)
	} else {
		if blocks, err = pickIgnoreBlocks(ctx); err != nil {
			return err
		}
	}

	merged, added := mergeIgnore(string(existing), blocks)
	if added == 0 {
		ctx.UI.Success("%s already has everything you asked for.", gitignoreFile)
		return nil
	}
	if err := os.WriteFile(gitignoreFile, []byte(merged), 0o644); err != nil {
		return err
	}
	ctx.UI.Success("Added %d %s to %s", added, plural(added, "entry", "entries"), ui.Bold(gitignoreFile))
	return nil
}

// blocksFromArgs turns arguments into preset blocks and a block of custom patterns.
func blocksFromArgs(args []string) [][]string {
	var blocks [][]string
	var custom []string
	for _, arg := range args {
		if preset, ok := templates.IgnorePreset(arg); ok {
			blocks = append(blocks, preset)
		} else {
			custom = append(custom, arg)
		}
	}
	if len(custom) > 0 {
		blocks = append(blocks, append([]string{"# Added with rgit"}, custom...))
	}
	return blocks
}

// pickIgnoreBlocks lets the user choose entries of the current directory and presets.
func pickIgnoreBlocks(ctx *cli.Context) ([][]string, error) {
	u := ctx.UI
	var blocks [][]string

	entries, err := os.ReadDir(".")
	if err != nil {
		return nil, err
	}
	var names, patterns []string
	for _, e := range entries {
		if e.Name() == ".git" || e.Name() == gitignoreFile {
			continue
		}
		if e.IsDir() {
			names = append(names, e.Name()+"/")
			patterns = append(patterns, "/"+e.Name()+"/")
		} else {
			names = append(names, e.Name())
			patterns = append(patterns, "/"+e.Name())
		}
	}

	if len(names) > 0 {
		picked, err := u.ChooseMany("Which files or folders should git ignore?", names)
		if err != nil {
			return nil, err
		}
		if len(picked) > 0 {
			block := []string{"# Picked with rgit"}
			for _, i := range picked {
				block = append(block, patterns[i])
			}
			blocks = append(blocks, block)
		}
	}

	presets := templates.IgnorePresets()
	u.Info("Presets: %s", strings.Join(presets, ", "))
	answer, err := u.Ask("Add presets? (names separated by spaces, empty to skip)", "")
	if err != nil {
		return nil, err
	}
	for _, name := range strings.Fields(strings.ReplaceAll(answer, ",", " ")) {
		preset, ok := templates.IgnorePreset(name)
		if !ok {
			u.Warn("Unknown preset %q, skipped.", name)
			continue
		}
		blocks = append(blocks, preset)
	}
	return blocks, nil
}

// mergeIgnore appends blocks of .gitignore lines to existing content and
// returns the result with the number of patterns added. Patterns already
// present are skipped, and a "# comment" heading is written only when at
// least one pattern below it is new.
func mergeIgnore(existing string, blocks [][]string) (string, int) {
	existing = strings.ReplaceAll(existing, "\r\n", "\n")
	present := map[string]bool{}
	for _, line := range strings.Split(existing, "\n") {
		present[strings.TrimSpace(line)] = true
	}

	var b strings.Builder
	b.WriteString(existing)
	added := 0
	for _, block := range blocks {
		heading, sectionOpen := "", false
		for _, raw := range block {
			line := strings.TrimSpace(raw)
			switch {
			case line == "" || present[line]:
				continue
			case strings.HasPrefix(line, "#"):
				heading = line
				continue
			}
			if heading != "" || !sectionOpen {
				startSection(&b)
				if heading != "" {
					b.WriteString(heading + "\n")
				}
				heading, sectionOpen = "", true
			}
			b.WriteString(line + "\n")
			present[line] = true
			added++
		}
	}
	return b.String(), added
}

// startSection separates a new group of lines from what came before.
func startSection(b *strings.Builder) {
	if b.Len() == 0 {
		return
	}
	if !strings.HasSuffix(b.String(), "\n") {
		b.WriteString("\n")
	}
	b.WriteString("\n")
}
