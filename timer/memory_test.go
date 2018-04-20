package timer

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMemoryTimeTracker(t *testing.T) {
	tt := &Memory{}
	tt.Start()
	d := tt.Finish()

	assert.True(t, d > time.Duration(0))
}
