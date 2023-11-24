package log

// the interface all our loggers should implement,
// regardless of internal implementation
type abstractLogger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	With(args ...interface{}) abstractLogger
	Named(name string) abstractLogger
}

// the main struct for all logging. This generic interface
// will decouple the rest of the application from the actual
// logging implementation, and enable easy hot-swapping of
// logging libraries and other internal details.
type Logger struct {
	abstractLogger
}

// convenience constructor
func New(name string) Logger {
	internal := newSlogLogger("logger")
	return Logger{
		&internal,
	}
}
