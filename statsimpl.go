package action

// statsImpl implements Stats
type statsImpl struct {
	average map[string]Average
}

// addAction implements Stats.AddAction
func (impl *statsImpl) addAction(msg InputMessage) {
	if len(impl.average) == 0 {
		impl.average = make(map[string]Average)
	}
	avg := impl.average[*msg.Action]
	avg.Add(*msg.Time)
	impl.average[*msg.Action] = avg
}

// getStats implements Stats.GetStats
func (impl *statsImpl) getStats() []OutputMessage {
	ret := make([]OutputMessage, 0, len(impl.average))

	for key, avg := range impl.average {
		ret = append(ret, OutputMessage{key, avg.Value})
	}

	return ret
}
