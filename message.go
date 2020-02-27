package action

// Message defines the schema for messages passed to AddAction.
type Message struct {
	Action *string
	Time   *float64
}
