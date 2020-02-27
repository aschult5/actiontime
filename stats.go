package action

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Stats tracks passed actions' average times.
type Stats struct {
}

var ErrMissingInput = errors.New("action: Missing input data")

// AddAction takes json input and updates the action's time average.
// Keys are case insensitive.
// String values are case sensitive.
func (a Stats) AddAction(input string) error {
	var m Message

	err := json.Unmarshal([]byte(input), &m)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if m.Action == nil || m.Time == nil {
		fmt.Println(ErrMissingInput)
		return ErrMissingInput
	}

	fmt.Println(m)
	return nil
}

// GetStats returns the averages of all actions as json.
func (a Stats) GetStats() string {
	return "{}"
}
