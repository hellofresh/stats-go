package stats

import (
	"testing"

	"github.com/hellofresh/stats-go/client"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	statsClient, err := NewClient("log://")
	assert.NoError(t, err)
	assert.IsType(t, &client.Log{}, statsClient)

	statsClient, err = NewClient("memory://")
	assert.NoError(t, err)
	assert.IsType(t, &client.Memory{}, statsClient)

	statsClient, err = NewClient("noop://")
	assert.NoError(t, err)
	assert.IsType(t, &client.Noop{}, statsClient)

	statsClient, err = NewClient("unknown://")
	assert.Nil(t, statsClient)
	assert.Error(t, err)
	assert.Equal(t, ErrUnknownClient, err)
}
