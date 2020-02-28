package average

import "testing"

func TestNaturalAverage(t *testing.T) {
	var avg Running

	avg.Add(5)
	avg.Add(15)
	expected := 10.0

	if avg.Value != expected {
		t.Errorf("Expected average of %f not %f", expected, avg.Value)
	}
}

func TestRationalAverage(t *testing.T) {
	var avg Running

	avg.Add(3)
	avg.Add(4)
	expected := 3.5

	if avg.Value != expected {
		t.Errorf("Expected average of %f not %f", expected, avg.Value)
	}
}

func TestEmptyAverage(t *testing.T) {
	var avg Running

	expected := 0.0

	if avg.Value != expected {
		t.Errorf("Expected average of %f not %f", expected, avg.Value)
	}
}
