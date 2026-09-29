package templates

import (
	"strings"
	"testing"
)

func TestEveryLicenseRenders(t *testing.T) {
	for _, l := range Licenses {
		text, err := l.Render("2026", "Jane Doe / Acme & Co")
		if err != nil {
			t.Errorf("%s: %v", l.ID, err)
			continue
		}
		if strings.Contains(text, "[year]") || strings.Contains(text, "[fullname]") {
			t.Errorf("%s: placeholders left in the output", l.ID)
		}
		if strings.Contains(text, "\r") {
			t.Errorf("%s: output should use \\n line endings", l.ID)
		}
		if detected, ok := DetectLicense(text); !ok || detected.ID != l.ID {
			t.Errorf("%s: DetectLicense = %v, %v", l.ID, detected.ID, ok)
		}
	}
}

func TestLicenseHolderIsInserted(t *testing.T) {
	l, _ := FindLicense("mit")
	text, _ := l.Render("2026", "Jane Doe / Acme & Co")
	if !strings.Contains(text, "Copyright (c) 2026 Jane Doe / Acme & Co") {
		t.Errorf("copyright line missing:\n%s", text[:200])
	}
}

func TestFindLicense(t *testing.T) {
	tests := map[string]string{"MIT": "mit", "apache": "apache-2.0", "gpl": "gpl-3.0", "cc": "cc-by-4.0", " Unlicense ": "unlicense"}
	for query, want := range tests {
		if l, ok := FindLicense(query); !ok || l.ID != want {
			t.Errorf("FindLicense(%q) = %q, %v; want %q", query, l.ID, ok, want)
		}
	}
	if _, ok := FindLicense("nope"); ok {
		t.Error("unknown license should not be found")
	}
	if _, ok := FindLicense(""); ok {
		t.Error("empty query should not match")
	}
}

func TestIgnorePresets(t *testing.T) {
	presets := IgnorePresets()
	if len(presets) < 10 {
		t.Fatalf("expected at least 10 presets, got %v", presets)
	}
	for _, name := range presets {
		lines, ok := IgnorePreset(name)
		if !ok || len(lines) == 0 {
			t.Errorf("preset %q is empty or missing", name)
		}
	}
	if lines, ok := IgnorePreset("JS"); !ok || !contains(lines, "node_modules/") {
		t.Error("alias js should resolve to the node preset")
	}
}

func TestRenderReadme(t *testing.T) {
	out, err := RenderReadme(Readme{
		Title: "Demo", Description: "A demo.", RunCommand: "go run .",
		Structure: ".\n└── main.go\n", License: "MIT License",
		Owner: "jane", RepoName: "demo", RepoURL: "https://github.com/jane/demo",
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"# Demo", "A demo.", "go run .", "└── main.go", "git clone https://github.com/jane/demo.git", "**MIT License**", "[@jane]"} {
		if !strings.Contains(out, want) {
			t.Errorf("README is missing %q:\n%s", want, out)
		}
	}

	minimal, err := RenderReadme(Readme{Title: "Solo", Structure: ".\n"})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(minimal, "github.com") || !strings.Contains(minimal, "rgit license") {
		t.Errorf("minimal README should skip GitHub parts and suggest a license:\n%s", minimal)
	}
}

func contains(lines []string, want string) bool {
	for _, l := range lines {
		if l == want {
			return true
		}
	}
	return false
}
