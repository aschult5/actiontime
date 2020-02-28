package action

import (
	"encoding/csv"
	"io"
	"os"
	"strconv"
	"testing"
)

const tcprefix string = "statsimpl_tc_"

func TestEmpty(t *testing.T) {
	testExecutor(tcprefix+"empty.csv", t)
}

func TestOneActionOneAdd(t *testing.T) {
	testExecutor(tcprefix+"one_one.csv", t)
}

func testExecutor(fn string, t *testing.T) {
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
		if len(record) != 3 {
			t.Errorf("Test case statements must be of form [<cmd>,<name>,<value>]")
			t.FailNow()
		}

		// Parse fields
		action := record[1]
		num, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			t.Error(err)
			continue
		}

		// Intepret command
		switch record[0] {

		case "add":
			msg := InputMessage{&action, &num}
			impl.addAction(msg)

		case "get":
			stats := impl.getStats()
			if action == "_len_" {
				// Testing length
				if len(stats) != int(num) {
					t.Errorf("stats length %d != %d", len(stats), int(num))
					break
				}
			} else {
				// Testing action average
				found := false
				for _, msg := range stats {
					if msg.Action != action {
						continue
					}
					if msg.Average != num {
						t.Errorf(`"%s" average %f != %f`, msg.Action, msg.Average, num)
						break
					}
					found = true
					break
				}
				if !found {
					t.Errorf(`Action "%s" not found`, action)
				}
			}

		default:
			t.Errorf("Unexpected command %s", record[0])
		}
	}
}
