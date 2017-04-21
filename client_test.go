package stats

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	client, err := NewClient("", "")
	assert.NoError(t, err)
	assert.IsType(t, &StatsdClient{}, client)

	statsdClient, _ := client.(*StatsdClient)
	assert.True(t, statsdClient.muted)

	client, err = NewClient("statsd://", "")
	assert.NoError(t, err)
	assert.IsType(t, &StatsdClient{}, client)

	statsdClient, _ = client.(*StatsdClient)
	assert.True(t, statsdClient.muted)

	client, err = NewClient("log://", "")
	assert.NoError(t, err)
	assert.IsType(t, &LogClient{}, client)

	client, err = NewClient("memory://", "")
	assert.NoError(t, err)
	assert.IsType(t, &MemoryClient{}, client)

	client, err = NewClient("noop://", "")
	assert.NoError(t, err)
	assert.IsType(t, &NoopClient{}, client)

	client, err = NewClient("unknown://", "")
	assert.Nil(t, client)
	assert.Error(t, err)
	assert.Equal(t, ErrUnknownClient, err)
}
