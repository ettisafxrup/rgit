package git

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// ErrDetachedHead is returned when HEAD does not point at a branch.
var ErrDetachedHead = errors.New("you are in 'detached HEAD' state; switch to a branch first (rgit switch)")

// IsRepo reports whether the working directory is inside a git work tree.
func (c *Client) IsRepo() bool {
	out, err := c.Output("rev-parse", "--is-inside-work-tree")
	return err == nil && out == "true"
}

// Init creates a new repository whose first branch is named branch. This
// works on every git version, unlike "git init -b".
func (c *Client) Init(branch string) error {
	if _, err := c.Output("init"); err != nil {
		return err
	}
	_, err := c.Output("symbolic-ref", "HEAD", "refs/heads/"+branch)
	return err
}

// CurrentBranch returns the name of the checked-out branch.
func (c *Client) CurrentBranch() (string, error) {
	branch, err := c.Output("symbolic-ref", "--quiet", "--short", "HEAD")
	if err != nil {
		return "", ErrDetachedHead
	}
	return branch, nil
}

// HasCommits reports whether the repository has at least one commit.
func (c *Client) HasCommits() bool {
	return c.Succeeds("rev-parse", "--verify", "--quiet", "HEAD")
}

// Branches lists local branch names.
func (c *Client) Branches() ([]string, error) {
	return c.Lines("for-each-ref", "--format=%(refname:short)", "refs/heads")
}

// BranchExists reports whether a local branch exists.
func (c *Client) BranchExists(name string) bool {
	return c.Succeeds("show-ref", "--verify", "--quiet", "refs/heads/"+name)
}

// MergedBranches lists local branches fully merged into base.
func (c *Client) MergedBranches(base string) ([]string, error) {
	return c.Lines("for-each-ref", "--merged="+base, "--format=%(refname:short)", "refs/heads")
}

// Upstream returns the upstream of the current branch, e.g. "origin/main".
func (c *Client) Upstream() (string, bool) {
	out, err := c.Output("rev-parse", "--abbrev-ref", "--symbolic-full-name", "@{u}")
	return out, err == nil && out != ""
}

// AheadBehind counts commits the current branch is ahead of and behind its upstream.
func (c *Client) AheadBehind() (ahead, behind int, err error) {
	out, err := c.Output("rev-list", "--left-right", "--count", "HEAD...@{u}")
	if err != nil {
		return 0, 0, err
	}
	fields := strings.Fields(out)
	if len(fields) != 2 {
		return 0, 0, fmt.Errorf("unexpected rev-list output %q", out)
	}
	ahead, _ = strconv.Atoi(fields[0])
	behind, _ = strconv.Atoi(fields[1])
	return ahead, behind, nil
}

// HasStagedChanges reports whether anything is staged for commit.
func (c *Client) HasStagedChanges() bool {
	// "diff --cached --quiet" exits with 1 when there are differences.
	_, err := c.Output("diff", "--cached", "--quiet")
	return exitCode(err) == 1
}

// DefaultBranch guesses the repository's main line of development: the
// remote's HEAD if known, otherwise "main" or "master", otherwise the
// current branch.
func (c *Client) DefaultBranch() string {
	if remote := c.DefaultRemote(); remote != "" {
		ref, err := c.Output("symbolic-ref", "--quiet", "--short", "refs/remotes/"+remote+"/HEAD")
		if err == nil && ref != "" {
			return strings.TrimPrefix(ref, remote+"/")
		}
	}
	for _, name := range []string{"main", "master"} {
		if c.BranchExists(name) {
			return name
		}
	}
	branch, _ := c.CurrentBranch()
	return branch
}

// Commit is a short summary of a single commit.
type Commit struct {
	Hash    string
	Subject string
	When    string
}

// LastCommit describes the commit HEAD points at.
func (c *Client) LastCommit() (Commit, error) {
	out, err := c.Output("log", "-1", "--format=%h%x00%s%x00%cr")
	if err != nil {
		return Commit{}, err
	}
	parts := strings.SplitN(out, "\x00", 3)
	if len(parts) != 3 {
		return Commit{}, fmt.Errorf("unexpected log output %q", out)
	}
	return Commit{Hash: parts[0], Subject: parts[1], When: parts[2]}, nil
}

// StashCount returns the number of stash entries.
func (c *Client) StashCount() int {
	lines, _ := c.Lines("stash", "list")
	return len(lines)
}

// Config reads a configuration value, returning "" when it is not set.
func (c *Client) Config(key string) string {
	out, _ := c.Output("config", "--get", key)
	return out
}

// GlobalConfig reads a value from the user's global configuration.
func (c *Client) GlobalConfig(key string) string {
	out, _ := c.Output("config", "--global", "--get", key)
	return out
}

// SetGlobalConfig writes a value to the user's global configuration.
func (c *Client) SetGlobalConfig(key, value string) error {
	_, err := c.Output("config", "--global", key, value)
	return err
}
