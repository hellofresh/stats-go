package state

// Factory interface for making new state instances
type Factory interface {
	Create() State
}
