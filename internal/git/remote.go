package git

import (
	"fmt"
	"net/url"
	"strings"
)

// Remotes lists the configured remote names.
func (c *Client) Remotes() ([]string, error) {
	return c.Lines("remote")
}

// DefaultRemote returns "origin" when it exists, otherwise the first remote,
// or "" when the repository has no remotes.
func (c *Client) DefaultRemote() string {
	remotes, _ := c.Remotes()
	for _, r := range remotes {
		if r == "origin" {
			return r
		}
	}
	if len(remotes) > 0 {
		return remotes[0]
	}
	return ""
}

// RemoteURL returns the fetch URL of a remote.
func (c *Client) RemoteURL(remote string) (string, error) {
	return c.Output("remote", "get-url", remote)
}

// AddRemote adds a new remote.
func (c *Client) AddRemote(name, rawURL string) error {
	_, err := c.Output("remote", "add", name, rawURL)
	return err
}

// RemoteBranchExists asks the remote whether it has the given branch.
func (c *Client) RemoteBranchExists(remote, branch string) (bool, error) {
	_, err := c.Output("ls-remote", "--exit-code", "--heads", remote, "refs/heads/"+branch)
	switch exitCode(err) {
	case 0:
		return true, nil
	case 2: // --exit-code: the remote answered but has no matching ref
		return false, nil
	default:
		return false, fmt.Errorf("could not reach remote %q: %w", remote, err)
	}
}

// RemoteDefaultBranch asks the remote which branch its HEAD points at.
// It returns "" for an empty remote repository.
func (c *Client) RemoteDefaultBranch(remote string) (string, error) {
	out, err := c.Output("ls-remote", "--symref", remote, "HEAD")
	if err != nil {
		return "", fmt.Errorf("could not reach remote %q: %w", remote, err)
	}
	for _, line := range strings.Split(out, "\n") {
		// Looks like: "ref: refs/heads/main	HEAD"
		if ref, ok := strings.CutPrefix(line, "ref: refs/heads/"); ok {
			if fields := strings.Fields(ref); len(fields) > 0 {
				return fields[0], nil
			}
		}
	}
	return "", nil
}

// CleanURL tidies a repository URL typed or pasted by a user.
func CleanURL(raw string) string {
	cleaned := strings.TrimSpace(raw)
	cleaned = strings.TrimRight(cleaned, "?/")
	return cleaned
}

// WebURL converts a remote URL (HTTPS, SSH or scp-like) into the address
// of the repository's web page, e.g.
//
//	git@github.com:user/repo.git  ->  https://github.com/user/repo
func WebURL(remote string) (string, error) {
	remote = CleanURL(remote)
	var host, path string

	switch {
	case strings.Contains(remote, "://"):
		u, err := url.Parse(remote)
		if err != nil || u.Host == "" {
			return "", fmt.Errorf("cannot understand remote URL %q", remote)
		}
		host, path = u.Hostname(), u.Path
		if u.Scheme == "http" || u.Scheme == "https" {
			host = u.Host // keep an explicit web port
		}
	case strings.Contains(remote, ":") && !isWindowsPath(remote):
		// scp-like syntax: [user@]host:path
		hostPart, pathPart, _ := strings.Cut(remote, ":")
		if at := strings.LastIndex(hostPart, "@"); at >= 0 {
			hostPart = hostPart[at+1:]
		}
		host, path = hostPart, pathPart
	default:
		return "", fmt.Errorf("remote %q is a local path, not a web address", remote)
	}

	path = strings.TrimSuffix(strings.Trim(path, "/"), ".git")
	if host == "" || path == "" {
		return "", fmt.Errorf("cannot understand remote URL %q", remote)
	}
	return "https://" + host + "/" + path, nil
}

// isWindowsPath reports whether s starts with a drive letter, like C:\ or C:/.
func isWindowsPath(s string) bool {
	return len(s) >= 3 && s[1] == ':' && (s[2] == '\\' || s[2] == '/')
}

// GitHubRepo extracts the owner and repository name from a GitHub remote URL.
func GitHubRepo(remote string) (owner, repo string, ok bool) {
	web, err := WebURL(remote)
	if err != nil {
		return "", "", false
	}
	rest, found := strings.CutPrefix(web, "https://github.com/")
	if !found {
		return "", "", false
	}
	owner, repo, found = strings.Cut(rest, "/")
	if !found || owner == "" || repo == "" || strings.Contains(repo, "/") {
		return "", "", false
	}
	return owner, repo, true
}

// ExpandCloneURL turns the GitHub shorthand "owner/repo" into a full clone
// URL. Anything else is returned unchanged.
func ExpandCloneURL(target string) string {
	target = CleanURL(target)
	if strings.Contains(target, ":") || strings.HasPrefix(target, ".") || strings.Count(target, "/") != 1 {
		return target
	}
	owner, repo, _ := strings.Cut(target, "/")
	if owner == "" || repo == "" {
		return target
	}
	return "https://github.com/" + owner + "/" + strings.TrimSuffix(repo, ".git") + ".git"
}
