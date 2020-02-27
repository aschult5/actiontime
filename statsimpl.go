package action

// statsImpl implements Stats
type statsImpl map[string]Average

// addAction implements Stats.AddAction
func (impl statsImpl) addAction(msg InputMessage) {
}

// getStats implements Stats.GetStats
func (impl statsImpl) getStats() string {
	return `{}`
}
