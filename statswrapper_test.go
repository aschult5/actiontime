package actiontime

import (
	"fmt"
	"testing"
)

func testStatsImpl(t *testing.T, csvFn string) {
	var obj statsImplWrap
	runner := testRunner{CsvFn: csvFn, Obj: &obj}
	t.Run(csvFn, runner.Run)
}

func testStats(t *testing.T, csvFn string) {
	var obj statsWrap
	runner := testRunner{CsvFn: csvFn, Obj: &obj}
	t.Run(csvFn, runner.Run)
}

func TestEmpty(t *testing.T) {
	tc := "testdata/tc_empty.csv"
	testStatsImpl(t, tc)
	testStats(t, tc)
}

func TestWrOne(t *testing.T) {
	tc := "testdata/tc_wr_one_one.csv"
	testStatsImpl(t, tc)
	testStats(t, tc)
}

func TestWrFew(t *testing.T) {
	tc := "testdata/tc_wr_few_few.csv"
	testStatsImpl(t, tc)
	testStats(t, tc)
}

func TestWrFewAsync(t *testing.T) {
	tc := "testdata/tc_wr_few_few_async.csv"
	testStatsImpl(t, tc)
	testStats(t, tc)
}

func TestBalFewAsync(t *testing.T) {
	tc := "testdata/tc_bal_few_few_async.csv"
	testStatsImpl(t, tc)
	testStats(t, tc)
}

func TestWrMilAsync(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	tc := "testdata/gen/tc_wr_mil_few_async.csv"
	if !fileExists(tc) {
		cmd := fmt.Sprintf("python3 tools/testgenerator.py --balance write --csv %s --add 1000000 jump run sit stand", tc)
		t.Errorf("Please generate %s with...\n%s", tc, cmd)
	} else {
		testStatsImpl(t, tc)
		testStats(t, tc)
	}
}
