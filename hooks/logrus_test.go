package hooks

import (
	"io/ioutil"
	"testing"

	"github.com/hellofresh/stats-go/bucket"
	"github.com/hellofresh/stats-go/client"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNewLogrusHook(t *testing.T) {
	statsClient := client.NewMemory(true)
	section := "errors"
	b := bucket.NewPlain(section, bucket.MetricOperation{log.ErrorLevel.String()}, true, true)

	hook := NewLogrusHook(statsClient, section)
	log.AddHook(hook)
	log.SetLevel(log.DebugLevel)
	log.SetOutput(ioutil.Discard)

	log.Debug("debug")
	log.Info("info")
	log.Warn("warn")
	log.Error("error")

	assert.Equal(t, 2, len(statsClient.CountMetrics))
	assert.Equal(t, 1, statsClient.CountMetrics[b.Metric()])
	assert.Equal(t, 1, statsClient.CountMetrics[b.MetricTotal()])

	log.Errorf("error section: %s", section)
	log.Errorln("error line")

	assert.Equal(t, 2, len(statsClient.CountMetrics))
	assert.Equal(t, 3, statsClient.CountMetrics[b.Metric()])
	assert.Equal(t, 3, statsClient.CountMetrics[b.MetricTotal()])
}
