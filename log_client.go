package stats

// LogClient is Client implementation for debug log
type LogClient struct {
	*StatsdClient
}

// NewLogClient builds and returns new LogClient instance
func NewLogClient(unicode bool) *LogClient {
	client := &LogClient{&StatsdClient{muted: true, unicode: unicode}}
	client.ResetHTTPRequestSection()

	return client
}
