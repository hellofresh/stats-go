package incrementer

// CounterFactory interface for making new CounterVec instances
type CounterFactory interface {
	Create(metric string, labelKeys []string) CounterVec
}

// PrometheusIncrementerFactory interface for making new incrementer instances
type IncrementerFactory interface {
	Create() Incrementer
}
