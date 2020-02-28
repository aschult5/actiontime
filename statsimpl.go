package action

import (
	average "github.com/aschult5/actiontime/average"
	"sync"
)

// statsImpl implements Stats
type statsImpl struct {
	avg   map[string]average.Running
	avgMu sync.RWMutex
}

// addAction implements Stats.AddAction
func (impl *statsImpl) addAction(msg inputMessage) {
	impl.avgMu.Lock()
	defer impl.avgMu.Unlock()

	if len(impl.avg) == 0 {
		impl.avg = make(map[string]average.Running)
	}
	avg := impl.avg[*msg.Action]
	avg.Add(*msg.Time)
	impl.avg[*msg.Action] = avg
}

// getStats implements Stats.GetStats
func (impl *statsImpl) getStats() []outputMessage {
	impl.avgMu.RLock()
	defer impl.avgMu.RUnlock()

	ret := make([]outputMessage, 0, len(impl.avg))

	for key, avg := range impl.avg {
		ret = append(ret, outputMessage{key, avg.Value})
	}

	return ret
}
