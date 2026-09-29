package commands

import (
	"fmt"
	"net/url"
	"os/exec"
	"runtime"
)

// openBrowser opens an http(s) address in the default web browser.
func openBrowser(address string) error {
	u, err := url.Parse(address)
	if err != nil || (u.Scheme != "https" && u.Scheme != "http") {
		return fmt.Errorf("refusing to open %q: not a web address", address)
	}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", address)
	case "darwin":
		cmd = exec.Command("open", address)
	default:
		cmd = exec.Command("xdg-open", address)
	}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("could not open a browser: %w (address: %s)", err, address)
	}
	return nil
}
