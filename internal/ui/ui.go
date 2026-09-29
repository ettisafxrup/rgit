// Package ui handles everything rgit shows to, and reads from, the user.
package ui

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// UI writes messages to the terminal and asks the user questions.
type UI struct {
	Out io.Writer
	Err io.Writer
	in  *bufio.Reader

	// AssumeYes is set by --yes: confirmations are answered "yes" and
	// questions with a default value take that value.
	AssumeYes bool
}

func New(in io.Reader, out, errOut io.Writer) *UI {
	return &UI{Out: out, Err: errOut, in: bufio.NewReader(in)}
}

func (u *UI) Println(a ...any)               { fmt.Fprintln(u.Out, a...) }
func (u *UI) Printf(format string, a ...any) { fmt.Fprintf(u.Out, format, a...) }

// Title prints a bold heading surrounded by blank lines.
func (u *UI) Title(text string) { fmt.Fprintf(u.Out, "\n %s\n\n", Bold(text)) }

// Message helpers, each with its own symbol: • info, › step, ✓ success, ! warning, ✗ error.
func (u *UI) Info(format string, a ...any)    { u.line(u.Out, Blue("•"), format, a...) }
func (u *UI) Step(format string, a ...any)    { u.line(u.Out, Cyan("›"), format, a...) }
func (u *UI) Success(format string, a ...any) { u.line(u.Out, Green("✓"), format, a...) }
func (u *UI) Warn(format string, a ...any)    { u.line(u.Err, Yellow("!"), format, a...) }
func (u *UI) Error(format string, a ...any)   { u.line(u.Err, Red("✗"), format, a...) }

// Hint prints a dimmed, indented suggestion, typically a command to run.
func (u *UI) Hint(format string, a ...any) {
	fmt.Fprintf(u.Err, "   %s\n", Dim(fmt.Sprintf(format, a...)))
}

// Field prints an aligned "label: value" pair.
func (u *UI) Field(label, value string) {
	fmt.Fprintf(u.Out, "   %-12s %s\n", Dim(label), value)
}

func (u *UI) line(w io.Writer, symbol, format string, a ...any) {
	fmt.Fprintf(w, " %s %s\n", Bold(symbol), fmt.Sprintf(format, a...))
}

// Indent prefixes every line of text with the given number of spaces.
func Indent(text string, spaces int) string {
	pad := strings.Repeat(" ", spaces)
	lines := strings.Split(strings.TrimRight(text, "\n"), "\n")
	for i, l := range lines {
		if l != "" {
			lines[i] = pad + l
		}
	}
	return strings.Join(lines, "\n")
}
