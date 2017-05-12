package state

import (
	"io/ioutil"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/Sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func TestLog_Set(t *testing.T) {
	hook := test.NewGlobal()

	log.SetLevel(log.DebugLevel)
	log.SetOutput(ioutil.Discard)

	metric1 := "metric1"
	metric2 := "metric2"
	metricState1 := 10
	metricState2 := 33

	logger := &Log{}
	logger.Set(metric1, metricState1)
	logger.Set(metric2, metricState2)

	assert.Equal(t, 2, len(hook.Entries))
	assert.Equal(t, "Muted stats state send", hook.Entries[0].Message)
	assert.Equal(t, metric1, hook.Entries[0].Data["bucket"])
	assert.Equal(t, metricState1, hook.Entries[0].Data["state"])

	assert.Equal(t, "Muted stats state send", hook.Entries[1].Message)
	assert.Equal(t, metric2, hook.Entries[1].Data["bucket"])
	assert.Equal(t, metricState2, hook.Entries[1].Data["state"])
}
