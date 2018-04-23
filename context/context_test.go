package context

import (
	"context"
	"testing"

	"github.com/hellofresh/stats-go/client"
	"github.com/stretchr/testify/require"
)

func TestContext(t *testing.T) {
	t.Parallel()

	tests := []struct {
		scenario string
		function func(*testing.T)
	}{
		{
			scenario: "when set the context",
			function: testSetContext,
		},
		{
			scenario: "when the client is set on context",
			function: testGetFromContextSuccess,
		},
		{
			scenario: "when the client is not set on context",
			function: testGetFromContextFail,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			test.function(t)
		})
	}
}

func testSetContext(t *testing.T) {
	statsClient := client.NewNoop(false)

	ctx := context.Background()
	ctx = New(ctx, statsClient)

	client := WithContext(ctx)
	require.NotNil(t, client)
}

func testGetFromContextSuccess(t *testing.T) {
	statsClient := client.NewNoop(false)

	ctx := context.Background()
	ctx = New(ctx, statsClient)

	client := WithContext(ctx)
	require.NotNil(t, client)
}

func testGetFromContextFail(t *testing.T) {
	client := WithContext(context.Background())
	require.NotNil(t, client)
}
