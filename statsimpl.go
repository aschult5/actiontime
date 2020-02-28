package action

import "sync"

// statsImpl implements Stats
type statsImpl struct {
	average   map[string]Average
	averageMu sync.RWMutex
}

// addAction implements Stats.AddAction
func (impl *statsImpl) addAction(msg InputMessage) {
	impl.averageMu.Lock()
	defer impl.averageMu.Unlock()

	if len(impl.average) == 0 {
		impl.average = make(map[string]Average)
	}
	avg := impl.average[*msg.Action]
	avg.Add(*msg.Time)
	impl.average[*msg.Action] = avg
}

// getStats implements Stats.GetStats
func (impl *statsImpl) getStats() []OutputMessage {
	impl.averageMu.RLock()
	defer impl.averageMu.RUnlock()

	ret := make([]OutputMessage, 0, len(impl.average))

	for key, avg := range impl.average {
		ret = append(ret, OutputMessage{key, avg.Value})
	}

	return ret
}
