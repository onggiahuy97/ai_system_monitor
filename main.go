package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"
)

// getActiveWindowInfo uses AppleScript to get name of frontmost application
// and the title of its active window
func getActiveWindowInfo() (string, string, error) {
	// This is AppleScript command
	script := `
		tell application "System Events"
			set frontApp to first application process whose frontmost is true
			set frontAppName to name of frontApp
			set windowTitle to ""
			try
				set windowTitle to name of window 1 of frontApp
			end try
			return frontAppName & "," & windowTitle
		end tell
	`
	cmd := exec.Command("osascript", "-e", script)

	// Capture the output
	var out bytes.Buffer
	cmd.Stdout = &out

	// Run the command and check for error
	err := cmd.Run()
	if err != nil {
		return "", "", fmt.Errorf("failed to run Apple Script: %w", err)
	}

	// output out.String()
	output := strings.TrimSpace(out.String())

	parts := strings.SplitN(output, ",", 2)
	if len(parts) != 2 {
		return parts[0], "", nil // Return app name even if the title missing
	}

	return parts[0], parts[1], nil
}

func main() {
	for {
		appName, windowTitle, err := getActiveWindowInfo()
		if err != nil {
			log.Printf("Error: %v", err)
		} else {
			fmt.Printf("Time: %s | App: %s | Title: %s\n", time.Now(), appName, windowTitle)
		}

	}
}
