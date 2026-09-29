// Package git is a thin wrapper around the git command-line tool.
//
// rgit deliberately drives the real git binary instead of re-implementing
// it, so it behaves exactly like git does, including credential helpers,
// hooks and user configuration.
package git

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// ErrNotInstalled is returned when the git executable cannot be found.
var ErrNotInstalled = errors.New("git is not installed or not in your PATH (download it from https://git-scm.com)")

// Error describes a failed git invocation.
type Error struct {
	Args     []string
	ExitCode int
	Stderr   string
}

func (e *Error) Error() string {
	msg := fmt.Sprintf("git %s failed (exit code %d)", strings.Join(e.Args, " "), e.ExitCode)
	if e.Stderr != "" {
		msg += ": " + e.Stderr
	}
	return msg
}

// Client runs git commands inside a working directory.
type Client struct {
	Dir    string    // working directory; empty means the current one
	Stdout io.Writer // where Run streams git's output
	Stderr io.Writer
}

// New returns a client working in dir that streams output to the terminal.
func New(dir string) *Client {
	return &Client{Dir: dir, Stdout: os.Stdout, Stderr: os.Stderr}
}

// Installed reports whether the git executable is available.
func Installed() error {
	if _, err := exec.LookPath("git"); err != nil {
		return ErrNotInstalled
	}
	return nil
}

// Run executes git and lets its output flow to the user, for long-running
// commands such as push and pull where progress output matters.
func (c *Client) Run(args ...string) error {
	cmd := c.command(args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = c.Stdout
	cmd.Stderr = c.Stderr
	return wrapError(args, cmd.Run(), "")
}

// Output executes git and returns its trimmed standard output.
func (c *Client) Output(args ...string) (string, error) {
	out, err := c.RawOutput(args...)
	return strings.TrimSpace(out), err
}

// RawOutput executes git and returns its standard output untouched, for
// formats where leading whitespace is meaningful (such as status --porcelain).
func (c *Client) RawOutput(args ...string) (string, error) {
	var stdout, stderr bytes.Buffer
	cmd := c.command(args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), wrapError(args, err, strings.TrimSpace(stderr.String()))
}

// Succeeds runs git quietly and reports whether it exited with status 0.
func (c *Client) Succeeds(args ...string) bool {
	_, err := c.Output(args...)
	return err == nil
}

// Lines runs git and splits its output into non-empty lines.
func (c *Client) Lines(args ...string) ([]string, error) {
	out, err := c.Output(args...)
	if err != nil || out == "" {
		return nil, err
	}
	return strings.Split(strings.ReplaceAll(out, "\r\n", "\n"), "\n"), nil
}

func (c *Client) command(args ...string) *exec.Cmd {
	// core.quotepath=false keeps non-ASCII file names readable.
	full := append([]string{"-c", "core.quotepath=false"}, args...)
	cmd := exec.Command("git", full...)
	cmd.Dir = c.Dir
	return cmd
}

func wrapError(args []string, err error, stderr string) error {
	if err == nil {
		return nil
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return &Error{Args: args, ExitCode: exitErr.ExitCode(), Stderr: stderr}
	}
	if errors.Is(err, exec.ErrNotFound) {
		return ErrNotInstalled
	}
	return err
}

// exitCode returns the exit status of a git run: 0 on success, the
// status carried by err, or -1 when git could not run at all.
func exitCode(err error) int {
	if err == nil {
		return 0
	}
	var gitErr *Error
	if errors.As(err, &gitErr) {
		return gitErr.ExitCode
	}
	return -1
}
