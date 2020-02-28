package action

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"testing"
	"time"
)

var wg sync.WaitGroup

func TestEmpty(t *testing.T) {
	runTestCase("tc_empty.csv", t)
}

func TestOneAdd(t *testing.T) {
	runTestCase("tc_one_one.csv", t)
}

func TestFewAdds(t *testing.T) {
	runTestCase("tc_few_few.csv", t)
}

func TestFewAsync(t *testing.T) {
	runTestCase("tc_few_few_async.csv", t)
}

func TestMilliAsync(t *testing.T) {
	tc := "tc_mil_few_async.csv"
	genTest(tc, 5, 1e6)

	runTestCase(tc, t)
}

func genTest(fn string, numactions int, numAdds uint) error {
	const actionLen uint = 5

	// randSeq was taken from https://stackoverflow.com/a/22892986
	rand.Seed(time.Now().UnixNano())
	var letters = []rune("abcdefghijklmnopqrstuvwxyz")

	randSeq := func(n uint) string {
		b := make([]rune, n)
		for i := range b {
			b[i] = letters[rand.Intn(len(letters))]
		}
		return string(b)
	}

	var args = []string{fmt.Sprintf("--add %d", numAdds), fmt.Sprintf("--csv %s", fn)}
	for i := 0; i < numactions; i++ {
		args = append(args, randSeq(actionLen))
	}

	// Build command to call testgenerator
	cmd := exec.Command(
		"./tools/testgenerator.py",
		args...)
	fmt.Println(cmd)

	// Execute command
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(out)
		fmt.Println(err.Error())
		return err
	}

	return nil
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
