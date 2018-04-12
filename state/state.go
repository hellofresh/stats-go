package state

// State is a metric state interface
type State interface {
	// Set sets metric state
	Set(metric string, n int, labels ...map[string]string)
}
