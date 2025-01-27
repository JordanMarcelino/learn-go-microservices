package logger

import "io"

type Logger interface {
	GetWriter() io.Writer
	Printf(format string, args ...any)
	Error(args ...any)
	Errorf(format string, args ...any)
	Fatal(args ...any)
	Fatalf(format string, args ...any)
	Info(args ...any)
	Infof(format string, args ...any)
	Warn(args ...any)
	Warnf(format string, args ...any)
	Debug(args ...any)
	Debugf(format string, args ...any)
	WithField(key string, value any) Logger
	WithFields(fields map[string]any) Logger
}
