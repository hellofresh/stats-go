package context

import (
	"context"

	"github.com/hellofresh/stats-go/client"
)

type statsKeyType int

const statsKey statsKeyType = iota

// New returns a context that has a stats Client
func New(ctx context.Context, client client.Client) context.Context {
	return context.WithValue(ctx, statsKey, client)
}

// WithContext returns a stats Client with as much context as possible
func WithContext(ctx context.Context) client.Client {
	ctxStats, ok := ctx.Value(statsKey).(client.Client)
	if !ok {
		return client.NewNoop()
	}

	return ctxStats
}
