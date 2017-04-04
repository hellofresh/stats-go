package stats

// LogClient is Client implementation for debug log
type LogClient struct {
	*StatsdClient
}

// NewLogClient builds and returns new LogClient instance
func NewLogClient() *LogClient {
	client := &LogClient{&StatsdClient{muted: true}}
	client.ResetHTTPRequestSection()

	return client
}
