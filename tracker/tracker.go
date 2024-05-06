package tracker

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"time"
)

// Tracker stores information about an App's Usage since the
// program launched. `Usage` is an array so that we can keep
// the AppUsage in the order that the programs were used.
type Tracker struct {
	FocusedApp string
	Date       time.Time
	Usage      []*AppUsage
	Key        map[string]int // Maps the AppNames to their Usage index
	Size       int            // Size of the `Usage` array
}

func (t *Tracker) getFilepath() (string, error) {
	date := t.Date.Format(time.DateOnly)
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/.tracker/%s.tracker", homedir, date), nil
}

func new() *Tracker {
	return &Tracker{
		Usage: make([]*AppUsage, 0),
		Date:  time.Now(),
		Key:   make(map[string]int),
		Size:  0,
	}
}

func Init() *Tracker {
	t, err := loadExistingTracker()
	if err != nil {
		return new()
	}

	return t
}

func loadExistingTracker() (*Tracker, error) {
	var t *Tracker = new()

	path, err := t.getFilepath()
	if err != nil {
		return nil, err
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	b := make([]byte, 1024)
	_, err = f.Read(b)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(b)
	dec := gob.NewDecoder(buf)
	err = dec.Decode(t)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (t *Tracker) Save() error {
	var buf bytes.Buffer

	enc := gob.NewEncoder(&buf)
	err := enc.Encode(t)
	if err != nil {
		return err
	}

	path, err := t.getFilepath()
	if err != nil {
		return err
	}

	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	f.Write(buf.Bytes())
	return nil
}

func (t *Tracker) RecordUsage(app string) {
	idx, ok := t.Key[app]

	if !ok {
		au := &AppUsage{
			Name:          app,
			FocusDuration: 0,
		}

		idx = t.Size
		t.Usage = append(t.Usage, au)
		t.Key[app] = idx

		t.Size += 1
	}

	t.Usage[idx].FocusDuration += 1

	if app != t.FocusedApp {
		t.FocusedApp = app
	}
}
