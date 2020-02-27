package action

// InputMessage defines the schema for messages passed to AddAction.
type InputMessage struct {
	Action *string
	Time   *float64
}
