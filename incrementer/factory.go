package incrementer

// CounterFactory interface for making new CounterVec instances
type CounterFactory interface {
	Create(metric string, labelKeys []string) CounterVec
}

// Factory interface for making new incrementer instances
type Factory interface {
	Create() Incrementer
}
