package actionstat

import (
	"encoding/json"
	"testing"
)

func TestAddAction(t *testing.T) {
	// Convert valid message to byte array
	m := ActionMessage{"jump", 100}
	b, err := json.Marshal(m)
	if err != nil {
		t.Error(err)
	}

	// Verify valid ActionMessage doesn't produce an error
	obj := ActionStat{}
	err = obj.AddAction(string(b))
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
