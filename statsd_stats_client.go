package stats

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	statsd "gopkg.in/alexcesaro/statsd.v2"
)

type StatsdStatsClient struct {
	client             *statsd.Client
	muted              bool
	httpMetricCallback HttpMetricNameAlterCallback
}

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

	return &StatsdStatsClient{client: statsdClient, muted: muted}
}

func (sc *StatsdStatsClient) BuildTimeTracker() TimeTracker {
	return NewTimeTracker(sc.client, sc.muted)
}

func (sc *StatsdStatsClient) Close() {
	sc.client.Close()
}

func (sc *StatsdStatsClient) TrackRequest(r *http.Request, tt TimeTracker, success bool) {
	b := NewBucketHttpRequest(r, success, sc.httpMetricCallback)
	i := NewIncrementer(sc.client, sc.muted)

	tt.Finish(b.Request())
	i.Increment(b.Request())
	i.Increment(b.TotalRequests())

	i.Increment(b.RequestsWithSuffix())
	i.Increment(b.TotalRequestsWithSuffix())
}

func (sc *StatsdStatsClient) TrackRoundTrip(r *http.Request, tt TimeTracker, success bool) {
	b := NewBucketHttpRequest(r, success, sc.httpMetricCallback)
	i := NewIncrementer(sc.client, sc.muted)

	tt.Finish(b.RoundTrip())
	i.Increment(b.RoundTrip())
	i.Increment(b.TotalRoundTrip())

	i.Increment(b.RoundTripWithSuffix())
	i.Increment(b.TotalRoundTripWithSuffix())
}

func (sc *StatsdStatsClient) TrackOperation(section string, operation MetricOperation, tt TimeTracker, success bool) {
	b := NewBucketPlain(section, operation, success)
	i := NewIncrementer(sc.client, sc.muted)

	if nil != tt {
		tt.Finish(b.MetricWithSuffix())
	}
	i.Increment(b.Metric())
	i.Increment(b.MetricWithSuffix())
	i.Increment(b.MetricTotal())
	i.Increment(b.MetricTotalWithSuffix())
}

func (sc *StatsdStatsClient) SetHttpMetricCallback(callback HttpMetricNameAlterCallback) StatsClient {
	sc.httpMetricCallback = callback
	return sc
}
