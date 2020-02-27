package actionstat

import "testing"

func TestAddAction(t *testing.T) {
	obj := ActionStat{}
	err = obj.AddAction(b)
	if err != nil {
		t.Error(err)
	}
}

func TestGetStats(t *testing.T) {
	obj := ActionStat{}
	s := obj.GetStats()
	if s != `{}` {
		t.Errorf("Expected empty json object, not %s", s)
	}
}
