package log

import hclog "github.com/hashicorp/go-hclog"

type hcLogLogger struct {
	internal hclog.Logger
}

var (
	_ abstractLogger = &hcLogLogger{}
)

// convenience constructor
func newHCLogLogger(name string) hcLogLogger {
	return hcLogLogger{
		internal: hclog.New(&hclog.LoggerOptions{
			Name:  name,
			Level: hclog.LevelFromString("DEBUG"),
		}),
	}
}

func (l *hcLogLogger) Debug(msg string, args ...interface{}) {
	l.internal.Debug(msg, args...)
}
func (l *hcLogLogger) Info(msg string, args ...interface{}) {
	l.internal.Info(msg, args...)
}
func (l *hcLogLogger) Warn(msg string, args ...interface{}) {
	l.internal.Warn(msg, args...)
}
func (l *hcLogLogger) Error(msg string, args ...interface{}) {
	l.internal.Error(msg, args...)
}

func (l *hcLogLogger) With(fields ...interface{}) abstractLogger {
	newLogger := hcLogLogger{
		internal: l.internal.With(fields...),
	}
	return &newLogger
}

func (l *hcLogLogger) Named(name string) abstractLogger {
	newLogger := hcLogLogger{
		internal: l.internal.Named(name),
	}
	return &newLogger
}
