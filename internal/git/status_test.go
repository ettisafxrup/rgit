package git

import "testing"

func TestParseStatus(t *testing.T) {
	out := "M  staged.go\x00 M modified.go\x00MM both.go\x00?? new.txt\x00" +
		"R  renamed.go\x00original.go\x00UU conflict.go\x00D  gone.go\x00"
	s := parseStatus(out)

	if len(s.Changes) != 7 {
		t.Fatalf("got %d changes, want 7: %+v", len(s.Changes), s.Changes)
	}
	if s.Staged != 4 || s.Modified != 2 || s.Untracked != 1 || s.Conflicts != 1 {
		t.Errorf("counts: staged=%d modified=%d untracked=%d conflicts=%d",
			s.Staged, s.Modified, s.Untracked, s.Conflicts)
	}

	wantLabels := []string{"modified", "modified", "modified", "new", "renamed", "conflict", "deleted"}
	for i, c := range s.Changes {
		if c.Label() != wantLabels[i] {
			t.Errorf("change %d (%s): label %q, want %q", i, c.Path, c.Label(), wantLabels[i])
		}
	}
	if s.Changes[4].Path != "renamed.go" {
		t.Errorf("rename path = %q", s.Changes[4].Path)
	}
}

func TestParseStatusClean(t *testing.T) {
	if s := parseStatus(""); !s.Clean() {
		t.Error("empty output should be clean")
	}
}
