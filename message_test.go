package action

import (
	"bytes"
	"encoding/json"
	"testing"
)

var goodInput []byte = []byte(`{"Action":"jump","Time":100}`)

func TestInputUnmarshal(t *testing.T) {
	var msg InputMessage

	err := json.Unmarshal(goodInput, &msg)
	if err != nil {
		t.Error(err)
	}

	if msg.Action == nil {
		t.Error("Failed to parse Action")
	}

	if msg.Time == nil {
		t.Error("Failed to parse Time")
	}

	if *msg.Action != "jump" {
		t.Errorf("%s != %s", *msg.Action, "jump")
	}

	if *msg.Time != 100 {
		t.Errorf("%f != %d", *msg.Time, 100)
	}
}

func TestInputMarshal(t *testing.T) {
	action := "jump"
	var num float64 = 100
	msg := InputMessage{&action, &num}

	b, err := json.Marshal(msg)

	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(b, goodInput) {
		t.Errorf("%s != %s", b, goodInput)
	}
}
