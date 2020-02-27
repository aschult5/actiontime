package action

import "testing"

func TestImplGetStats(t *testing.T) {
	var impl statsImpl

	stats := impl.getStats()
	if len(stats) != 0 {
		t.Errorf("Expected empty stats, not %v", stats)
	}
}

func TestImplAddAndGet(t *testing.T) {
	var impl statsImpl

	// Form InputMessage
	action := "jump"
	time := 100.0
	msg := InputMessage{&action, &time}

	// Verify action was added
	impl.addAction(msg)
	stats := impl.getStats()
	if len(stats) != 1 {
		t.Errorf("Expected stats with 1 entry, not %v", stats)
	} else {
		// Verify OutputMessage matches expected
		expected := OutputMessage{action, time}
		if stats[0] != expected {
			t.Errorf("%v did not match expected %v", stats[0], expected)
		}
	}
}
