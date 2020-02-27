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

func TestInvalidJson(t *testing.T) {
	obj := ActionStat{}
	err := obj.AddAction("{{")
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Error("Failed to detect json SyntaxError")
	}
}

func TestExtraJson(t *testing.T) {
	obj := ActionStat{}
	err := obj.AddAction(`{"action": "jump", "time": 100, "extra": "value"}`)
	if err != nil {
		t.Error(err)
	}
}

func TestMissingJson(t *testing.T) {
	obj := ActionStat{}
	err := obj.AddAction(`{"action": "jump"}`)
	if err == nil {
		t.Error("Didn't detect missing parameter")
	}
}

func TestUnexpectedJson(t *testing.T) {
	obj := ActionStat{}
	err := obj.AddAction(`{"action": 1, "time": 1}`)
	if _, ok := err.(*json.UnmarshalTypeError); !ok {
		t.Error("Didn't detect unexpected json")
	}
}

func TestNullJson(t *testing.T) {
	obj := ActionStat{}
	err := obj.AddAction("null")
	if err == nil {
		t.Error("Didn't detect null json")
	}
}

func TestGetStats(t *testing.T) {
	obj := ActionStat{}
	s := obj.GetStats()
	if s != `{}` {
		t.Errorf("Expected empty json object, not %s", s)
	}
}
