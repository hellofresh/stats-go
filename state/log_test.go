package state

import (
	"testing"

	"github.com/hellofresh/stats-go/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLog_Set(t *testing.T) {
	var (
		logMessages []string
		logFields   []map[string]interface{}
	)
	log.SetHandler(func(msg string, fields map[string]interface{}, err error) {
		logMessages = append(logMessages, msg)
		logFields = append(logFields, fields)
	})

	metric1 := "metric1"
	metric2 := "metric2"
	metricState1 := 10
	metricState2 := 33

	logger := &Log{}
	logger.Set(metric1, metricState1)
	logger.Set(metric2, metricState2)

	require.Equal(t, 2, len(logMessages))
	assert.Equal(t, "Stats state set", logMessages[0])
	assert.Equal(t, metric1, logFields[0]["bucket"])
	assert.Equal(t, metricState1, logFields[0]["state"])

	assert.Equal(t, "Stats state set", logMessages[1])
	assert.Equal(t, metric2, logFields[1]["bucket"])
	assert.Equal(t, metricState2, logFields[1]["state"])
}
