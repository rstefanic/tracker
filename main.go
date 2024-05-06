package main

import (
	"os"
	"os/signal"
	"syscall"

	"tracker/display"
	"tracker/tracker"

	"github.com/progrium/macdriver/macos"
	"github.com/progrium/macdriver/macos/appkit"
	"github.com/progrium/macdriver/macos/foundation"
)

func main() {
	macos.RunApp(runTracker)
}

func runTracker(a appkit.Application, ad *appkit.ApplicationDelegate) {
	go handleExitSignal()

	tracker := tracker.NewTracker()
	display := display.NewDisplay()

	var workspace appkit.Workspace = appkit.Workspace_SharedWorkspace()
	var oneSec foundation.TimeInterval = 1.0
	var repeat bool = true

	foundation.Timer_ScheduledTimerWithTimeIntervalRepeatsBlock(oneSec, repeat, func(_ foundation.Timer) {
		app := workspace.FrontmostApplication()
		tracker.RecordUsage(app.LocalizedName())
		display.PrintUsage(tracker.Usage)
	})
}

func handleExitSignal() {
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	os.Exit(0)
}
