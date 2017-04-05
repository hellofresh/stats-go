package stats

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

// ErrUnknownClient is an error returned when trying to create stats client of unknown type
var ErrUnknownClient = errors.New("Unknown stats client type")

// Client is an interface for different methods of gathering stats
type Client interface {
	// BuildTimeTracker builds timer to track metric timings
	BuildTimeTracker() TimeTracker
	// Close closes underlying client connection if any
	Close() error

	// TrackRequest tracks HTTP Request stats
	TrackRequest(r *http.Request, tt TimeTracker, success bool) Client

	// TrackOperation tracks custom operation
	TrackOperation(section string, operation MetricOperation, tt TimeTracker, success bool) Client
	// TrackOperation tracks custom operation with n diff
	TrackOperationN(section string, operation MetricOperation, tt TimeTracker, n int, success bool) Client

	// SetHTTPMetricCallback sets callback handler that allows metric operation alteration for HTTP Request
	SetHTTPMetricCallback(callback HTTPMetricNameAlterCallback) Client

	// SetHTTPRequestSection sets metric section for HTTP Request metrics
	SetHTTPRequestSection(section string) Client

	// ResetHTTPRequestSection resets metric section for HTTP Request metrics to default value that is "request"
	ResetHTTPRequestSection() Client
}

// NewClient creates and builds new stats client instance by given dsn
func NewClient(dsn, prefix string) (Client, error) {
	// for backward compatibility
	if dsn == "" {
		return NewStatsdClient(dsn, prefix), nil
	}

	dsnURL, err := url.Parse(dsn)
	if err != nil {
		return nil, err
	}

	// backward compatibility statsd dsn - "<statsd.host>:<port>"
	if fmt.Sprintf("%s:%s", dsnURL.Scheme, dsnURL.Opaque) == dsn {
		return NewStatsdClient(dsn, prefix), nil
	}

	switch dsnURL.Scheme {
	case "statsd":
		return NewStatsdClient(dsnURL.Host, prefix), nil
	case "log":
		return NewLogClient(), nil
	case "memory":
		return NewMemoryClient(), nil
	}

	return nil, ErrUnknownClient
}
