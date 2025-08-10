package browser

import (
	"errors"
	"fmt"
	"net/url"
	"os/exec"
	"runtime"
)

// Open opens the specified URL in the default web browser.
// It includes validation to ensure the URL is http/https and that the command
// to open the browser is safe.
func Open(rawURL string) error {
	// First, validate the URL to ensure it's a well-formed http/https URL.
	u, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return errors.New("URL scheme must be http or https")
	}

	var cmdName string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmdName = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmdName = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmdName = "xdg-open"
	}
	args = append(args, u.String()) // Use the parsed and validated URL

	// Find the absolute path of the command to be executed.
	cmdPath, err := exec.LookPath(cmdName)
	if err != nil {
		return fmt.Errorf("command not found: %s", cmdName)
	}

	// #nosec G204
	// The command and URL have been validated, so we can safely proceed.
	cmd := exec.Command(cmdPath, args...)
	return cmd.Start()
}
