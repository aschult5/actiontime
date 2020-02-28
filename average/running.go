package average

// Running maintains a running average by tracking the count of entries
type Running struct {
	Value float64 `json:"avg"`
	Count uint64  `json:"-"`
}

// Add updates the average as Avg[n] = Avg[n-1]*(n-1)/n + Val[n]/n
func (avg *Running) Add(val float64) {
	avg.Value *= float64(avg.Count) / float64(avg.Count+1)
	avg.Value += val / float64(avg.Count+1)
	avg.Count++
}
