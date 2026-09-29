//go:build !windows

package ui

import "os"

// PrepareConsole is a no-op outside Windows; terminals already speak UTF-8.
func PrepareConsole() {}

func isTerminal(f *os.File) bool {
	info, err := f.Stat()
	return err == nil && info.Mode()&os.ModeCharDevice != 0
}

func enableVirtualTerminal(*os.File) bool { return true }
