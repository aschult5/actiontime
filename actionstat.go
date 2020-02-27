package actionstat

// ActionStat tracks passed actions' average times.
type ActionStat struct {
}

// AddAction takes json input and updates the action's time average.
// Keys are case insensitive.
// String values are case sensitive.
func (a ActionStat) AddAction(json string) error {
	return nil
}

// GetStats returns the averages of all actions as json.
func (a ActionStat) GetStats() string {
	return "{}"
}
