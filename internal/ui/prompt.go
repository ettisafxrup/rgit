package ui

import (
	"errors"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
)

// ErrNoInput is returned when a question is asked but stdin has no more data.
var ErrNoInput = errors.New("no input available (stdin closed); use --yes to accept defaults")

// Ask asks a free-text question. An empty answer returns def.
func (u *UI) Ask(question, def string) (string, error) {
	if u.AssumeYes && def != "" {
		u.printQuestion(question, def)
		fmt.Fprintln(u.Out, def)
		return def, nil
	}
	u.printQuestion(question, def)
	answer, err := u.readLine()
	if err != nil {
		return "", err
	}
	if answer == "" {
		return def, nil
	}
	return answer, nil
}

// AskRequired asks a question until a non-empty answer is given.
func (u *UI) AskRequired(question, def string) (string, error) {
	for {
		answer, err := u.Ask(question, def)
		if err != nil || answer != "" {
			return answer, err
		}
		u.Warn("A value is required.")
	}
}

// Confirm asks a yes/no question. An empty answer returns def.
func (u *UI) Confirm(question string, def bool) (bool, error) {
	hint := "y/N"
	if def {
		hint = "Y/n"
	}
	prompt := fmt.Sprintf(" %s %s %s ", Bold(Cyan("?")), question, Dim("["+hint+"]"))

	if u.AssumeYes {
		fmt.Fprintln(u.Out, prompt+"yes")
		return true, nil
	}
	for {
		fmt.Fprint(u.Out, prompt)
		answer, err := u.readLine()
		if err != nil {
			return false, err
		}
		switch strings.ToLower(answer) {
		case "":
			return def, nil
		case "y", "yes":
			return true, nil
		case "n", "no":
			return false, nil
		}
		u.Warn("Please answer yes or no.")
	}
}

// Choose shows a numbered list and returns the index of the picked option.
func (u *UI) Choose(question string, options []string) (int, error) {
	fmt.Fprintf(u.Out, " %s %s\n", Bold(Cyan("?")), question)
	u.printOptions(options)
	for {
		fmt.Fprintf(u.Out, "   %s ", Green(">"))
		answer, err := u.readLine()
		if err != nil {
			return 0, err
		}
		n, convErr := strconv.Atoi(answer)
		if convErr == nil && n >= 1 && n <= len(options) {
			return n - 1, nil
		}
		u.Warn("Enter a number between 1 and %d.", len(options))
	}
}

// ChooseMany shows a numbered list and returns the indexes the user picked,
// written as a list like "1,3,5-7". An empty answer picks nothing.
func (u *UI) ChooseMany(question string, options []string) ([]int, error) {
	fmt.Fprintf(u.Out, " %s %s %s\n", Bold(Cyan("?")), question, Dim("(e.g. 1,3,5-7 — empty to skip)"))
	u.printOptions(options)
	for {
		fmt.Fprintf(u.Out, "   %s ", Green(">"))
		answer, err := u.readLine()
		if err != nil {
			return nil, err
		}
		picked, parseErr := ParseSelection(answer, len(options))
		if parseErr == nil {
			return picked, nil
		}
		u.Warn("%v", parseErr)
	}
}

// ParseSelection turns input such as "1, 3, 5-7" into sorted zero-based
// indexes. Every number must be between 1 and max.
func ParseSelection(input string, max int) ([]int, error) {
	seen := map[int]bool{}
	for _, part := range strings.FieldsFunc(input, func(r rune) bool { return r == ',' || r == ' ' }) {
		from, to, isRange := strings.Cut(part, "-")
		if !isRange {
			to = from
		}
		start, err1 := strconv.Atoi(strings.TrimSpace(from))
		end, err2 := strconv.Atoi(strings.TrimSpace(to))
		if err1 != nil || err2 != nil {
			return nil, fmt.Errorf("%q is not a number or range", part)
		}
		if start > end {
			start, end = end, start
		}
		if start < 1 || end > max {
			return nil, fmt.Errorf("%q is out of range (1-%d)", part, max)
		}
		for n := start; n <= end; n++ {
			seen[n-1] = true
		}
	}

	picked := make([]int, 0, len(seen))
	for i := range seen {
		picked = append(picked, i)
	}
	sort.Ints(picked)
	return picked, nil
}

func (u *UI) printQuestion(question, def string) {
	if def != "" {
		fmt.Fprintf(u.Out, " %s %s %s ", Bold(Cyan("?")), question, Dim("("+def+")"))
		return
	}
	fmt.Fprintf(u.Out, " %s %s ", Bold(Cyan("?")), question)
}

func (u *UI) printOptions(options []string) {
	width := len(strconv.Itoa(len(options)))
	for i, option := range options {
		fmt.Fprintf(u.Out, "   %s %s\n", Yellow(fmt.Sprintf("%*d)", width, i+1)), option)
	}
}

// readLine reads one line of input without the trailing newline.
func (u *UI) readLine() (string, error) {
	line, err := u.in.ReadString('\n')
	if err == io.EOF && line == "" {
		fmt.Fprintln(u.Out)
		return "", ErrNoInput
	}
	if err != nil && err != io.EOF {
		return "", err
	}
	return strings.TrimSpace(line), nil
}
