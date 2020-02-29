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

// testRunner encapsulates the functionality to run test case files
type testRunner struct {
	// CsvFn is the path to the test file
	CsvFn string

	// TestT is the corresponding test handle
	TestT *testing.T

	// Obj is the object under test
	Obj statsWrapper

	// wg is the WaitGroup for addasync and getasync commands
	wg sync.WaitGroup

	// cget communicates getasync results
	cget chan []outputMessage
}

// Run loops over CsvFn, parses lines, and executes commands
func (r *testRunner) Run() {
	r.cget = make(chan []outputMessage)

	csvfile, err := os.Open(r.CsvFn)
	if err != nil {
		r.TestT.Error(err)
		r.TestT.FailNow()
	}
	defer csvfile.Close()

	// Read file line by line
	rdr := csv.NewReader(csvfile)
	for {
		record, err := rdr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			r.TestT.Error(err)
			continue
		}

		// Parse fields
		cmd, err := parse(record)
		if err != nil {
			r.TestT.Error(err)
		}

		// Execute cmd
		err = r.execute(cmd)
		if err != nil {
			r.TestT.Error(err)
		}
	}
}

// execute executes the given command against the object under test
func (r *testRunner) execute(cmd testCommand) error {
	// Intepret command
	switch cmd.Command {

	case "sync":
		r.wg.Wait()

	case "addasync":
		r.wg.Add(1)
		// Copies of action and value made as `cmd` may go out of scope
		go func(action string, value float64) {
			defer r.wg.Done()
			r.Obj.Add(action, value)
		}(cmd.Action, cmd.Value)

	case "add":
		r.Obj.Add(cmd.Action, cmd.Value)

	case "getasync":
		r.wg.Add(2)
		go func() {
			defer r.wg.Done()
			// Ensure getStats is called by reading its return value.
			// The value isn't checked because the code isn't designed
			// to enforce order, so we can't reliably expect a certain avg.
			r.cget <- r.Obj.Get()
		}()
		go func() {
			defer r.wg.Done()
			<-r.cget
		}()

	case "get":
		return check(r.Obj.Get(), cmd.Action, cmd.Value)

	default:
		return fmt.Errorf("Unexpected command %s", cmd.Command)
	}

	return nil
}

func check(stats []outputMessage, action string, value float64) error {
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

// testCommand represents a line in a given test case
type testCommand struct {
	Command string
	Action  string
	Value   float64
}

// parse attempts to parse an array of strings into a testCommand
func parse(record []string) (testCommand, error) {
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

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
