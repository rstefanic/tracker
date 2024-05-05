package display

import (
	"fmt"

	"tracker/tracker"
)

// Display keeps information about the terminal output.
type Display struct {
	maxNameLen int
	linesLen   int
}

func NewDisplay() *Display {
	return &Display{
		maxNameLen: 0,
		linesLen:   0,
	}
}

func (d *Display) PrintUsage(usage []*tracker.AppUsage) {
	// Adjust display if there's a new AppUsage
	usageLen := len(usage)
	if usageLen > d.linesLen {
		d.linesLen = usageLen

		// Add a newline so moving the cursor to
		// the start includes this new AppUsage
		fmt.Println()

		newStat := usage[usageLen-1]
		newNameLen := len(newStat.Name)

		if newNameLen > d.maxNameLen {
			d.maxNameLen = newNameLen
		}
	}

	// Print the usage
	d.moveCursorToStart()
	for _, u := range usage {
		s := fmt.Sprintf("%-*s\t%s", d.maxNameLen, u.Name, u.GetFocusDuration())
		d.printLine(s)
	}
}

func (d *Display) printLine(s string) {
	d.clearLine() // Clear any existing text on this line
	fmt.Print(s)
	d.moveCursorDown()
}

func (d *Display) moveCursorToStart() {
	fmt.Printf("\x1b[%dF", d.linesLen)
}

func (d *Display) clearLine() {
	fmt.Print("\x1b[2K")
}

func (d *Display) moveCursorDown() {
	fmt.Print("\x1b[1E")
}
