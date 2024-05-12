package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"tracker/display"
	"tracker/launchd"
	tr "tracker/tracker"

	"github.com/progrium/macdriver/macos"
	"github.com/progrium/macdriver/macos/appkit"
	"github.com/progrium/macdriver/macos/foundation"
)

func main() {
	args := os.Args[1:] // Skip the program name
	if len(args) == 0 {
		help()
		os.Exit(0)
	}

	for _, cmd := range args {
		switch cmd {
		case "launchd":
			launchd.AddAgent()
			fmt.Println("Launch Agent installed.")
			break
		case "start":
			macos.RunApp(startTracker)
			break
		default:
			help()
		}
	}
}

func startTracker(a appkit.Application, ad *appkit.ApplicationDelegate) {
	workspace := appkit.Workspace_SharedWorkspace()
	tracker := tr.Init()
	display := display.NewDisplay(
		tracker.LongestAppNameLen(),
		len(tracker.Usage),
	)

	go saveOnExitSignal(tracker)

	newTrackerOnDayChange(workspace, tracker)
	trackUsageEverySecond(workspace, tracker)
	updateDisplayEverySecond(display, tracker)
	saveTrackerEveryMinute(tracker)
}

// Creates a new tracker when the calendar day changes.
func newTrackerOnDayChange(ws appkit.Workspace, t *tr.Tracker) {
	appkit.Workspace.NotificationCenter(ws).AddObserverForNameObjectQueueUsingBlock(
		foundation.CalendarDayChangedNotification,
		nil,
		foundation.OperationQueue_MainQueue(),
		func(notification foundation.Notification) {
			t.Save()
			t = tr.Init()
		},
	)
}

func trackUsageEverySecond(ws appkit.Workspace, t *tr.Tracker) {
	var oneSec foundation.TimeInterval = 1.0
	var repeat bool = true
	foundation.Timer_ScheduledTimerWithTimeIntervalRepeatsBlock(oneSec, repeat, func(_ foundation.Timer) {
		app := ws.FrontmostApplication()

		// `app` can be nil when no app is focused; guard against it
		if !app.IsNil() {
			t.RecordUsage(app.LocalizedName())
		}
	})
}

func updateDisplayEverySecond(d *display.Display, t *tr.Tracker) {
	var oneSec foundation.TimeInterval = 1.0
	var repeat bool = true
	foundation.Timer_ScheduledTimerWithTimeIntervalRepeatsBlock(oneSec, repeat, func(_ foundation.Timer) {
		d.PrintUsage(t.Usage)
	})
}

func saveTrackerEveryMinute(t *tr.Tracker) {
	var sixtySecs foundation.TimeInterval = 60.0
	var repeat bool = true
	foundation.Timer_ScheduledTimerWithTimeIntervalRepeatsBlock(sixtySecs, repeat, func(_ foundation.Timer) {
		t.Save()
	})
}

func saveOnExitSignal(t *tr.Tracker) {
	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs

	err := t.Save()
	if err != nil {
		fmt.Println("Saving of tracker failed.")
		fmt.Println(err)
	}

	os.Exit(0)
}

func help() {
	fmt.Print(`Usage: tracker [command]
command:
	launchd	Add the program as a Launch Agent for the current user.
	start	Start the tracker.
	help	Print this usage message.
`)
}
