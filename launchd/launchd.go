package launchd

import (
	_ "embed"
	"fmt"
	"os"
	"text/template"
)

//go:embed .plist.template
var plist string

type LaunchAgent struct {
	Program string
}

// Adds a property list (plist) for launchd as a
// third-party agent for the logged in user.
func AddAgent() error {
	exePath, err := os.Executable()
	if err != nil {
		return err
	}

	la := LaunchAgent{
		Program: exePath,
	}

	tmpl, err := template.New("launchdConfig").Parse(plist)
	if err != nil {
		return err
	}

	homedir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	pListPath := fmt.Sprintf("%s/Library/LaunchAgents/com.robertstefanic.tracker.plist", homedir)
	f, err := os.OpenFile(pListPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	err = tmpl.Execute(f, la)
	return err
}
