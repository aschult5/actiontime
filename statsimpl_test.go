package action

import "testing"

func TestImplAddAction(t *testing.T) {
	var impl statsImpl
	action := "jump"
	time := 100.0

	msg := InputMessage{&action, &time}

	impl.addAction(msg)
}

func TestImplGetStats(t *testing.T) {
	var impl statsImpl

	stats := impl.getStats()
	if stats != `{}` {
		t.Errorf("Expected empty json object, not %s", stats)
	}
}
