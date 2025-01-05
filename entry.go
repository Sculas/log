package log

import (
	"fmt"
	"os"
	"time"
)

// assert interface compliance.
var _ Interface = (*Entry)(nil)

// Now returns the current time.
var Now = time.Now

// Entry represents a single log entry.
type Entry struct {
	Logger    *Logger   `json:"-"`
	Name      string    `json:"name"`
	Fields    Fields    `json:"fields"`
	Level     Level     `json:"level"`
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
	start     time.Time
	fields    []Fields
}

// NewEntry returns a new entry for `log`.
func NewEntry(log *Logger) *Entry {
	return &Entry{
		Logger: log,
	}
}

// WithFields returns a new entry with `fields` set.
func (e *Entry) WithFields(fields Fielder) Interface {
	return e.withFields(fields)
}

// withFields returns a new entry with `fields` set.
// But to satisfy the interface, WithFields casts *Entry
// back into an `Interface`
func (e *Entry) withFields(fields Fielder) *Entry {
	var f []Fields
	f = append(f, e.fields...)
	f = append(f, fields.Fields())

	return &Entry{
		Logger: e.Logger,
		Name:   e.Name,
		fields: f,
	}
}

// WithField returns a new entry with the `key` and `value` set.
func (e *Entry) WithField(key string, value interface{}) Interface {
	return e.WithFields(Fields{key: value})
}

// WithDuration returns a new entry with the "duration" field set
// to the given duration in milliseconds.
func (e *Entry) WithDuration(d time.Duration) Interface {
	return e.WithField("duration", d.Milliseconds())
}

// WithError returns a new entry with the "error" set to `err`.
//
// The given error may implement .Fielder, if it does the method
// will add all its `.Fields()` into the returned entry.
func (e *Entry) WithError(err error) Interface {
	if err == nil {
		return e
	}
	ctx := e.WithField("error", err)
	if f, ok := err.(Fielder); ok {
		ctx = ctx.WithFields(f.Fields())
	}

	return ctx
}

// Named returns a new entry with the `name` set.
func (e *Entry) Named(name string) Interface {
	// If the entry already has a name, return a new entry with the new name.
	if e.Name != "" {
		return &Entry{
			Logger: e.Logger,
			Name:   e.Name + "/" + name,
			fields: e.fields,
		}
	}

	// If the entry does not have a name, return the entry itself with the new name.
	e.Name = name

	return e
}

// Debug level message.
func (e *Entry) Debug(msg string) {
	e.Logger.log(DebugLevel, e, msg)
}

// Info level message.
func (e *Entry) Info(msg string) {
	e.Logger.log(InfoLevel, e, msg)
}

// Warn level message.
func (e *Entry) Warn(msg string) {
	e.Logger.log(WarnLevel, e, msg)
}

// Error level message.
func (e *Entry) Error(msg string) {
	e.Logger.log(ErrorLevel, e, msg)
}

// Fatal level message, followed by an exit.
func (e *Entry) Fatal(msg string) {
	e.Logger.log(FatalLevel, e, msg)
	os.Exit(1)
}

// Debugf level formatted message.
func (e *Entry) Debugf(msg string, v ...interface{}) {
	e.Debug(fmt.Sprintf(msg, v...))
}

// Infof level formatted message.
func (e *Entry) Infof(msg string, v ...interface{}) {
	e.Info(fmt.Sprintf(msg, v...))
}

// Warnf level formatted message.
func (e *Entry) Warnf(msg string, v ...interface{}) {
	e.Warn(fmt.Sprintf(msg, v...))
}

// Errorf level formatted message.
func (e *Entry) Errorf(msg string, v ...interface{}) {
	e.Error(fmt.Sprintf(msg, v...))
}

// Fatalf level formatted message, followed by an exit.
func (e *Entry) Fatalf(msg string, v ...interface{}) {
	e.Fatal(fmt.Sprintf(msg, v...))
}

// Trace returns a new entry with a Stop method to fire off
// a corresponding completion log, useful with defer.
func (e *Entry) Trace(msg string) Interface {
	e.Info(msg)
	v := e.withFields(e.Fields)
	v.Message = msg
	v.start = time.Now()
	return v
}

// Stop should be used with Trace, to fire off the completion message. When
// an `err` is passed the "error" field is set, and the log level is error.
func (e *Entry) Stop(err *error) {
	if err == nil || *err == nil {
		e.WithDuration(time.Since(e.start)).Info(e.Message)
	} else {
		e.WithDuration(time.Since(e.start)).WithError(*err).Error(e.Message)
	}
}

// mergedFields returns the fields list collapsed into a single map.
func (e *Entry) mergedFields() Fields {
	f := Fields{}

	for _, fields := range e.fields {
		for k, v := range fields {
			f[k] = v
		}
	}

	return f
}

// finalize returns a copy of the Entry with Fields merged.
func (e *Entry) finalize(level Level, msg string) *Entry {
	return &Entry{
		Logger:    e.Logger,
		Name:      e.Name,
		Fields:    e.mergedFields(),
		Level:     level,
		Message:   msg,
		Timestamp: Now(),
	}
}
