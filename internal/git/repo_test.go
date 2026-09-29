package git

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/ettisafxrup/rgit/internal/testenv"
)

func quietClient(dir string) *Client {
	return &Client{Dir: dir, Stdout: io.Discard, Stderr: io.Discard}
}

func TestRepositoryBasics(t *testing.T) {
	testenv.Isolate(t)
	dir := t.TempDir()
	c := quietClient(dir)

	if c.IsRepo() {
		t.Fatal("empty temp dir should not be a repository")
	}
	if err := c.Init("main"); err != nil {
		t.Fatal(err)
	}
	if !c.IsRepo() {
		t.Fatal("expected a repository after Init")
	}
	if branch, err := c.CurrentBranch(); err != nil || branch != "main" {
		t.Errorf("CurrentBranch = %q, %v", branch, err)
	}
	if c.HasCommits() {
		t.Error("new repository should have no commits")
	}

	testenv.WriteFile(t, filepath.Join(dir, "a.txt"), "hello")
	if s, _ := c.Status(); s.Untracked != 1 {
		t.Errorf("expected one untracked file, got %+v", s)
	}
	if c.HasStagedChanges() {
		t.Error("nothing is staged yet")
	}
	c.Output("add", "a.txt")
	if !c.HasStagedChanges() {
		t.Error("a.txt should be staged")
	}
	if _, err := c.Output("commit", "-m", "first"); err != nil {
		t.Fatal(err)
	}
	if last, err := c.LastCommit(); err != nil || last.Subject != "first" {
		t.Errorf("LastCommit = %+v, %v", last, err)
	}
	if got := c.DefaultBranch(); got != "main" {
		t.Errorf("DefaultBranch = %q", got)
	}
}

func TestRemoteQueries(t *testing.T) {
	testenv.Isolate(t)
	remoteDir := testenv.BareRemote(t, "trunk")
	dir := t.TempDir()
	c := quietClient(dir)
	c.Init("trunk")
	c.AddRemote("origin", remoteDir)

	if got := c.DefaultRemote(); got != "origin" {
		t.Errorf("DefaultRemote = %q", got)
	}
	if def, err := c.RemoteDefaultBranch("origin"); err != nil || def != "" {
		t.Errorf("empty remote: RemoteDefaultBranch = %q, %v", def, err)
	}

	testenv.WriteFile(t, filepath.Join(dir, "a.txt"), "hello")
	c.Output("add", ".")
	c.Output("commit", "-m", "first")
	if err := c.Run("push", "-u", "origin", "trunk"); err != nil {
		t.Fatal(err)
	}

	if ok, err := c.RemoteBranchExists("origin", "trunk"); !ok || err != nil {
		t.Errorf("RemoteBranchExists(trunk) = %v, %v", ok, err)
	}
	if ok, err := c.RemoteBranchExists("origin", "nope"); ok || err != nil {
		t.Errorf("RemoteBranchExists(nope) = %v, %v", ok, err)
	}
	if _, err := c.RemoteBranchExists("missing-remote", "trunk"); err == nil {
		t.Error("expected an error for a missing remote")
	}
	if def, err := c.RemoteDefaultBranch("origin"); err != nil || def != "trunk" {
		t.Errorf("RemoteDefaultBranch = %q, %v", def, err)
	}
	if up, ok := c.Upstream(); !ok || up != "origin/trunk" {
		t.Errorf("Upstream = %q, %v", up, ok)
	}

	testenv.WriteFile(t, filepath.Join(dir, "b.txt"), "more")
	c.Output("add", ".")
	c.Output("commit", "-m", "second")
	if ahead, behind, err := c.AheadBehind(); ahead != 1 || behind != 0 || err != nil {
		t.Errorf("AheadBehind = %d, %d, %v", ahead, behind, err)
	}
}

func TestErrorMessage(t *testing.T) {
	testenv.Isolate(t)
	_, err := quietClient(os.TempDir()).Output("definitely-not-a-git-command")
	gitErr, ok := err.(*Error)
	if !ok || gitErr.ExitCode == 0 || gitErr.Stderr == "" {
		t.Errorf("expected a descriptive *Error, got %#v", err)
	}
}
