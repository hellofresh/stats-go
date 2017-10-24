package hooks

import (
	"io/ioutil"
	"testing"

	"github.com/hellofresh/stats-go"
	"github.com/hellofresh/stats-go/bucket"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNewLogrusHook(t *testing.T) {
	client := stats.NewMemoryClient(true)
	section := "errors"
	b := bucket.NewPlain(section, bucket.MetricOperation{log.ErrorLevel.String()}, true, true)

	hook := NewLogrusHook(client, section)
	log.AddHook(hook)
	log.SetLevel(log.DebugLevel)
	log.SetOutput(ioutil.Discard)

	log.Debug("debug")
	log.Info("info")
	log.Warn("warn")
	log.Error("error")

	assert.Equal(t, 2, len(client.CountMetrics))
	assert.Equal(t, 1, client.CountMetrics[b.Metric()])
	assert.Equal(t, 1, client.CountMetrics[b.MetricTotal()])

	log.Errorf("error section: %s", section)
	log.Errorln("error line")

	assert.Equal(t, 2, len(client.CountMetrics))
	assert.Equal(t, 3, client.CountMetrics[b.Metric()])
	assert.Equal(t, 3, client.CountMetrics[b.MetricTotal()])
}
