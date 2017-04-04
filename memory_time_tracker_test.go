package stats

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMemoryTimeTracker(t *testing.T) {
	bucket := "test-bucket"

	tt := &MemoryTimeTracker{}
	tt.Start()
	tt.Finish(bucket)

	metric := tt.Elapsed()
	assert.Equal(t, bucket, metric.Bucket)
	assert.True(t, metric.Elapsed > time.Duration(0))
}
