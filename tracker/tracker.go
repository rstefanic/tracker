package tracker

import (
	"time"
)

// Tracker stores information about an App's Usage since the
// program launched. `Usage` is an array so that we can keep
// the AppUsage in the order that the programs were used.
type Tracker struct {
	FocusedApp string
	Date       time.Time
	Usage      []*AppUsage
	key        map[string]int // Maps the AppNames to their Usage index
	size       int            // Size of the `Usage` array
}

}

func NewTracker() *Tracker {
	return &Tracker{
		Usage: make([]*AppUsage, 0),
		Date:  time.Now(),
		key:   make(map[string]int),
		size:  0,
	}
}

func (t *Tracker) RecordUsage(app string) {
	idx, ok := t.key[app]

	if !ok {
		au := &AppUsage{
			Name:          app,
			FocusDuration: 0,
		}

		idx = t.size
		t.Usage = append(t.Usage, au)
		t.key[app] = idx

		t.size += 1
	}

	t.Usage[idx].FocusDuration += 1

	if app != t.FocusedApp {
		t.FocusedApp = app
	}
}
