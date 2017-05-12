package state

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	instance := New(nil, true)
	assert.IsType(t, &Log{}, instance)

	instance = New(nil, false)
	assert.IsType(t, &Statsd{}, instance)
}
