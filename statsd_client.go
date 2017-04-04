package stats

import (
	"net/http"
	"sync"

	log "github.com/Sirupsen/logrus"
	statsd "gopkg.in/alexcesaro/statsd.v2"
)

// StatsdClient is Client implementation for statsd
type StatsdClient struct {
	sync.Mutex
	client             *statsd.Client
	muted              bool
	httpMetricCallback HTTPMetricNameAlterCallback
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

// BuildTimeTracker builds timer to track metric timings
func (sc *StatsdClient) BuildTimeTracker() TimeTracker {
	return NewTimeTracker(sc.client, sc.muted)
}

// Close statsd connection
func (sc *StatsdClient) Close() error {
	sc.client.Close()
	return nil
}

// TrackRequest tracks HTTP Request stats
func (sc *StatsdClient) TrackRequest(r *http.Request, tt TimeTracker, success bool) Client {
	b := NewBucketHTTPRequest(sc.httpRequestSection, r, success, sc.httpMetricCallback)
	i := NewIncrementer(sc.client, sc.muted)

	tt.Finish(b.Metric())
	i.IncrementAll(b)

	return sc
}

// TrackOperation tracks custom operation
func (sc *StatsdClient) TrackOperation(section string, operation MetricOperation, tt TimeTracker, success bool) Client {
	b := NewBucketPlain(section, operation, success)
	i := NewIncrementer(sc.client, sc.muted)

	if nil != tt {
		tt.Finish(b.MetricWithSuffix())
	}
	i.IncrementAll(b)

	return sc
}

// TrackOperationN tracks custom operation with n diff
func (sc *StatsdClient) TrackOperationN(section string, operation MetricOperation, tt TimeTracker, n int, success bool) Client {
	b := NewBucketPlain(section, operation, success)
	i := NewIncrementer(sc.client, sc.muted)

	if nil != tt {
		tt.Finish(b.MetricWithSuffix())
	}
	i.IncrementAllN(b, n)

	return sc
}

// SetHTTPMetricCallback sets callback handler that allows metric operation alteration for HTTP Request
func (sc *StatsdClient) SetHTTPMetricCallback(callback HTTPMetricNameAlterCallback) Client {
	sc.Lock()
	defer sc.Unlock()

	sc.httpMetricCallback = callback
	return sc
}

// SetHTTPRequestSection sets metric section for HTTP Request metrics
func (sc *StatsdClient) SetHTTPRequestSection(section string) Client {
	sc.Lock()
	defer sc.Unlock()

	sc.httpRequestSection = section
	return sc
}

// ResetHTTPRequestSection resets metric section for HTTP Request metrics to default value that is "request"
func (sc *StatsdClient) ResetHTTPRequestSection() Client {
	return sc.SetHTTPRequestSection(sectionRequest)
}
