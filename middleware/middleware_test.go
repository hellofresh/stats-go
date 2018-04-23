package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hellofresh/stats-go/client"
	"github.com/stretchr/testify/assert"
)

func TestMiddleware(t *testing.T) {
	t.Parallel()

	tests := []struct {
		scenario string
		function func(*testing.T, *http.Request, *httptest.ResponseRecorder)
	}{
		{
			scenario: "when a request is recorded",
			function: testRecorded,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()
			test.function(t, r, w)
		})
	}
}

func testRecorded(t *testing.T, r *http.Request, w *httptest.ResponseRecorder) {
	mClient := client.NewMemory(false)
	mw := New(mClient)

	mw(http.HandlerFunc(ping)).ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Len(t, mClient.TimerMetrics, 1)
}

// ping is a test handler
func ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte("OK\n"))
}
