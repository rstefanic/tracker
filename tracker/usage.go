package tracker

import "time"

type AppUsage struct {
	Name          string
	FocusDuration time.Duration
}

func (au AppUsage) GetFocusDuration() string {
	var t time.Time
	dur := time.Duration(au.FocusDuration)

	t = t.Add(dur * time.Second)
	return t.Format(time.TimeOnly)
}
