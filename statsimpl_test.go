package actiontime

import (
	"fmt"
	"testing"
)

func TestEmpty(t *testing.T) {
	var obj statsImplWrap
	runner := testRunner{CsvFn: "testdata/tc_empty.csv", TestT: t, Obj: &obj}
	runner.Run()
}

func TestWrOne(t *testing.T) {
	var obj statsImplWrap
	runner := testRunner{CsvFn: "testdata/tc_wr_one_one.csv", TestT: t, Obj: &obj}
	runner.Run()
}

func TestWrFew(t *testing.T) {
	var obj statsImplWrap
	runner := testRunner{CsvFn: "testdata/tc_wr_few_few.csv", TestT: t, Obj: &obj}
	runner.Run()
}

func TestWrFewAsync(t *testing.T) {
	var obj statsImplWrap
	runner := testRunner{CsvFn: "testdata/tc_wr_few_few_async.csv", TestT: t, Obj: &obj}
	runner.Run()
}

func TestBalFewAsync(t *testing.T) {
	var obj statsImplWrap
	runner := testRunner{CsvFn: "testdata/tc_bal_few_few_async.csv", TestT: t, Obj: &obj}
	runner.Run()
}

func TestWrMilAsync(t *testing.T) {
	tc := "testdata/gen/tc_wr_mil_few_async.csv"
	if !fileExists(tc) {
		cmd := fmt.Sprintf("python3 tools/testgenerator.py --csv %s --add 1000000 jump run sit stand", tc)
		t.Errorf("Please generate %s with...\n`%s`", tc, cmd)
	} else {
		var obj statsImplWrap
		runner := testRunner{CsvFn: tc, TestT: t, Obj: &obj}
		runner.Run()
	}
}
