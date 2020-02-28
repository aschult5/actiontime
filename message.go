package action

// inputMessage defines the schema for messages passed to AddAction.
type inputMessage struct {
	Action *string
	Time   *float64
}

// outputMessage defines the schema for messages returned from GetStats.
type outputMessage struct {
	Action  string  `json:"action"`
	Average float64 `json:"avg"`
}
