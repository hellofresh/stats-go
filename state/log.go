package state

import (
	log "github.com/sirupsen/logrus"
)

// Log struct is State interface implementation that writes all states to log
type Log struct{}

// Set sets metric state
func (s *Log) Set(metric string, n int) {
	log.WithFields(log.Fields{
		"bucket": metric,
		"state":  n,
	}).Debug("Muted stats state send")
}
