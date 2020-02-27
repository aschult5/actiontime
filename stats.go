package action

import (
	"encoding/json"
	"errors"
)

// Stats tracks passed actions' average times.
type Stats struct {
	statsImpl
}

var ErrMissingInput = errors.New("action: Missing input data")

// AddAction takes json input and updates the action's time average.
// Keys are case insensitive.
// String values are case sensitive.
func (a *Stats) AddAction(input string) error {
	var m Message

	err := json.Unmarshal([]byte(input), &m)
	if err != nil {
		return err
	}

	if m.Action == nil || m.Time == nil {
		return ErrMissingInput
	}

	a.addAction(m)
	return nil
}

// GetStats returns the averages of all actions as json.
func (a *Stats) GetStats() string {
	return a.getStats()
}
