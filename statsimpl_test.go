package action

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"testing"
)

const tcprefix string = "statsimpl_tc_"

func TestEmpty(t *testing.T) {
	runTestCase(tcprefix+"empty.csv", t)
}

func TestOneActionOneAdd(t *testing.T) {
	runTestCase(tcprefix+"one_one.csv", t)
}

func TestTwoActionFourAdd(t *testing.T) {
	runTestCase(tcprefix+"two_four.csv", t)
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

	case "add":
		msg := InputMessage{&cmd.Action, &cmd.Value}
		impl.addAction(msg)

	case "get":
		stats := impl.getStats()
		if cmd.Action == "_len_" {
			// Testing length
			if len(stats) != int(cmd.Value) {
				return fmt.Errorf("stats length %d != %d", len(stats), int(cmd.Value))
			}
		} else {
			// Testing action average
			found := false
			for _, msg := range stats {
				if msg.Action != cmd.Action {
					continue
				}
				if msg.Average != cmd.Value {
					return fmt.Errorf(`"%s" average %f != %f`, msg.Action, msg.Average, cmd.Value)
				}
				found = true
				break
			}
			if !found {
				return fmt.Errorf(`Action "%s" not found`, cmd.Action)
			}
		}

	default:
		return fmt.Errorf("Unexpected command %s", cmd.Command)
	}

	return nil
}
