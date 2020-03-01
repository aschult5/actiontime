package actiontime

import (
	"encoding/json"
	"fmt"
)

// statsWrapper offers a common test interface across Stats and statsImpl
type statsWrapper interface {
	Add(action string, value float64)
	Get() []outputMessage
}

//
// statsWrap
//
type statsWrap struct {
	Stats
}

func (s *statsWrap) Add(action string, value float64) {
	// Convert raw values into json-encoded string expected by Stats.AddAction
	err := s.AddAction(fmt.Sprintf(`{"action":"%s","time":%f}`, action, value))
	if err != nil {
		fmt.Println(err)
	}
}

func (s *statsWrap) Get() []outputMessage {
	// Convert Stats.GetStats from string to []outputMessage
	var msgs []outputMessage
	err := json.Unmarshal([]byte(s.GetStats()), &msgs)
	if err != nil {
		fmt.Println(err)
	}
	return msgs
}

//
// statsImplWrap
//
type statsImplWrap struct {
	statsImpl
}

func (s *statsImplWrap) Add(action string, value float64) {
	msg := inputMessage{action, value}
	s.addAction(msg)
}

func (s *statsImplWrap) Get() []outputMessage {
	return s.getStats()
}
