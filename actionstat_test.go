package actionstat

import "testing"

func TestAddAction(t *testing.T) {
	obj := ActionStat{}
	e := obj.AddAction(`{"action": "jump", "time": 100}`)
	if e != nil {
		t.Fail()
	}
}
