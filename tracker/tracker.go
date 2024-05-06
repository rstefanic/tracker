package tracker

import (
	"time"
)

// Tracker stores information about an App's Usage since the
// program launched. `Usage` is an array so that we can keep
// the AppUsage in the order that the programs were used.
type Tracker struct {
	FocusedApp string
	Usage      []*AppUsage
	key        map[string]int // Maps the AppNames to their Usage index
	size       int            // Size of the `Usage` array
}

type AppUsage struct {
	Name          string
	focusDuration time.Duration
}

func (au AppUsage) GetFocusDuration() string {
	var t time.Time
	dur := time.Duration(au.focusDuration)

	t = t.Add(dur * time.Second)
	return t.Format(time.TimeOnly)
}

func NewTracker() *Tracker {
	return &Tracker{
		Usage: make([]*AppUsage, 0),
		key:   make(map[string]int),
		size:  0,
	}
}

func (t *Tracker) RecordUsage(app string) {
	idx, ok := t.key[app]

	if !ok {
		au := &AppUsage{
			Name:          app,
			focusDuration: 0,
		}

		idx = t.size
		t.Usage = append(t.Usage, au)
		t.key[app] = idx

		t.size += 1
	}

	t.Usage[idx].focusDuration += 1

	if app != t.FocusedApp {
		t.FocusedApp = app
	}
}
