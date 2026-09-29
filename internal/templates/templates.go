// Package templates provides the files rgit generates: licenses,
// .gitignore presets and the README skeleton. Everything is embedded in
// the binary, so rgit works as a single self-contained executable.
package templates

import (
	"bytes"
	"embed"
	"fmt"
	"sort"
	"strings"
	"text/template"
)

//go:embed licenses/*.txt gitignore/*.gitignore readme.md.tmpl
var files embed.FS

// License is one of the bundled license texts.
type License struct {
	ID   string // SPDX-style identifier, also the template file name
	Name string
}

// Licenses lists the bundled licenses in the order they are offered.
var Licenses = []License{
	{"mit", "MIT License"},
	{"apache-2.0", "Apache License 2.0"},
	{"gpl-3.0", "GNU General Public License v3.0"},
	{"bsd-3-clause", "BSD 3-Clause License"},
	{"bsl-1.0", "Boost Software License 1.0"},
	{"agpl-3.0", "GNU Affero General Public License v3.0"},
	{"lgpl-3.0", "GNU Lesser General Public License v3.0"},
	{"mpl-2.0", "Mozilla Public License 2.0"},
	{"unlicense", "The Unlicense"},
	{"cc-by-4.0", "Creative Commons Attribution 4.0"},
}

// FindLicense looks a license up by identifier. It is forgiving: "MIT",
// "apache" and "gpl" all work.
func FindLicense(query string) (License, bool) {
	q := strings.ToLower(strings.TrimSpace(query))
	if q == "" {
		return License{}, false
	}
	for _, l := range Licenses {
		if l.ID == q || strings.ToLower(l.Name) == q {
			return l, true
		}
	}
	for _, l := range Licenses {
		if strings.HasPrefix(l.ID, q) {
			return l, true
		}
	}
	return License{}, false
}

// Render returns the license text with the year and copyright holder filled in.
func (l License) Render(year, holder string) (string, error) {
	data, err := files.ReadFile("licenses/" + l.ID + ".txt")
	if err != nil {
		return "", fmt.Errorf("license template %q is missing: %w", l.ID, err)
	}
	text := strings.ReplaceAll(string(data), "\r\n", "\n")
	text = strings.ReplaceAll(text, "[year]", year)
	text = strings.ReplaceAll(text, "[fullname]", holder)
	return text, nil
}

// DetectLicense recognizes which bundled license a text is, by comparing
// its first line, which is distinctive for every bundled license.
func DetectLicense(text string) (License, bool) {
	first := normalizeSpace(FirstLine(text))
	if first == "" {
		return License{}, false
	}
	for _, l := range Licenses {
		data, err := files.ReadFile("licenses/" + l.ID + ".txt")
		if err == nil && normalizeSpace(FirstLine(string(data))) == first {
			return l, true
		}
	}
	return License{}, false
}

// FirstLine returns the first non-blank line of text, trimmed.
func FirstLine(text string) string {
	for _, line := range strings.Split(text, "\n") {
		if line = strings.TrimSpace(line); line != "" {
			return line
		}
	}
	return ""
}

func normalizeSpace(s string) string {
	return strings.ToLower(strings.Join(strings.Fields(s), " "))
}

// ignoreAliases maps common alternative names to preset names.
var ignoreAliases = map[string]string{
	"golang": "go", "js": "node", "javascript": "node", "ts": "node", "typescript": "node",
	"py": "python", "csharp": "dotnet", "c#": "dotnet", "c++": "cpp", "c": "cpp",
	"mac": "macos", "osx": "macos", "win": "windows", "idea": "jetbrains", "code": "vscode",
}

// IgnorePresets lists the names of the bundled .gitignore presets.
func IgnorePresets() []string {
	entries, _ := files.ReadDir("gitignore")
	names := make([]string, 0, len(entries))
	for _, e := range entries {
		names = append(names, strings.TrimSuffix(e.Name(), ".gitignore"))
	}
	sort.Strings(names)
	return names
}

// IgnorePreset returns the lines of a .gitignore preset.
func IgnorePreset(name string) ([]string, bool) {
	name = strings.ToLower(strings.TrimSpace(name))
	if alias, ok := ignoreAliases[name]; ok {
		name = alias
	}
	data, err := files.ReadFile("gitignore/" + name + ".gitignore")
	if err != nil {
		return nil, false
	}
	text := strings.ReplaceAll(string(data), "\r\n", "\n")
	return strings.Split(strings.TrimRight(text, "\n"), "\n"), true
}

// Readme holds the values used to fill in the README template.
type Readme struct {
	Title       string
	Description string
	RunCommand  string
	Structure   string
	License     string
	Owner       string // GitHub user or organization, optional
	RepoName    string
	RepoURL     string // web URL of the repository, optional
}

// RenderReadme produces a README.md from the given values.
func RenderReadme(r Readme) (string, error) {
	tmpl, err := template.ParseFS(files, "readme.md.tmpl")
	if err != nil {
		return "", err
	}
	var out bytes.Buffer
	if err := tmpl.Execute(&out, r); err != nil {
		return "", err
	}
	return out.String(), nil
}
