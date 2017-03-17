package stats

import (
	"net/http"
	"sync"

	log "github.com/Sirupsen/logrus"
	statsd "gopkg.in/alexcesaro/statsd.v2"
)

var setterSync sync.Mutex

// StatsdStatsClient is StatsClient implementation for statsd
type StatsdStatsClient struct {
	client             *statsd.Client
	muted              bool
	httpMetricCallback HttpMetricNameAlterCallback
	httpRequestSection string
}

// NewStatsdStatsClient builds and returns new StatsdStatsClient instance
func NewStatsdStatsClient(dsn, prefix string) *StatsdStatsClient {
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

	client := &StatsdStatsClient{client: statsdClient, muted: muted}
	client.ResetHttpRequestSection()

	return client
}

// BuildTimeTracker builds timer to track metric timings
func (sc *StatsdStatsClient) BuildTimeTracker() TimeTracker {
	return NewTimeTracker(sc.client, sc.muted)
}

// Close statsd connection
func (sc *StatsdStatsClient) Close() error {
	sc.client.Close()
	return nil
}

// TrackRequest tracks HTTP Request stats
func (sc *StatsdStatsClient) TrackRequest(r *http.Request, tt TimeTracker, success bool) StatsClient {
	b := NewBucketHttpRequest(sc.httpRequestSection, r, success, sc.httpMetricCallback)
	i := NewIncrementer(sc.client, sc.muted)

	tt.Finish(b.Metric())
	i.IncrementAll(b)

	return sc
}

// TrackOperation tracks custom operation
func (sc *StatsdStatsClient) TrackOperation(section string, operation MetricOperation, tt TimeTracker, success bool) StatsClient {
	b := NewBucketPlain(section, operation, success)
	i := NewIncrementer(sc.client, sc.muted)

	if nil != tt {
		tt.Finish(b.MetricWithSuffix())
	}
	i.IncrementAll(b)

	return sc
}

// TrackOperationN tracks custom operation with n diff
func (sc *StatsdStatsClient) TrackOperationN(section string, operation MetricOperation, tt TimeTracker, n int, success bool) StatsClient {
	b := NewBucketPlain(section, operation, success)
	i := NewIncrementer(sc.client, sc.muted)

	if nil != tt {
		tt.Finish(b.MetricWithSuffix())
	}
	i.IncrementAllN(b, n)

	return sc
}

// SetHttpMetricCallback sets callback handler that allows metric operation alteration for HTTP Request
func (sc *StatsdStatsClient) SetHttpMetricCallback(callback HttpMetricNameAlterCallback) StatsClient {
	setterSync.Lock()
	defer setterSync.Unlock()

	sc.httpMetricCallback = callback
	return sc
}

// SetHttpRequestSection sets metric section for HTTP Request metrics
func (sc *StatsdStatsClient) SetHttpRequestSection(section string) StatsClient {
	setterSync.Lock()
	defer setterSync.Unlock()

	sc.httpRequestSection = section
	return sc
}

// ResetHttpRequestSection resets metric section for HTTP Request metrics to default value that is "request"
func (sc *StatsdStatsClient) ResetHttpRequestSection() StatsClient {
	return sc.SetHttpRequestSection(sectionRequest)
}
