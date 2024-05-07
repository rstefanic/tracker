package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"tracker/display"
	"tracker/launchd"
	"tracker/tracker"

	"github.com/progrium/macdriver/macos"
	"github.com/progrium/macdriver/macos/appkit"
	"github.com/progrium/macdriver/macos/foundation"
)

func main() {
	launchdAgent := flag.Bool("launchd-agent", false, "Add the program as a Launch Agent for the current user.")
	flag.Parse()

	if *launchdAgent {
		launchd.AddAgent()
		os.Exit(0)
	}

	macos.RunApp(runTracker)
}

func runTracker(a appkit.Application, ad *appkit.ApplicationDelegate) {
	tracker := tracker.Init()
	display := display.NewDisplay()

	go handleExitSignal(func() {
		err := tracker.Save()
		if err != nil {
			fmt.Println("Saving of tracker failed.")
			fmt.Println(err)
		}
	})

	var workspace appkit.Workspace = appkit.Workspace_SharedWorkspace()
	var oneSec foundation.TimeInterval = 1.0
	var repeat bool = true

	foundation.Timer_ScheduledTimerWithTimeIntervalRepeatsBlock(oneSec, repeat, func(_ foundation.Timer) {
		app := workspace.FrontmostApplication()
		tracker.RecordUsage(app.LocalizedName())
		display.PrintUsage(tracker.Usage)
	})
}

func handleExitSignal(callback func()) {
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	<-sig
	callback()
	os.Exit(0)
}
