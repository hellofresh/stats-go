package middleware

import (
	"net/http"

	"github.com/felixge/httpsnoop"
	"github.com/hellofresh/stats-go/client"
	"github.com/hellofresh/stats-go/context"
	"github.com/hellofresh/stats-go/timer"
)

// New creates a new stats middleware
func New(s client.Client) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.New(r.Context(), s))

			mt := httpsnoop.CaptureMetrics(handler, w, r)
			t := timer.NewDuration(mt.Duration)

			success := mt.Code < http.StatusBadRequest
			s.TrackRequest(r, t, success)
		})
	}
}
