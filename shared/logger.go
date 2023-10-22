package shared

import (
	"github.com/hashicorp/go-hclog"
	"github.com/vgarvardt/gue/v5/adapter"
)

// ensure the gue.Logger interface is implemented by our logger
var (
	_ adapter.Logger = &AdaptedLogger{}
)

type AdaptedLogger struct {
	Logger hclog.Logger
}

// convenience constructor
func NewLogger(name string, log hclog.Logger) adapter.Logger {
	return &AdaptedLogger{
		Logger: log.Named(name),
	}
}

// splat out the lsit of gue structs into a list for hclog
func splat(fields []adapter.Field) []interface{} {
	var flat []interface{}
	for _, field := range fields {
		flat = append(flat, field.Key, field.Value)
	}
	return flat
}

// arbitrary helper to DRY up some code
// func adapt(msg string, fields []adapter.Field, logFn func(msg string, args ...interface{})) {
// 	logFn(msg, splat(fields)...)
// }

func (l *AdaptedLogger) Debug(msg string, fields ...adapter.Field) {
	l.Logger.Debug(msg, splat(fields)...)
}
func (l *AdaptedLogger) Info(msg string, fields ...adapter.Field) {
	l.Logger.Info(msg, splat(fields)...)
}
func (l *AdaptedLogger) Error(msg string, fields ...adapter.Field) {
	l.Logger.Error(msg, splat(fields)...)
}

func (l *AdaptedLogger) With(fields ...adapter.Field) adapter.Logger {
	l.Logger = l.Logger.With(splat(fields)...)
	return l
}

var L = hclog.New(&hclog.LoggerOptions{
	Level: hclog.LevelFromString("TRACE"), // @level
	// IncludeLocation: true,                           // @caller
	JSONFormat: false,
	Color:      hclog.AutoColor,
	// ColorHeaderOnly: true,
	// DisableTime: true,
})
