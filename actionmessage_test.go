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

	if m.Action != "jump" {
		t.Errorf("%s != %s", m.Action, "jump")
	}

	if m.Time != 100 {
		t.Errorf("%f != %d", m.Time, 100)
	}
}

func TestJsonMarchas(t *testing.T) {
	m := ActionMessage{"jump", 100}

	b, err := json.Marshal(m)

	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(b, goodJson) {
		t.Errorf("%s != %s", b, goodJson)
	}
}
