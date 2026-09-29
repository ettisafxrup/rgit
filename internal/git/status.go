package git

import "strings"

// Change is one entry of "git status --porcelain".
type Change struct {
	Index    byte // state in the staging area
	WorkTree byte // state in the working tree
	Path     string
}

// Status summarizes the working tree.
type Status struct {
	Changes   []Change
	Staged    int
	Modified  int
	Untracked int
	Conflicts int
}

// Clean reports whether there is nothing to commit.
func (s Status) Clean() bool { return len(s.Changes) == 0 }

// Status reads the state of the working tree.
func (c *Client) Status() (Status, error) {
	// Raw output: a leading space is part of the first entry (" M file").
	out, err := c.RawOutput("status", "--porcelain", "-z")
	if err != nil {
		return Status{}, err
	}
	return parseStatus(out), nil
}

// parseStatus parses NUL-separated porcelain v1 output.
func parseStatus(out string) Status {
	var s Status
	entries := strings.Split(out, "\x00")
	for i := 0; i < len(entries); i++ {
		entry := entries[i]
		if len(entry) < 4 {
			continue
		}
		change := Change{Index: entry[0], WorkTree: entry[1], Path: entry[3:]}
		if change.Index == 'R' || change.Index == 'C' {
			i++ // renames and copies are followed by the original path
		}
		s.Changes = append(s.Changes, change)

		switch {
		case change.Index == '?':
			s.Untracked++
		case isConflict(change):
			s.Conflicts++
		default:
			if change.Index != ' ' {
				s.Staged++
			}
			if change.WorkTree != ' ' {
				s.Modified++
			}
		}
	}
	return s
}

func isConflict(c Change) bool {
	pair := string([]byte{c.Index, c.WorkTree})
	return c.Index == 'U' || c.WorkTree == 'U' || pair == "AA" || pair == "DD"
}

// Label describes a change in a single human-friendly word.
func (c Change) Label() string {
	if c.Index == '?' {
		return "new"
	}
	if isConflict(c) {
		return "conflict"
	}
	state := c.Index
	if state == ' ' {
		state = c.WorkTree
	}
	switch state {
	case 'A':
		return "added"
	case 'D':
		return "deleted"
	case 'R':
		return "renamed"
	case 'C':
		return "copied"
	case 'T':
		return "typechange"
	default:
		return "modified"
	}
}
