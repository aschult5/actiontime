package actionstat

import (
	"encoding/json"
	"fmt"
)

// ActionStat tracks passed actions' average times.
type ActionStat struct {
}

// AddAction takes json input and updates the action's time average.
// Keys are case insensitive.
// String values are case sensitive.
func (a ActionStat) AddAction(input string) error {
	var m ActionMessage

	err := json.Unmarshal([]byte(input), &m)
	if err == nil {
		fmt.Println(m)
	}

	return err
}

// GetStats returns the averages of all actions as json.
func (a ActionStat) GetStats() string {
	return "{}"
}
