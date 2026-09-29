// Package testenv provides helpers for tests that run real git commands.
package testenv

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// Isolate points git at a fresh global configuration for the duration of
// the test, so tests never read or modify the developer's real settings.
func Isolate(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not installed")
	}
	config := filepath.Join(t.TempDir(), "gitconfig")
	WriteFile(t, config, "[user]\n\tname = Test User\n\temail = test@example.com\n"+
		"[init]\n\tdefaultBranch = master\n[advice]\n\tdetachedHead = false\n")
	t.Setenv("GIT_CONFIG_GLOBAL", config)
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")
	t.Setenv("GIT_TERMINAL_PROMPT", "0")
}

// BareRemote creates an empty bare repository whose HEAD names branch.
func BareRemote(t *testing.T, branch string) string {
	t.Helper()
	dir := filepath.Join(t.TempDir(), "remote.git")
	Git(t, "", "init", "--bare", dir)
	Git(t, dir, "symbolic-ref", "HEAD", "refs/heads/"+branch)
	return dir
}

// Git runs a git command and fails the test if it does not succeed.
func Git(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return string(out)
}

// WriteFile creates a file (and its parent folders) with the given content.
func WriteFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

// ReadFile returns the content of a file, failing the test if it is missing.
func ReadFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}
