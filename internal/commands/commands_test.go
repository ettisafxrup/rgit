package commands

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ettisafxrup/rgit/internal/cli"
	"github.com/ettisafxrup/rgit/internal/git"
	"github.com/ettisafxrup/rgit/internal/testenv"
	"github.com/ettisafxrup/rgit/internal/ui"
)

// rgit runs an rgit command line inside dir, feeding input to its prompts.
func rgit(t *testing.T, dir, input string, args ...string) (string, error) {
	t.Helper()
	t.Chdir(dir)
	ui.SetColor(false)
	var out bytes.Buffer
	u := ui.New(strings.NewReader(input), &out, &out)
	app := cli.NewApp(u, &git.Client{Stdout: &out, Stderr: &out}, All()...)
	err := app.Run(args)
	return out.String(), err
}

// mustRgit is rgit that fails the test on error.
func mustRgit(t *testing.T, dir, input string, args ...string) string {
	t.Helper()
	out, err := rgit(t, dir, input, args...)
	if err != nil {
		t.Fatalf("rgit %v failed: %v\n%s", args, err, out)
	}
	return out
}

// projectWithRemote returns a folder with a first commit already pushed to
// a bare remote on branch main.
func projectWithRemote(t *testing.T) (dir, remote string) {
	t.Helper()
	testenv.Isolate(t)
	remote = testenv.BareRemote(t, "main")
	dir = t.TempDir()
	testenv.WriteFile(t, filepath.Join(dir, "app.txt"), "v1")
	mustRgit(t, dir, "\n\n"+remote+"\n\n", "push", "Initial", "commit")
	return dir, remote
}

func remoteLog(t *testing.T, remote, branch string) string {
	return testenv.Git(t, remote, "log", branch, "--format=%s")
}

func TestPushCreatesRepositoryAndPushes(t *testing.T) {
	dir, remote := projectWithRemote(t)

	if got := remoteLog(t, remote, "main"); strings.TrimSpace(got) != "Initial commit" {
		t.Errorf("remote history = %q", got)
	}
	if up := testenv.Git(t, dir, "rev-parse", "--abbrev-ref", "@{u}"); strings.TrimSpace(up) != "origin/main" {
		t.Errorf("upstream = %q", up)
	}
}

func TestPushNothingToCommit(t *testing.T) {
	dir, _ := projectWithRemote(t)
	out := mustRgit(t, dir, "", "push")
	if !strings.Contains(out, "Nothing new to commit") {
		t.Errorf("unexpected output:\n%s", out)
	}
}

func TestPushUsesDefaultMessage(t *testing.T) {
	dir, remote := projectWithRemote(t)
	testenv.WriteFile(t, filepath.Join(dir, "app.txt"), "v2")
	mustRgit(t, dir, "", "push", "-y")
	if got := remoteLog(t, remote, "main"); !strings.HasPrefix(got, "Update ") {
		t.Errorf("expected the default 'Update <date>' message, got %q", got)
	}
}

func TestPushDeclinedStagingAborts(t *testing.T) {
	dir, _ := projectWithRemote(t)
	testenv.WriteFile(t, filepath.Join(dir, "app.txt"), "v2")
	_, err := rgit(t, dir, "n\n", "push", "msg")
	if !errors.Is(err, cli.ErrAborted) {
		t.Errorf("expected ErrAborted, got %v", err)
	}
}

func TestPushMatchesExistingRemoteBranch(t *testing.T) {
	testenv.Isolate(t)
	// The remote already has history on a branch called "trunk".
	remote := testenv.BareRemote(t, "trunk")
	seed := t.TempDir()
	testenv.Git(t, seed, "init")
	testenv.Git(t, seed, "symbolic-ref", "HEAD", "refs/heads/trunk")
	testenv.WriteFile(t, filepath.Join(seed, "README.md"), "hello")
	testenv.Git(t, seed, "add", ".")
	testenv.Git(t, seed, "commit", "-m", "Created on the server")
	testenv.Git(t, seed, "push", remote, "trunk")

	dir := t.TempDir()
	testenv.WriteFile(t, filepath.Join(dir, "main.go"), "package main")
	mustRgit(t, dir, remote+"\n", "push", "-y", "Local work")

	got := strings.Fields(remoteLog(t, remote, "trunk"))
	if strings.Join(got, " ") != "Local work Created on the server" {
		t.Errorf("remote trunk history = %v", got)
	}
	if _, err := os.Stat(filepath.Join(dir, "README.md")); err != nil {
		t.Error("the server's README should have been pulled in")
	}
}

func TestPullWithUncommittedChanges(t *testing.T) {
	dir, remote := projectWithRemote(t)

	// Someone else pushes a new file.
	other := filepath.Join(t.TempDir(), "other")
	testenv.Git(t, "", "clone", remote, other)
	testenv.WriteFile(t, filepath.Join(other, "theirs.txt"), "from a teammate")
	mustRgit(t, other, "", "push", "-y", "Teammate change")

	// We have local, uncommitted work; pull must keep it.
	testenv.WriteFile(t, filepath.Join(dir, "app.txt"), "local edit")
	mustRgit(t, dir, "", "pull")

	if testenv.ReadFile(t, filepath.Join(dir, "theirs.txt")) != "from a teammate" {
		t.Error("remote change was not pulled")
	}
	if testenv.ReadFile(t, filepath.Join(dir, "app.txt")) != "local edit" {
		t.Error("local uncommitted work was lost")
	}
}

func TestPullIntoEmptyFolder(t *testing.T) {
	_, remote := projectWithRemote(t)
	fresh := t.TempDir()
	mustRgit(t, fresh, "\n\n"+remote+"\n", "pull")
	if testenv.ReadFile(t, filepath.Join(fresh, "app.txt")) != "v1" {
		t.Error("expected the remote files in the new folder")
	}
}

func TestSyncDivergedBranches(t *testing.T) {
	dir, remote := projectWithRemote(t)
	other := filepath.Join(t.TempDir(), "other")
	testenv.Git(t, "", "clone", remote, other)
	testenv.WriteFile(t, filepath.Join(other, "b.txt"), "b")
	mustRgit(t, other, "", "push", "-y", "Their commit")

	testenv.WriteFile(t, filepath.Join(dir, "a.txt"), "a")
	testenv.Git(t, dir, "add", ".")
	testenv.Git(t, dir, "commit", "-m", "Our commit")
	mustRgit(t, dir, "", "sync")

	if got := strings.Fields(remoteLog(t, remote, "main")); len(got) != 6 || got[0] != "Our" {
		t.Errorf("remote history after sync = %v", got)
	}
}

func TestStatusAndLog(t *testing.T) {
	dir, _ := projectWithRemote(t)
	testenv.WriteFile(t, filepath.Join(dir, "app.txt"), "changed")
	testenv.WriteFile(t, filepath.Join(dir, "new.txt"), "x")

	out := mustRgit(t, dir, "", "status")
	// "modified   app.txt" also guards against the first letter of an
	// unstaged path being lost (porcelain lines start with a space).
	for _, want := range []string{"main", "origin/main", "up to date", "Initial commit", "1 modified", "1 new", "modified   app.txt", "new.txt"} {
		if !strings.Contains(out, want) {
			t.Errorf("status output is missing %q:\n%s", want, out)
		}
	}
	if out := mustRgit(t, dir, "", "log", "-n", "5"); !strings.Contains(out, "Initial commit") {
		t.Errorf("log output:\n%s", out)
	}
}

func TestUndo(t *testing.T) {
	testenv.Isolate(t)
	dir := t.TempDir()
	testenv.Git(t, dir, "init")
	testenv.WriteFile(t, filepath.Join(dir, "a.txt"), "a")
	testenv.Git(t, dir, "add", ".")
	testenv.Git(t, dir, "commit", "-m", "first")
	testenv.WriteFile(t, filepath.Join(dir, "b.txt"), "b")
	testenv.Git(t, dir, "add", ".")
	testenv.Git(t, dir, "commit", "-m", "second")

	mustRgit(t, dir, "\n", "undo")
	if got := strings.TrimSpace(testenv.Git(t, dir, "log", "--format=%s")); got != "first" {
		t.Errorf("history after undo = %q", got)
	}
	if staged := testenv.Git(t, dir, "diff", "--cached", "--name-only"); strings.TrimSpace(staged) != "b.txt" {
		t.Errorf("b.txt should still be staged, got %q", staged)
	}

	// Undoing the very first commit leaves an unborn branch with staged files.
	testenv.Git(t, dir, "commit", "-m", "second again")
	mustRgit(t, dir, "", "undo", "-y")
	mustRgit(t, dir, "", "undo", "-y")
	if testenv.Git(t, dir, "status", "--porcelain") == "" {
		t.Error("files should remain after undoing the first commit")
	}
	if _, err := rgit(t, dir, "", "undo", "-y"); err == nil {
		t.Error("undo without commits should fail")
	}
}

func TestNewSwitchRenameCleanup(t *testing.T) {
	dir, remote := projectWithRemote(t)

	mustRgit(t, dir, "", "new", "feature/one")
	testenv.WriteFile(t, filepath.Join(dir, "one.txt"), "1")
	mustRgit(t, dir, "", "push", "-y", "Feature one")

	// Rename on the remote too: the new name appears, the old one is gone.
	mustRgit(t, dir, "", "rename", "feature/one:feature/first", "-y")
	branches := testenv.Git(t, remote, "branch", "--format=%(refname:short)")
	if !strings.Contains(branches, "feature/first") || strings.Contains(branches, "feature/one") {
		t.Errorf("remote branches after rename:\n%s", branches)
	}

	// Pick "main" from the switch menu (the only other branch).
	mustRgit(t, dir, "1\n", "switch")
	if b := strings.TrimSpace(testenv.Git(t, dir, "branch", "--show-current")); b != "main" {
		t.Fatalf("current branch = %q", b)
	}

	testenv.Git(t, dir, "merge", "--ff-only", "feature/first")
	mustRgit(t, dir, "", "new", "unmerged")
	testenv.WriteFile(t, filepath.Join(dir, "wip.txt"), "wip")
	testenv.Git(t, dir, "add", ".")
	testenv.Git(t, dir, "commit", "-m", "wip")
	mustRgit(t, dir, "", "switch", "main")

	out := mustRgit(t, dir, "", "cleanup", "-y")
	local := testenv.Git(t, dir, "branch", "--format=%(refname:short)")
	if strings.Contains(local, "feature/first") || !strings.Contains(local, "unmerged") || !strings.Contains(local, "main") {
		t.Errorf("branches after cleanup:\n%s\noutput:\n%s", local, out)
	}

	if _, err := rgit(t, dir, "", "new", "bad..name"); err == nil {
		t.Error("an invalid branch name should be rejected")
	}
}

func TestIgnore(t *testing.T) {
	testenv.Isolate(t)
	dir := t.TempDir()
	mustRgit(t, dir, "", "ignore", "go", "*.pem")
	mustRgit(t, dir, "", "ignore", "go", "*.pem") // nothing new the second time

	content := testenv.ReadFile(t, filepath.Join(dir, ".gitignore"))
	if strings.Count(content, "*.pem") != 1 || strings.Count(content, "go.work\n") != 1 {
		t.Errorf("duplicate or missing entries:\n%s", content)
	}
	if strings.Contains(content, ".gitignore\n") {
		t.Error(".gitignore must never ignore itself")
	}

	// Interactive: pick the second entry, then add a preset.
	os.Mkdir(filepath.Join(dir, "build"), 0o755)
	testenv.WriteFile(t, filepath.Join(dir, "notes.txt"), "")
	mustRgit(t, dir, "2\nmacos\n", "ignore")
	content = testenv.ReadFile(t, filepath.Join(dir, ".gitignore"))
	if !strings.Contains(content, "/notes.txt\n") || !strings.Contains(content, ".DS_Store") {
		t.Errorf("interactive picks missing:\n%s", content)
	}
}

func TestMergeIgnore(t *testing.T) {
	merged, added := mergeIgnore("node_modules/", [][]string{
		{"# Dependencies", "node_modules/", "# Logs", "*.log"},
	})
	want := "node_modules/\n\n# Logs\n*.log\n"
	if merged != want || added != 1 {
		t.Errorf("got %q (%d), want %q", merged, added, want)
	}
}

func TestLicense(t *testing.T) {
	testenv.Isolate(t)
	dir := t.TempDir()
	mustRgit(t, dir, "", "license", "mit", "--name", "Jane / Co", "--year", "2020")
	if text := testenv.ReadFile(t, filepath.Join(dir, "LICENSE")); !strings.Contains(text, "Copyright (c) 2020 Jane / Co") {
		t.Errorf("LICENSE:\n%s", text)
	}

	// An existing license is not overwritten without consent.
	if _, err := rgit(t, dir, "n\n", "license", "apache"); !errors.Is(err, cli.ErrAborted) {
		t.Errorf("expected ErrAborted, got %v", err)
	}
	if _, err := rgit(t, dir, "y\n", "license", "nonsense"); err == nil {
		t.Error("unknown license should fail")
	}

	// Interactive menu, default year, holder from git config.
	mustRgit(t, dir, "y\n3\n\n\n", "license")
	text := testenv.ReadFile(t, filepath.Join(dir, "LICENSE"))
	if !strings.Contains(text, "GNU GENERAL PUBLIC LICENSE") || !strings.Contains(text, "Test User") {
		t.Errorf("interactive LICENSE is wrong:\n%.300s", text)
	}
}

func TestReadmeAndArchi(t *testing.T) {
	testenv.Isolate(t)
	dir := t.TempDir()
	testenv.Git(t, dir, "init")
	testenv.Git(t, dir, "remote", "add", "origin", "git@github.com:jane/demo.git")
	testenv.WriteFile(t, filepath.Join(dir, "src", "main.go"), "package main")
	testenv.WriteFile(t, filepath.Join(dir, "dist", "app.exe"), "binary")
	testenv.WriteFile(t, filepath.Join(dir, ".gitignore"), "dist/\n")

	mustRgit(t, dir, "", "archi")
	archi := testenv.ReadFile(t, filepath.Join(dir, "ARCHITECTURE"))
	if !strings.Contains(archi, "main.go") || strings.Contains(archi, "dist") {
		t.Errorf("ARCHITECTURE should follow .gitignore:\n%s", archi)
	}

	// Title, description, run command, then "add a license?" -> mit, year, holder.
	mustRgit(t, dir, "Demo\nA demo app.\ngo run ./src\ny\n1\n\n\n", "readme")
	readme := testenv.ReadFile(t, filepath.Join(dir, "README.md"))
	for _, want := range []string{"# Demo", "A demo app.", "go run ./src", "https://github.com/jane/demo", "MIT License", "main.go"} {
		if !strings.Contains(readme, want) {
			t.Errorf("README is missing %q:\n%s", want, readme)
		}
	}
	if strings.Contains(readme, "README.md") {
		t.Error("the README should not list itself in the structure")
	}
}

func TestLogin(t *testing.T) {
	testenv.Isolate(t)
	dir := t.TempDir()
	mustRgit(t, dir, "Jane Doe\nnot-an-email\njane@example.com\n", "login")
	if name := strings.TrimSpace(testenv.Git(t, dir, "config", "--global", "user.name")); name != "Jane Doe" {
		t.Errorf("user.name = %q", name)
	}
	if out := mustRgit(t, dir, "", "whoami"); !strings.Contains(out, "jane@example.com") {
		t.Errorf("whoami output:\n%s", out)
	}
}

func TestCloneShorthandAndLocal(t *testing.T) {
	_, remote := projectWithRemote(t)
	parent := t.TempDir()
	mustRgit(t, parent, "", "clone", remote, "copy")
	if testenv.ReadFile(t, filepath.Join(parent, "copy", "app.txt")) != "v1" {
		t.Error("clone did not produce the files")
	}
}

func TestCommandsOutsideRepository(t *testing.T) {
	testenv.Isolate(t)
	dir := t.TempDir()
	for _, cmd := range []string{"status", "log", "undo", "sync", "cleanup", "open", "switch"} {
		if _, err := rgit(t, dir, "", cmd); !errors.Is(err, errNotRepo) {
			t.Errorf("rgit %s outside a repository: got %v", cmd, err)
		}
	}
}

func TestDoctorAndVersion(t *testing.T) {
	dir, _ := projectWithRemote(t)
	out, _ := rgit(t, dir, "", "doctor") // may report a missing credential helper
	for _, want := range []string{"installed", "Test User", "On branch main", "Remote origin"} {
		if !strings.Contains(out, want) {
			t.Errorf("doctor output is missing %q:\n%s", want, out)
		}
	}
	if out := mustRgit(t, dir, "", "version"); !strings.HasPrefix(out, "rgit ") {
		t.Errorf("version output: %q", out)
	}
}
