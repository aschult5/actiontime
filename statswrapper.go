package actiontime

import (
	"encoding/json"
	"fmt"
)

// statsWrapper offers a common test interface across Stats and statsImpl
type statsWrapper interface {
	Add(action string, value float64) error
	Get() []outputMessage
}

//
// statsWrap
//
type statsWrap struct {
	Stats
}

func (s *statsWrap) Add(action string, value float64) error {
	// Convert raw values into json-encoded string expected by Stats.AddAction
	return s.AddAction(fmt.Sprintf(`{"action":"%s","time":%f}`, action, value))
}

func (s *statsWrap) Get() []outputMessage {
	// Convert Stats.GetStats from string to []outputMessage
	var msgs []outputMessage
	// msgs will be nil on error
	json.Unmarshal([]byte(s.GetStats()), &msgs)
	return msgs
}

//
// statsImplWrap
//
type statsImplWrap struct {
	statsImpl
}

func (s *statsImplWrap) Add(action string, value float64) error {
	msg := inputMessage{action, value}
	s.addAction(msg)
	return nil
}

func (s *statsImplWrap) Get() []outputMessage {
	return s.getStats()
}
