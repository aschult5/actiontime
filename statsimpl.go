package action

// statsImpl implements Stats
type statsImpl map[string]Average

// addAction implements Stats.AddAction
func (impl statsImpl) addAction(msg InputMessage) {
	avg := impl[*msg.Action]

	avg.Add(*msg.Time)
}

// getStats implements Stats.GetStats
func (impl statsImpl) getStats() []OutputMessage {
	ret := make([]OutputMessage, 0, len(impl))

	for key, avg := range impl {
		ret = append(ret, OutputMessage{key, avg.Value})
	}

	return ret
}
