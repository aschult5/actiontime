package actionstat

import (
	"bytes"
	"encoding/json"
	"testing"
)

var goodJson []byte = []byte(`{"Action":"jump","Time":100}`)

func TestJsonUnmarshal(t *testing.T) {
	var m ActionMessage

	err := json.Unmarshal(goodJson, &m)
	if err != nil {
		t.Error(err)
	}

	if m.Action == nil {
		t.Error("Failed to parse Action")
	}

	if m.Time == nil {
		t.Error("Failed to parse Time")
	}

	if *m.Action != "jump" {
		t.Errorf("%s != %s", *m.Action, "jump")
	}

	if *m.Time != 100 {
		t.Errorf("%f != %d", *m.Time, 100)
	}
}

func TestJsonMarshal(t *testing.T) {
	action := "jump"
	var num float64 = 100
	m := ActionMessage{&action, &num}

	b, err := json.Marshal(m)

	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(b, goodJson) {
		t.Errorf("%s != %s", b, goodJson)
	}
}
