package action

// statsImpl implements Stats
type statsImpl struct {
}

// addAction implements Stats.AddAction
func (impl *statsImpl) addAction(msg Message) {
}

// getStats implements Stats.GetStats
func (impl *statsImpl) getStats() string {
	return `{}`
}
