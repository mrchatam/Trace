package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// openBrowserFn opens a URL in the default browser. Injectable for tests.
var openBrowserFn = openBrowser

// openBrowser best-effort opens url via the platform default handler.
// darwin: open; windows: cmd /c start; else: xdg-open. No new module deps.
func openBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return err
	}
	_ = cmd.Process.Release()
	return nil
}

func tipOpenManually(url string) {
	fmt.Fprintf(os.Stderr, "trace gui: could not open browser; open manually: %s\n", url)
}

// tipOpenManuallyFn is called when browser open fails; injectable for tests.
var tipOpenManuallyFn = tipOpenManually
