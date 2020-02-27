package actionstat

// ActionMessage defines the schema for messages passed to AddAction.
type ActionMessage struct {
	Action *string
	Time   *float64
}
