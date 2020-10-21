package state

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMemory(t *testing.T) {
	memory := NewMemory()

	metric1 := "metric1"
	metric2 := "metric2"
	metricState1 := 10
	metricState2 := 33
	metricState12 := 42

	memory.Set(metric1, metricState1)
	memory.Set(metric2, metricState2)

	metrics := memory.Metrics()

	assert.Equal(t, metricState1, metrics[metric1])
	assert.Equal(t, metricState2, metrics[metric2])

	memory.Set(metric1, metricState12)
	assert.Equal(t, metricState12, metrics[metric1])
}
