package ui

import "os"

// ANSI escape sequences used by rgit.
const (
	reset  = "\x1b[0m"
	bold   = "\x1b[1m"
	dim    = "\x1b[2m"
	red    = "\x1b[31m"
	green  = "\x1b[32m"
	yellow = "\x1b[33m"
	blue   = "\x1b[34m"
	cyan   = "\x1b[36m"
)

var colorEnabled = true

// SetColor turns colored output on or off.
func SetColor(enabled bool) { colorEnabled = enabled }

// ColorEnabled reports whether colored output is on.
func ColorEnabled() bool { return colorEnabled }

// DetectColor decides whether colors should be used by default: only when
// stdout is an interactive terminal and NO_COLOR is not set (https://no-color.org).
func DetectColor() bool {
	if _, set := os.LookupEnv("NO_COLOR"); set {
		return false
	}
	if os.Getenv("TERM") == "dumb" {
		return false
	}
	return isTerminal(os.Stdout) && enableVirtualTerminal(os.Stdout)
}

func paint(code, s string) string {
	if !colorEnabled || s == "" {
		return s
	}
	return code + s + reset
}

func Bold(s string) string   { return paint(bold, s) }
func Dim(s string) string    { return paint(dim, s) }
func Red(s string) string    { return paint(red, s) }
func Green(s string) string  { return paint(green, s) }
func Yellow(s string) string { return paint(yellow, s) }
func Blue(s string) string   { return paint(blue, s) }
func Cyan(s string) string   { return paint(cyan, s) }
