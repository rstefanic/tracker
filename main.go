package main

import (
	"tracker/display"
	"tracker/tracker"

	"github.com/progrium/macdriver/macos"
	"github.com/progrium/macdriver/macos/appkit"
	"github.com/progrium/macdriver/macos/foundation"
)

func main() {
	tracker := tracker.NewTracker()
	display := display.NewDisplay()

	macos.RunApp(func(a appkit.Application, ad *appkit.ApplicationDelegate) {
		var workspace appkit.Workspace = appkit.Workspace_SharedWorkspace()
		var oneSec foundation.TimeInterval = 1.0
		var repeat bool = true

		foundation.Timer_ScheduledTimerWithTimeIntervalRepeatsBlock(oneSec, repeat, func(_ foundation.Timer) {
			app := workspace.FrontmostApplication()
			tracker.RecordUsage(app.LocalizedName())
			display.PrintUsage(tracker.Usage)
		})
	})
}
