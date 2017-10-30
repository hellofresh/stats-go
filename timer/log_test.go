package timer

import (
	"testing"

	"github.com/hellofresh/stats-go/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	var (
		logMessages []string
		logFields   []map[string]interface{}
	)
	log.SetHandler(func(msg string, fields map[string]interface{}, err error) {
		logMessages = append(logMessages, msg)
		logFields = append(logFields, fields)
	})

	b := "foo.bar.bucket"

	tr := &Log{}
	tr.Start().Finish(b)

	require.Equal(t, 1, len(logMessages))
	assert.Equal(t, "Stats timer finished", logMessages[0])
	assert.Equal(t, b, logFields[0]["bucket"])
}
