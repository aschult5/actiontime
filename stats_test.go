package actiontime

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"testing"
)

//
// Benchmarks
//

// Results must be read globally to prevent benchmark optimization
// Ref: https://dave.cheney.net/2013/06/30/how-to-write-benchmarks-in-go
var getStatsRes string

func BenchmarkGetStats(b *testing.B) {
	// Define GetStats benchmark table
	benchmarks := []struct {
		name   string
		numGet int
		numAct int
		numGo  int
	}{
		{"100Get_100Action_0Go", 100, 100, 0},
		{"100Get_100Action_1Go", 100, 100, 1},
		{"100Get_100Action_2Go", 100, 100, 2},
		{"100Get_100Action_4Go", 100, 100, 4},
		{"100Get_100Action_8Go", 100, 100, 8},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			var s Stats
			// Different action names will produce a larger output from GetStats
			err := addDifferent(bm.numAct, &s)
			if err != nil {
				b.Error(err)
				return
			}

			var getFun func()
			if bm.numGo <= 0 {
				// Don't spawn goroutines
				getFun = func() {
					for i := 0; i < bm.numGet; i++ {
						getStatsRes = s.GetStats()
					}
				}
			} else {
				// Limit number of running goroutines with a semaphore
				sem := make(chan bool, bm.numGo)
				getFun = func() {
					for i := 0; i < bm.numGet; i++ {
						sem <- true
						go func() {
							defer func() { <-sem }()
							getStatsRes = s.GetStats()
						}()
					}
					// Wait for remaining goroutines
					for i := 0; i < cap(sem); i++ {
						sem <- true
					}
					// Clear the sem
					for i := 0; i < cap(sem); i++ {
						<-sem
					}
				}
			}

			for n := 0; n < b.N; n++ {
				getFun()
			}
		})
	}
}

func BenchmarkAddAction(b *testing.B) {
	// Define AddAction benchmark table
	benchmarks := []struct {
		name    string
		numAdd  int
		numGo   int
		actions []string
	}{
		{"100Add_0Go", 100, 0, []string{"jump"}},
		{"100Add_1Go", 100, 1, []string{"jump"}},
		{"100Add_2Go", 100, 2, []string{"jump"}},
		{"100Add_4Go", 100, 4, []string{"jump"}},
		{"100Add_8Go", 100, 8, []string{"jump"}},
		{"100AddDiff_0Go", 100, 0, []string{"jump", "fall", "run", "walk", "sit", "stand"}},
		{"100AddDiff_1Go", 100, 1, []string{"jump", "fall", "run", "walk", "sit", "stand"}},
		{"100AddDiff_2Go", 100, 2, []string{"jump", "fall", "run", "walk", "sit", "stand"}},
		{"100AddDiff_4Go", 100, 4, []string{"jump", "fall", "run", "walk", "sit", "stand"}},
		{"100AddDiff_8Go", 100, 8, []string{"jump", "fall", "run", "walk", "sit", "stand"}},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			var s Stats
			// Limit number of running goroutines with a semaphore
			var addFun func()
			if bm.numGo <= 0 {
				// Don't spawn goroutines
				addFun = func() {
					for i := 0; i < bm.numAdd; i++ {
						// Loop over actions
						arg := fmt.Sprintf(`{"action":"%s","time":%d}`,
							bm.actions[i%len(bm.actions)], i+1)
						if err := s.AddAction(arg); err != nil {
							b.Log(arg)
							b.Error(err)
							break
						}
					}
				}
			} else {
				sem := make(chan bool, bm.numGo)
				addFun = func() {
					for i := 0; i < bm.numAdd; i++ {
						arg := fmt.Sprintf(`{"action":"%s","time":%d}`,
							bm.actions[i%len(bm.actions)], i+1)

						// Call AddAction concurrently, up to numGo at once
						sem <- true
						go func(a string) {
							defer func() { <-sem }()
							if err := s.AddAction(a); err != nil {
								b.Log(arg)
								b.Error(err)
							}
						}(arg)
					}
					// Wait for remaining goroutines
					for i := 0; i < cap(sem); i++ {
						sem <- true
					}
					// Clear the sem
					for i := 0; i < cap(sem); i++ {
						<-sem
					}
				}
			}

			for n := 0; n < b.N; n++ {
				addFun()
			}
		})
	}
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

//
// Tests
//
func TestAddAction(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		err   error
	}{
		{"Good", `{"action":"jump","time":100}`, nil},
		{"ExtraJson", `{"action": "jump", "time": 100, "extra": "value"}`, nil},
		{"NullJson", `null`, ErrBadInput},
		{"MissingField", `{"action": "jump"}`, ErrBadInput},
		{"EmptyAction", `{"action": "", "time": 100}`, ErrBadInput},
		{"LongAction", fmt.Sprintf(`{"action": "%s", "time": 100}`, strings.Repeat("a", MaxActionLen)), nil},
		{"TooLongAction", fmt.Sprintf(`{"action": "%s", "time": 100}`, strings.Repeat("b", MaxActionLen+1)), ErrBadInput},
		{"NegativeTime", `{"action": "jump", "time": -1}`, ErrBadInput},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Verify result of input matches expected error
			var obj Stats
			if err := obj.AddAction(tc.input); err != tc.err {
				t.Errorf(`Got: "%s"\nExpected: "%s"`, err, tc.err)
			}
		})
	}
}

func TestInvalidJson(t *testing.T) {
	obj := Stats{}
	err := obj.AddAction("{{")
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Error("Failed to detect json SyntaxError")
	}
}

func TestUnexpectedJson(t *testing.T) {
	obj := Stats{}
	err := obj.AddAction(`{"action": 1, "time": 1}`)
	if _, ok := err.(*json.UnmarshalTypeError); !ok {
		t.Error("Didn't detect unexpected json")
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
