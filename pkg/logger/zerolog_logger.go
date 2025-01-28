package logger

import (
	"fmt"
	"io"
	"os"

	"github.com/rs/zerolog"
)

type ZeroLogLogger struct {
	Log zerolog.Logger
}

func NewZeroLogLogger(level int) Logger {
	log := zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "15:04",
	}).
		Level(zerolog.Level(level)).
		With().
		Timestamp().
		Int("pid", os.Getpid()).
		Logger()

	return &ZeroLogLogger{Log: log}
}

func (l *ZeroLogLogger) GetWriter() io.Writer {
	return l.Log
}

func (l *ZeroLogLogger) Printf(format string, args ...any) {
	l.Log.Printf(format, args...)
}

func (l *ZeroLogLogger) Error(args ...any) {
	l.Log.Error().Caller(1).Msg(fmt.Sprint(args...))
}

func (l *ZeroLogLogger) Errorf(format string, args ...any) {
	l.Log.Error().Caller(1).Msgf(format, args...)
}

func (l *ZeroLogLogger) Fatal(args ...any) {
	l.Log.Fatal().Caller(1).Msg(fmt.Sprint(args...))
}

func (l *ZeroLogLogger) Fatalf(format string, args ...any) {
	l.Log.Fatal().Caller(1).Msgf(format, args...)
}

func (l *ZeroLogLogger) Info(args ...any) {
	l.Log.Info().Caller(1).Msg(fmt.Sprint(args...))
}

func (l *ZeroLogLogger) Infof(format string, args ...any) {
	l.Log.Info().Caller(1).Msgf(format, args...)
}

func (l *ZeroLogLogger) Warn(args ...any) {
	l.Log.Warn().Caller(1).Msg(fmt.Sprint(args...))
}

func (l *ZeroLogLogger) Warnf(format string, args ...any) {
	l.Log.Warn().Caller(1).Msgf(format, args...)
}

func (l *ZeroLogLogger) Debug(args ...any) {
	l.Log.Debug().Caller(1).Msg(fmt.Sprint(args...))
}

func (l *ZeroLogLogger) Debugf(format string, args ...any) {
	l.Log.Debug().Caller(1).Msgf(format, args...)
}

func (l *ZeroLogLogger) WithField(key string, value any) Logger {
	var log zerolog.Logger
	if err, ok := value.(error); ok {
		log = l.Log.With().AnErr(key, err).Logger()
	} else {
		log = l.Log.With().Any(key, value).Logger()
	}

	return &ZeroLogLogger{
		Log: log,
	}
}

func (l *ZeroLogLogger) WithFields(fields map[string]any) Logger {
	logCtx := l.Log.With()
	for k, v := range fields {
		if errs, ok := v.([]error); ok {
			logCtx = logCtx.Errs(k, errs)
		} else if err, ok := v.(error); ok {
			logCtx = logCtx.AnErr(k, err)
		} else {
			logCtx = logCtx.Any(k, v)
		}
	}

	return &ZeroLogLogger{
		Log: logCtx.Logger(),
	}
}
