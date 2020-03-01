// Package actiontime takes actions and times as json, tracking average times.
// Input is received as a json string, per requirements.
package actiontime

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Stats tracks passed actions' average times.
type Stats struct {
	statsImpl
}

// ErrBadInput indicates malformed input
var ErrBadInput = errors.New("actiontime: Malformed input data")

// AddAction takes json input and updates the action's time average.
// Example input: `{"action":"jump", "time": 100}`
// Keys are case insensitive.
// String values are case sensitive.
// Returns json.UnmarshalTypeError or actiontime.ErrBadInput on failure.
func (a *Stats) AddAction(input string) error {
	var msg inputMessage

	err := json.Unmarshal([]byte(input), &msg)
	if err != nil {
		return err
	}

	if msg.Action == nil || msg.Time == nil {
		return ErrBadInput
	}

	a.addAction(msg)
	return nil
}

// GetStats returns the averages of all action times as a json-encoded string.
func (a *Stats) GetStats() string {
	b, err := json.Marshal(a.getStats())
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return string(b)
}
