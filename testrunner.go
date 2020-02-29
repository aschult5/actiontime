package actiontime

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"sync"
	"testing"
)

var wg sync.WaitGroup

// testCommand represents a line in a given test case
type testCommand struct {
	Command string
	Action  string
	Value   float64
}

// runTestCase reads fn as a csv file containing test commands
func runTestCase(fn string, t *testing.T) {
	var impl statsImpl
	var cstats = make(chan []outputMessage)

	csvfile, err := os.Open(fn)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer csvfile.Close()

	// Read file line by line
	r := csv.NewReader(csvfile)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Error(err)
			continue
		}

		// Parse fields
		cmd, err := parseRecord(record)
		if err != nil {
			t.Error(err)
		}

		// Execute cmd
		err = executeCommand(cmd, &impl, cstats)
		if err != nil {
			t.Error(err)
		}
	}
}

// parseRecord attempts to parse an array of strings into a testCommand
func parseRecord(record []string) (testCommand, error) {
	var ret testCommand
	if len(record) != 3 {
		return ret, fmt.Errorf("Test case statements must be of form [<cmd>,<name>,<value>]")
	}

	cmd := record[0]
	action := record[1]
	val, err := strconv.ParseFloat(record[2], 64)
	if err != nil {
		return ret, err
	}

	ret = testCommand{cmd, action, val}
	return ret, nil
}

// executeCommand executes the given command on the passed *statsImpl
func executeCommand(cmd testCommand, impl *statsImpl, cstats chan []outputMessage) error {
	// Intepret command
	switch cmd.Command {

	case "sync":
		wg.Wait()

	case "addasync":
		wg.Add(1)
		// Copies of action and value made as `cmd` may go out of scope
		go func(action string, value float64) {
			defer wg.Done()
			msg := inputMessage{&action, &value}
			impl.addAction(msg)
		}(cmd.Action, cmd.Value)

	case "add":
		msg := inputMessage{&cmd.Action, &cmd.Value}
		impl.addAction(msg)

	case "getasync":
		wg.Add(2)
		go func() {
			defer wg.Done()
			// Ensure getStats is called by reading its return value
			cstats <- impl.getStats()
		}()
		go func() {
			defer wg.Done()
			<-cstats
		}()

	case "get":
		return handleStats(impl.getStats(), cmd.Action, cmd.Value)

	default:
		return fmt.Errorf("Unexpected command %s", cmd.Command)
	}

	return nil
}

func handleStats(stats []outputMessage, action string, value float64) error {
	const TOLERANCE = 0.000001
	if action == "_len_" {
		// Testing length
		if len(stats) != int(value) {
			return fmt.Errorf("stats length %d != %d", len(stats), int(value))
		}
	} else {
		// Testing action average
		found := false
		for _, msg := range stats {
			if msg.Action != action {
				continue
			}
			found = true

			if diff := math.Abs(msg.Average - value); diff > TOLERANCE {
				return fmt.Errorf(`"%s": average %f outside tolerance %f of %f`,
					msg.Action, msg.Average, TOLERANCE, value)
			}
			break
		}
		if !found {
			return fmt.Errorf(`Action "%s" not found`, action)
		}
	}

	return nil
}

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
