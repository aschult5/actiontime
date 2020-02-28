package actiontime

import (
	"bytes"
	"encoding/json"
	"testing"
)

var goodInput []byte = []byte(`{"Action":"jump","Time":100}`)
var goodOutput []byte = []byte(`{"action":"jump","avg":100}`)

func TestInputUnmarshal(t *testing.T) {
	var msg inputMessage

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
	msg := inputMessage{&action, &num}

	b, err := json.Marshal(msg)

	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(b, goodInput) {
		t.Errorf("%s != %s", b, goodInput)
	}
}

func TestOutputUnmarshal(t *testing.T) {
	var msg outputMessage

	err := json.Unmarshal(goodOutput, &msg)
	if err != nil {
		t.Error(err)
	}

	if msg.Action != "jump" {
		t.Errorf("%s != %s", msg.Action, "jump")
	}

	if msg.Average != 100 {
		t.Errorf("%f != %d", msg.Average, 100)
	}
}

func TestOutputMarshal(t *testing.T) {
	msg := outputMessage{"jump", 100}

	b, err := json.Marshal(msg)

	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(b, goodOutput) {
		t.Errorf("%s != %s", b, goodOutput)
	}
}
