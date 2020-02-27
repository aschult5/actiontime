package action

// InputMessage defines the schema for messages passed to AddAction.
type InputMessage struct {
	Action *string
	Time   *float64
}

// OutputMessage defines the schema for messages returned from GetStats.
type OutputMessage struct {
	Action  string  `json:"action"`
	Average float64 `json:"avg"`
}
