package bucket

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMetricStorage_LooksLikeID(t *testing.T) {
	storage := newMetricStorage()
	firstSection := time.Now().Format(time.RFC3339Nano)

	for i := 0; i < maxUniqueMetrics-1; i++ {
		assert.False(t, storage.LooksLikeID(firstSection, time.Now().Format(time.RFC3339Nano)))
	}

	assert.True(t, storage.LooksLikeID(firstSection, time.Now().Format(time.RFC3339Nano)))
	assert.True(t, storage.LooksLikeID(firstSection, time.Now().Format(time.RFC3339Nano)))

	assert.Equal(t, maxUniqueMetrics, len(storage.metrics[firstSection]))
}
