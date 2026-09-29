package ui

import (
	"bytes"
	"errors"
	"reflect"
	"strings"
	"testing"
)

func newTestUI(input string) (*UI, *bytes.Buffer) {
	var out bytes.Buffer
	SetColor(false)
	return New(strings.NewReader(input), &out, &out), &out
}

func TestParseSelection(t *testing.T) {
	tests := []struct {
		input   string
		max     int
		want    []int
		wantErr bool
	}{
		{"", 5, []int{}, false},
		{"1", 5, []int{0}, false},
		{"1,3", 5, []int{0, 2}, false},
		{" 2 , 4 ", 5, []int{1, 3}, false},
		{"2-4", 5, []int{1, 2, 3}, false},
		{"4-2", 5, []int{1, 2, 3}, false},
		{"1,1,1-2", 5, []int{0, 1}, false},
		{"0", 5, nil, true},
		{"6", 5, nil, true},
		{"a", 5, nil, true},
		{"1-x", 5, nil, true},
	}
	for _, tt := range tests {
		got, err := ParseSelection(tt.input, tt.max)
		if (err != nil) != tt.wantErr {
			t.Errorf("ParseSelection(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			continue
		}
		if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
			t.Errorf("ParseSelection(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestAsk(t *testing.T) {
	u, _ := newTestUI("hello\n\n")
	if got, _ := u.Ask("q", "def"); got != "hello" {
		t.Errorf("got %q, want hello", got)
	}
	if got, _ := u.Ask("q", "def"); got != "def" {
		t.Errorf("empty answer should give the default, got %q", got)
	}
	if _, err := u.Ask("q", ""); !errors.Is(err, ErrNoInput) {
		t.Errorf("closed input should give ErrNoInput, got %v", err)
	}
}

func TestAskRequiredRetries(t *testing.T) {
	u, out := newTestUI("\n  \nvalue\n")
	got, err := u.AskRequired("q", "")
	if err != nil || got != "value" {
		t.Fatalf("got %q, %v", got, err)
	}
	if strings.Count(out.String(), "A value is required") != 2 {
		t.Errorf("expected two warnings, output:\n%s", out)
	}
}

func TestConfirm(t *testing.T) {
	u, _ := newTestUI("y\nNO\nmaybe\nyes\n\n")
	want := []bool{true, false, true, false}
	for i, w := range want {
		got, err := u.Confirm("ok?", false)
		if err != nil || got != w {
			t.Errorf("answer %d: got %v, %v; want %v", i, got, err, w)
		}
	}
}

func TestAssumeYes(t *testing.T) {
	u, _ := newTestUI("") // no input at all
	u.AssumeYes = true
	if ok, err := u.Confirm("ok?", false); !ok || err != nil {
		t.Errorf("Confirm with AssumeYes = %v, %v", ok, err)
	}
	if got, err := u.Ask("q", "def"); got != "def" || err != nil {
		t.Errorf("Ask with AssumeYes = %q, %v", got, err)
	}
}

func TestChoose(t *testing.T) {
	u, _ := newTestUI("9\nx\n2\n")
	got, err := u.Choose("pick", []string{"a", "b", "c"})
	if err != nil || got != 1 {
		t.Errorf("Choose = %d, %v; want 1", got, err)
	}
}

func TestColorDisabled(t *testing.T) {
	SetColor(false)
	if Red("x") != "x" {
		t.Error("colors should be off")
	}
	SetColor(true)
	defer SetColor(false)
	if Red("x") == "x" || Red("") != "" {
		t.Error("colors should be on, and empty strings stay empty")
	}
}

func TestIndent(t *testing.T) {
	if got := Indent("a\n\nb\n", 2); got != "  a\n\n  b" {
		t.Errorf("Indent = %q", got)
	}
}
