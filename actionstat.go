package actionstat

type ActionStat struct {
}

func (a ActionStat) AddAction(json string) error {
	return nil
}

func (a ActionStat) GetStats() string {
	return "{}"
}
