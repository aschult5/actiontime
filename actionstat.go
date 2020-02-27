package actionstat

import (
	"encoding/json"
	"errors"
	"fmt"
)

// ActionStat tracks passed actions' average times.
type ActionStat struct {
}

var ErrMissingInput = errors.New("actionstat: Missing input data")

// AddAction takes json input and updates the action's time average.
// Keys are case insensitive.
// String values are case sensitive.
func (a ActionStat) AddAction(input string) error {
	var m ActionMessage

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
func (a ActionStat) GetStats() string {
	return "{}"
}
