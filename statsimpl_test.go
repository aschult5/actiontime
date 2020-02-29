package actiontime

import (
	"fmt"
	"testing"
)

func TestEmpty(t *testing.T) {
	runTestCase("testdata/tc_empty.csv", t)
}

func TestWrOne(t *testing.T) {
	runTestCase("testdata/tc_wr_one_one.csv", t)
}

func TestWrFew(t *testing.T) {
	runTestCase("testdata/tc_wr_few_few.csv", t)
}

func TestWrFewAsync(t *testing.T) {
	runTestCase("testdata/tc_wr_few_few_async.csv", t)
}

func TestBalFewAsync(t *testing.T) {
	runTestCase("testdata/tc_bal_few_few_async.csv", t)
}

func TestWrMilAsync(t *testing.T) {
	tc := "testdata/gen/tc_wr_mil_few_async.csv"
	if !fileExists(tc) {
		cmd := fmt.Sprintf("python3 tools/testgenerator.py --csv %s --add 1000000 jump run sit stand", tc)
		t.Errorf("Please generate %s with...\n`%s`", tc, cmd)
	} else {
		runTestCase(tc, t)
	}
}
