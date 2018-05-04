package incrementer

import (
	"testing"

	"github.com/hellofresh/stats-go/bucket"
	"github.com/hellofresh/stats-go/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLog(t *testing.T) {
	var (
		logMessages []string
		logFields   []map[string]interface{}
	)
	log.SetHandler(func(msg string, fields map[string]interface{}, err error) {
		logMessages = append(logMessages, msg)
		logFields = append(logFields, fields)
	})

	b := "foo.bar.bucket"
	n := 42

	i := Log{}
	i.Increment(b)
	i.IncrementN(b, n)

	require.Equal(t, 2, len(logMessages))
	assert.Equal(t, "Stats counter incremented", logMessages[0])
	assert.Equal(t, b, logFields[0]["metric"])
	assert.Equal(t, "Stats counter incremented by n", logMessages[1])
	assert.Equal(t, b, logFields[1]["metric"])
	assert.Equal(t, n, logFields[1]["n"])

	logMessages = make([]string, 0)
	logFields = make([]map[string]interface{}, 0)

	bb := bucket.NewPlain("section", bucket.NewMetricOperation("o1", "o2", "o3"), true, true)
	i.IncrementAll(bb)

	assert.Equal(t, 4, len(logMessages))
	for j := 0; j < 4; j++ {
		assert.Equal(t, "Stats counter incremented", logMessages[j])
	}

	logMessages = make([]string, 0)
	logFields = make([]map[string]interface{}, 0)

	i.IncrementAllN(bb, n)

	assert.Equal(t, 4, len(logMessages))
	for j := 0; j < 4; j++ {
		assert.Equal(t, "Stats counter incremented by n", logMessages[j])
		assert.Equal(t, n, logFields[j]["n"])
	}
}
