package actiontime

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"testing"
)

var getStatsRes string

func BenchmarkGetStats100(b *testing.B) {
	var s Stats
	// Different action names will produce a larger output from GetStats
	err := addDifferent(100, &s)
	if err != nil {
		b.Error(err)
		b.FailNow()
	}

	// Read result locally and again globally to prevent optimization
	// Ref: https://dave.cheney.net/2013/06/30/how-to-write-benchmarks-in-go
	var r string
	for n := 0; n < b.N; n++ {
		r = s.GetStats()
	}
	getStatsRes = r
}

func BenchmarkAddAction100(b *testing.B) {
	for n := 0; n < b.N; n++ {
		s := Stats{}
		err := addSame(100, &s)
		if err != nil {
			b.Error(err)
			break
		}
	}
}

// addSame calls AddAction n times with the same action but different times
func addSame(n int, s *Stats) error {
	for i := 1; i <= n; i++ {
		err := s.AddAction(fmt.Sprintf(`{"action":"stand","time":%d}`, i))
		if err != nil {
			return err
		}
	}
	return nil
}

// addDifferent calls AddAction n times with different actions and times
func addDifferent(n int, s *Stats) error {
	// Always seed with same value, we don't actually want random results
	rand.Seed(42)

	for i := 1; i <= n; i++ {
		action := randStringBytes(5)
		err := s.AddAction(fmt.Sprintf(`{"action":"%s","time":%d}`, action, i))
		if err != nil {
			return err
		}
	}
	return nil
}

// randStringBytes produces an English-alphabet string of length n
func randStringBytes(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func TestAddAction(t *testing.T) {
	// Form a valid input string
	action := "jump"
	var time float64 = 100
	str := getInputMessageString(action, time)

	// Verify valid inputMessage doesn't produce an error
	obj := Stats{}
	err := obj.AddAction(str)
	if err != nil {
		t.Error(err)
	}
}

func TestInvalidJson(t *testing.T) {
	obj := Stats{}
	err := obj.AddAction("{{")
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Error("Failed to detect json SyntaxError")
	}
}

func TestExtraJson(t *testing.T) {
	obj := Stats{}
	err := obj.AddAction(`{"action": "jump", "time": 100, "extra": "value"}`)
	if err != nil {
		t.Error(err)
	}
}

func TestBadJson(t *testing.T) {
	obj := Stats{}
	err := obj.AddAction(`{"action": "jump"}`)
	if err != ErrBadInput {
		t.Error("Didn't detect missing parameter")
	}
}

func TestUnexpectedJson(t *testing.T) {
	obj := Stats{}
	err := obj.AddAction(`{"action": 1, "time": 1}`)
	if _, ok := err.(*json.UnmarshalTypeError); !ok {
		t.Error("Didn't detect unexpected json")
	}
}

func TestEmptyAction(t *testing.T) {
	obj := Stats{}
	err := obj.AddAction(`{"action": "", "time": 1}`)
	if err != ErrBadInput {
		t.Error("Didn't detect empty action string")
	}
}

func TestLongActionName(t *testing.T) {
	obj := Stats{}

	// Test long but acceptable action name
	long := strings.Repeat("a", MaxActionLen)
	err := obj.AddAction(fmt.Sprintf(`
		{
		"action": "%s",
		"time": 1
		}`, long))
	if err != nil {
		t.Error("Didn't allow MaxActionLen string")
	}

	// Test toolong of an action name
	toolong := strings.Repeat("b", MaxActionLen+1)
	err = obj.AddAction(fmt.Sprintf(`
		{
		"action": "%s",
		"time": 1
		}`, toolong))
	if err != ErrBadInput {
		t.Error("Didn't detect long action string")
	}
}

func TestNegativeTime(t *testing.T) {
	obj := Stats{}
	err := obj.AddAction(`{"action": "jump", "time": -1}`)
	if err != ErrBadInput {
		t.Error("Didn't detect negative time")
	}
}

func TestNullJson(t *testing.T) {
	obj := Stats{}
	err := obj.AddAction("null")
	if err != ErrBadInput {
		t.Error("Didn't detect null json")
	}
}

func TestGetStats(t *testing.T) {
	obj := Stats{}
	s := obj.GetStats()
	if s != `[]` {
		t.Errorf("Expected empty json object, not %s", s)
	}
}

func TestAddAndGet(t *testing.T) {
	// Form a valid input string
	action := "jump"
	var time float64 = 100
	istr := getInputMessageString(action, time)

	// Add the action
	obj := Stats{}
	obj.AddAction(istr)

	// Retrieve the stats and verify them
	ostr := obj.GetStats()
	var messages []outputMessage
	err := json.Unmarshal([]byte(ostr), &messages)
	if err != nil {
		t.Error(err)
	}

	if len(messages) != 1 {
		t.Errorf("Expected stats with 1 entry, not %v", messages)
	} else {
		expected := outputMessage{action, time}
		if messages[0] != expected {
			t.Errorf("%v did not match expected %v", messages[0], expected)
		}
	}
}

// getInputMessageString converts valid message values to a json string
func getInputMessageString(action string, time float64) string {
	msg := inputMessage{action, time}
	b, _ := json.Marshal(msg)
	return string(b)
}
