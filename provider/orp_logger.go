package provider

import (
	"github.com/go-gorp/gorp"
	"github.com/krostar/nebulo/log"
)

// ORPLogger implement gorp.GorpLogger to allow
// the use of a custom logger for SQL actions
type ORPLogger struct {
	gorp.GorpLogger
}

// Printf print a SQL action
func (pl *ORPLogger) Printf(format string, args ...interface{}) {
	log.Logf(log.DEBUG, -1, format, args...)
}
