package action

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"
	"testing"
)

const tcprefix string = "statsimpl_tc_"

var wg sync.WaitGroup

func TestEmpty(t *testing.T) {
	runTestCase(tcprefix+"empty.csv", t)
}

func TestOneActionOneAdd(t *testing.T) {
	runTestCase(tcprefix+"one_one.csv", t)
}

func TestTwoActionFourAdd(t *testing.T) {
	runTestCase(tcprefix+"two_four.csv", t)
}

func TestTwoActionFourAddAsync(t *testing.T) {
	runTestCase(tcprefix+"two_four_async.csv", t)
}

// testCommand represents a line in a given test case
type testCommand struct {
	Command string
	Action  string
	Value   float64
}

// runTestCase reads fn as a csv file containing test commands
func runTestCase(fn string, t *testing.T) {
	var impl statsImpl

	csvfile, err := os.Open(fn)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

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
		err = executeCommand(cmd, &impl)
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
func executeCommand(cmd testCommand, impl *statsImpl) error {
	// Intepret command
	switch cmd.Command {

	case "sync":
		wg.Wait()

	case "addasync":
		wg.Add(1)
		// Copies of action and value made as `cmd` may go out of scope
		go func(action string, value float64) {
			defer wg.Done()
			msg := InputMessage{&action, &value}
			impl.addAction(msg)
		}(cmd.Action, cmd.Value)

	case "add":
		msg := InputMessage{&cmd.Action, &cmd.Value}
		impl.addAction(msg)

	case "get":
		return handleStats(impl.getStats(), cmd.Action, cmd.Value)

	default:
		return fmt.Errorf("Unexpected command %s", cmd.Command)
	}

	return nil
}

func handleStats(stats []OutputMessage, action string, value float64) error {
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
			if msg.Average != value {
				return fmt.Errorf(`"%s" average %f != %f`, msg.Action, msg.Average, value)
			}
			found = true
			break
		}
		if !found {
			return fmt.Errorf(`Action "%s" not found`, action)
		}
	}

	return nil
}
