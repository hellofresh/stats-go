package bucket

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMetricStorage_LooksLikeID(t *testing.T) {
	storage := newMetricStorage(25)
	firstSection := time.Now().Format(time.RFC3339Nano)

	for i := uint(0); i < storage.threshold-1; i++ {
		assert.False(t, storage.LooksLikeID(firstSection, fmt.Sprint(i)))
	}

	assert.True(t, storage.LooksLikeID(firstSection, fmt.Sprint(storage.threshold+1)))
	assert.True(t, storage.LooksLikeID(firstSection, fmt.Sprint(storage.threshold+1)))

	assert.Equal(t, storage.threshold, uint(len(storage.metrics[firstSection])))
}
