package cli

import (
	"bytes"
	"errors"
	"flag"
	"reflect"
	"strings"
	"testing"

	"github.com/ettisafxrup/rgit/internal/ui"
)

func TestParseInterspersed(t *testing.T) {
	tests := []struct {
		args []string
		want []string
		yes  bool
	}{
		{[]string{"fix", "bug", "-y"}, []string{"fix", "bug"}, true},
		{[]string{"-y", "fix"}, []string{"fix"}, true},
		{[]string{"fix", "--", "-y"}, []string{"fix", "-y"}, false},
		{[]string{}, nil, false},
	}
	for _, tt := range tests {
		fs := newFlagSet("test")
		yes := fs.Bool("y", false, "")
		got, err := parseInterspersed(fs, tt.args)
		if err != nil {
			t.Fatalf("%v: %v", tt.args, err)
		}
		if !reflect.DeepEqual(got, tt.want) || *yes != tt.yes {
			t.Errorf("%v: got %v yes=%v, want %v yes=%v", tt.args, got, *yes, tt.want, tt.yes)
		}
	}
}

func TestSplitCommand(t *testing.T) {
	name, rest := splitCommand([]string{"--no-color", "push", "-y", "msg"})
	if name != "push" || !reflect.DeepEqual(rest, []string{"--no-color", "-y", "msg"}) {
		t.Errorf("got %q %v", name, rest)
	}
}

func TestDistance(t *testing.T) {
	tests := []struct {
		a, b string
		want int
	}{{"push", "push", 0}, {"psuh", "push", 2}, {"pul", "pull", 1}, {"", "abc", 3}}
	for _, tt := range tests {
		if got := distance(tt.a, tt.b); got != tt.want {
			t.Errorf("distance(%q, %q) = %d, want %d", tt.a, tt.b, got, tt.want)
		}
	}
}

func newTestApp(commands ...*Command) (*App, *bytes.Buffer) {
	var out bytes.Buffer
	ui.SetColor(false)
	return NewApp(ui.New(strings.NewReader(""), &out, &out), nil, commands...), &out
}

func TestRunDispatch(t *testing.T) {
	var gotArgs []string
	var gotFlag string
	cmd := &Command{
		Name: "greet", Aliases: []string{"hi"}, Group: "Everyday", Summary: "say hi", MaxArgs: 2,
		Flags: func(fs *flag.FlagSet) { fs.StringVar(&gotFlag, "name", "", "who to greet") },
		Run: func(ctx *Context) error {
			gotArgs = ctx.Args
			return nil
		},
	}
	app, out := newTestApp(cmd)

	if err := app.Run([]string{"hi", "a", "--name", "bob", "b"}); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(gotArgs, []string{"a", "b"}) || gotFlag != "bob" {
		t.Errorf("args=%v flag=%q", gotArgs, gotFlag)
	}

	var usageErr *UsageError
	if err := app.Run([]string{"greet", "1", "2", "3"}); !errors.As(err, &usageErr) {
		t.Errorf("too many args should be a usage error, got %v", err)
	}
	if err := app.Run([]string{"greet", "--bogus"}); !errors.As(err, &usageErr) {
		t.Errorf("unknown flag should be a usage error, got %v", err)
	}
	if err := app.Run([]string{"gret"}); err == nil || !strings.Contains(err.Error(), `did you mean "greet"`) {
		t.Errorf("expected a suggestion, got %v", err)
	}

	out.Reset()
	if err := app.Run([]string{"greet", "--help"}); err != nil || !strings.Contains(out.String(), "--name <string>") {
		t.Errorf("command help: %v\n%s", err, out)
	}
	out.Reset()
	if err := app.Run(nil); err != nil || !strings.Contains(out.String(), "greet") {
		t.Errorf("overview help: %v\n%s", err, out)
	}
	out.Reset()
	if err := app.Run([]string{"help", "hi"}); err != nil || !strings.Contains(out.String(), "rgit greet") {
		t.Errorf("help <alias>: %v\n%s", err, out)
	}
}

func TestUsageErrorGetsCommand(t *testing.T) {
	cmd := &Command{Name: "x", MaxArgs: -1, Run: func(*Context) error { return &UsageError{Message: "bad"} }}
	app, _ := newTestApp(cmd)
	var usageErr *UsageError
	if err := app.Run([]string{"x"}); !errors.As(err, &usageErr) || usageErr.Command != cmd {
		t.Errorf("usage error should carry its command, got %#v", err)
	}
}
