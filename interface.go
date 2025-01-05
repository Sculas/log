package log

import "time"

// Fielder is an interface for providing fields to custom types.
type Fielder interface {
	Fields() Fields
}

// Fields represents a map of entry level data used for structured logging.
type Fields map[string]interface{}

// Interface represents the API of both Logger and Entry.
type Interface interface {
	WithFields(Fielder) Interface
	WithField(string, interface{}) Interface
	WithDuration(time.Duration) Interface
	WithError(error) Interface
	Named(string) Interface
	Debug(string)
	Info(string)
	Warn(string)
	Error(string)
	Fatal(string)
	Debugf(string, ...interface{})
	Infof(string, ...interface{})
	Warnf(string, ...interface{})
	Errorf(string, ...interface{})
	Fatalf(string, ...interface{})
	Trace(string) Interface
}
