package log

import (
	"context"
	"log/slog"
	"os"
	"runtime"
	"time"

	"github.com/lmittmann/tint"
)

type slogLogger struct {
	internal slog.Logger
}

var (
	_ abstractLogger = &slogLogger{}
)

// convenience constructor
func newSlogLogger(name string) slogLogger {
	internal := slog.New(tint.NewHandler(os.Stdout, &tint.Options{
		Level:     slog.LevelDebug,
		AddSource: true,
	}))

	return slogLogger{
		internal: *internal,
	}
}

func (l *slogLogger) log(level slog.Level, msg string, args ...interface{}) {
	ctx := context.Background()

	// we pass a custom program counter to slog so that it can
	// determine the source file and line number of the caller.
	// the source would otherwise point to the wrong location,
	// such as somewhere in this file itself.
	pc, _, _, _ := runtime.Caller(2)

	record := slog.Record{
		Time:    time.Now(),
		Message: msg,
		PC:      pc,
		Level:   level,
	}

	attrs := []slog.Attr{}
	for i := 0; i < len(args); i += 2 {
		key := args[i].(string)
		value := args[i+1]
		attrs = append(attrs, slog.Any(key, value))

	}
	handler := l.internal.Handler().WithAttrs(attrs)
	handler.Handle(ctx, record)
}

func (l *slogLogger) Debug(msg string, args ...interface{}) {
	l.log(slog.LevelDebug, msg, args...)
}
func (l *slogLogger) Info(msg string, args ...interface{}) {
	l.log(slog.LevelInfo, msg, args...)
}
func (l *slogLogger) Warn(msg string, args ...interface{}) {
	l.log(slog.LevelWarn, msg, args...)
}
func (l *slogLogger) Error(msg string, args ...interface{}) {
	l.log(slog.LevelError, msg, args...)
}

func (l *slogLogger) With(fields ...interface{}) abstractLogger {
	newLogger := slogLogger{
		internal: *l.internal.With(fields...),
	}
	return &newLogger
}

func (l *slogLogger) Named(name string) abstractLogger {
	newLogger := slogLogger{
		internal: *l.internal.With("name", name),
	}
	return &newLogger
}
