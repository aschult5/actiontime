package actiontime

import (
	"fmt"
	"testing"
)

func TestEmpty(t *testing.T) {
	runner := testRunner{CsvFn: "testdata/tc_empty.csv", TestT: t}
	runner.Run()
}

func TestWrOne(t *testing.T) {
	runner := testRunner{CsvFn: "testdata/tc_wr_one_one.csv", TestT: t}
	runner.Run()
}

func TestWrFew(t *testing.T) {
	runner := testRunner{CsvFn: "testdata/tc_wr_few_few.csv", TestT: t}
	runner.Run()
}

func TestWrFewAsync(t *testing.T) {
	runner := testRunner{CsvFn: "testdata/tc_wr_few_few_async.csv", TestT: t}
	runner.Run()
}

func TestBalFewAsync(t *testing.T) {
	runner := testRunner{CsvFn: "testdata/tc_bal_few_few_async.csv", TestT: t}
	runner.Run()
}

func TestWrMilAsync(t *testing.T) {
	tc := "testdata/gen/tc_wr_mil_few_async.csv"
	if !fileExists(tc) {
		cmd := fmt.Sprintf("python3 tools/testgenerator.py --csv %s --add 1000000 jump run sit stand", tc)
		t.Errorf("Please generate %s with...\n`%s`", tc, cmd)
	} else {
		runner := testRunner{CsvFn: tc, TestT: t}
		runner.Run()
	}
}
