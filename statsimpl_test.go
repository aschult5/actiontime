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
	if len(stats) != 0 {
		t.Errorf("Expected empty stats, not %v", stats)
	}
}
