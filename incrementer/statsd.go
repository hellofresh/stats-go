package incrementer

import (
	"github.com/hellofresh/stats-go/bucket"
	statsd "gopkg.in/alexcesaro/statsd.v2"
)

// Statsd struct is Incrementer interface implementation that writes all metrics to statsd
type Statsd struct {
	c *statsd.Client
}

// Increment increments metric in statsd
func (i *Statsd) Increment(metric string) {
	i.c.Increment(metric)
}

// IncrementN increments metric by n in statsd
func (i *Statsd) IncrementN(metric string, n int) {
	i.c.Count(metric, n)
}

// IncrementAll increments all metrics for given bucket in statsd
func (i *Statsd) IncrementAll(b bucket.Bucket) {
	incrementAll(i, b)
}

// IncrementAllN increments all metrics for given bucket in statsd
func (i *Statsd) IncrementAllN(b bucket.Bucket, n int) {
	incrementAllN(i, b, n)
}
