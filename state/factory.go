package state

// StateFactory interface for making new state instances
type StateFactory interface {
	Create() State
}
