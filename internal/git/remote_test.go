package git

import "testing"

func TestWebURL(t *testing.T) {
	tests := map[string]string{
		"https://github.com/user/repo.git":          "https://github.com/user/repo",
		"https://github.com/user/repo":              "https://github.com/user/repo",
		"https://token@github.com/user/repo.git":    "https://github.com/user/repo",
		"https://git.example.com:8443/team/app.git": "https://git.example.com:8443/team/app",
		"git@github.com:user/repo.git":              "https://github.com/user/repo",
		"git@gitlab.com:group/sub/project.git":      "https://gitlab.com/group/sub/project",
		"ssh://git@github.com/user/repo.git":        "https://github.com/user/repo",
		"ssh://git@host.com:2222/user/repo.git":     "https://host.com/user/repo",
		"  https://github.com/user/repo.git?  ":     "https://github.com/user/repo",
	}
	for in, want := range tests {
		got, err := WebURL(in)
		if err != nil || got != want {
			t.Errorf("WebURL(%q) = %q, %v; want %q", in, got, err, want)
		}
	}

	for _, bad := range []string{"", "/local/path/repo", "../repo", `C:\repos\app.git`, "D:/repos/app"} {
		if got, err := WebURL(bad); err == nil {
			t.Errorf("WebURL(%q) = %q, want an error", bad, got)
		}
	}
}

func TestGitHubRepo(t *testing.T) {
	owner, repo, ok := GitHubRepo("git@github.com:ettisafxrup/rgit.git")
	if !ok || owner != "ettisafxrup" || repo != "rgit" {
		t.Errorf("got %q %q %v", owner, repo, ok)
	}
	if _, _, ok := GitHubRepo("https://gitlab.com/a/b.git"); ok {
		t.Error("gitlab URL must not be treated as GitHub")
	}
}

func TestExpandCloneURL(t *testing.T) {
	tests := map[string]string{
		"user/repo":                    "https://github.com/user/repo.git",
		"user/repo.git":                "https://github.com/user/repo.git",
		"https://github.com/user/repo": "https://github.com/user/repo",
		"git@github.com:user/repo.git": "git@github.com:user/repo.git",
		"./local/repo":                 "./local/repo",
		"a/b/c":                        "a/b/c",
		"C:/repos/project":             "C:/repos/project",
	}
	for in, want := range tests {
		if got := ExpandCloneURL(in); got != want {
			t.Errorf("ExpandCloneURL(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestCleanURL(t *testing.T) {
	if got := CleanURL("  https://x.com/a/b.git?\t"); got != "https://x.com/a/b.git" {
		t.Errorf("CleanURL = %q", got)
	}
}
