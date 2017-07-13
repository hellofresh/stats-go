package stats

import (
	"net/http"
	"sync"

	"github.com/hellofresh/stats-go/bucket"
	"github.com/hellofresh/stats-go/incrementer"
	"github.com/hellofresh/stats-go/state"
	"github.com/hellofresh/stats-go/timer"
	log "github.com/sirupsen/logrus"
	"gopkg.in/alexcesaro/statsd.v2"
)

// StatsdClient is Client implementation for statsd
type StatsdClient struct {
	sync.Mutex
	client             *statsd.Client
	muted              bool
	httpMetricCallback bucket.HTTPMetricNameAlterCallback
	httpRequestSection string
}

// NewStatsdClient builds and returns new StatsdClient instance
func NewStatsdClient(dsn, prefix string) *StatsdClient {
	var options []statsd.Option
	muted := false

	log.WithField("dsn", dsn).Info("Trying to connect to statsd instance")
	if len(dsn) == 0 {
		log.Debug("Statsd DSN not provided, client will be muted")
		options = append(options, statsd.Mute(true))
		muted = true
	} else {
		options = append(options, statsd.Address(dsn))
	}

	if len(prefix) > 0 {
		options = append(options, statsd.Prefix(prefix))
	}

	statsdClient, err := statsd.New(options...)
	if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"dsn":    dsn,
			"prefix": prefix,
		}).Warning("An error occurred while connecting to StatsD. Client will be muted.")
		muted = true
	}

	client := &StatsdClient{client: statsdClient, muted: muted}
	client.ResetHTTPRequestSection()

	return client
}

// BuildTimer builds timer to track metric timings
func (c *StatsdClient) BuildTimer() timer.Timer {
	return timer.New(c.client, c.muted)
}

// Close statsd connection
func (c *StatsdClient) Close() error {
	c.client.Close()
	return nil
}

// TrackRequest tracks HTTP Request stats
func (c *StatsdClient) TrackRequest(r *http.Request, t timer.Timer, success bool) Client {
	b := bucket.NewHTTPRequest(c.httpRequestSection, r, success, c.httpMetricCallback)
	i := incrementer.New(c.client, c.muted)

	t.Finish(b.Metric())
	i.IncrementAll(b)

	return c
}

// TrackOperation tracks custom operation
func (c *StatsdClient) TrackOperation(section string, operation bucket.MetricOperation, t timer.Timer, success bool) Client {
	b := bucket.NewPlain(section, operation, success)
	i := incrementer.New(c.client, c.muted)

	if nil != t {
		t.Finish(b.MetricWithSuffix())
	}
	i.IncrementAll(b)

	return c
}

// TrackOperationN tracks custom operation with n diff
func (c *StatsdClient) TrackOperationN(section string, operation bucket.MetricOperation, t timer.Timer, n int, success bool) Client {
	b := bucket.NewPlain(section, operation, success)
	i := incrementer.New(c.client, c.muted)

	if nil != t {
		t.Finish(b.MetricWithSuffix())
	}
	i.IncrementAllN(b, n)

	return c
}

// TrackMetric tracks custom metric, w/out ok/fail additional sections
func (c *StatsdClient) TrackMetric(section string, operation bucket.MetricOperation) Client {
	b := bucket.NewPlain(section, operation, true)
	i := incrementer.New(c.client, c.muted)

	i.Increment(b.Metric())
	i.Increment(b.MetricTotal())

	return c
}

// TrackMetricN tracks custom metric with n diff, w/out ok/fail additional sections
func (c *StatsdClient) TrackMetricN(section string, operation bucket.MetricOperation, n int) Client {
	b := bucket.NewPlain(section, operation, true)
	i := incrementer.New(c.client, c.muted)

	i.IncrementN(b.Metric(), n)
	i.IncrementN(b.MetricTotal(), n)

	return c
}

// TrackState tracks metric absolute value
func (c *StatsdClient) TrackState(section string, operation bucket.MetricOperation, value int) Client {
	b := bucket.NewPlain(section, operation, true)
	s := state.New(c.client, c.muted)

	s.Set(b.Metric(), value)

	return c
}

// SetHTTPMetricCallback sets callback handler that allows metric operation alteration for HTTP Request
func (c *StatsdClient) SetHTTPMetricCallback(callback bucket.HTTPMetricNameAlterCallback) Client {
	c.Lock()
	defer c.Unlock()

	c.httpMetricCallback = callback
	return c
}

// GetHTTPMetricCallback gets callback handler that allows metric operation alteration for HTTP Request
func (c *StatsdClient) GetHTTPMetricCallback() bucket.HTTPMetricNameAlterCallback {
	c.Lock()
	defer c.Unlock()

	return c.httpMetricCallback
}

// SetHTTPRequestSection sets metric section for HTTP Request metrics
func (c *StatsdClient) SetHTTPRequestSection(section string) Client {
	c.Lock()
	defer c.Unlock()

	c.httpRequestSection = section
	return c
}

// ResetHTTPRequestSection resets metric section for HTTP Request metrics to default value that is "request"
func (c *StatsdClient) ResetHTTPRequestSection() Client {
	return c.SetHTTPRequestSection(bucket.SectionRequest)
}
