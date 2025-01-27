package logger

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

type zerologLogger struct {
	log zerolog.Logger
}

func NewZeroLogLogger() Logger {
	log := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeLocation: time.UTC, TimeFormat: zerolog.TimeFormatUnix}).
		With().
		Timestamp().
		Int("pid", os.Getpid()).
		Logger()

	return &zerologLogger{log: log}
}

func (l *zerologLogger) GetWriter() io.Writer {
	return l.log
}

func (l *zerologLogger) Printf(format string, args ...any) {
	l.log.Printf(format, args...)
}

func (l *zerologLogger) Error(args ...any) {
	l.log.Error().Caller(1).Msg(fmt.Sprint(args...))
}

func (l *zerologLogger) Errorf(format string, args ...any) {
	l.log.Error().Caller(1).Msgf(format, args...)
}

func (l *zerologLogger) Fatal(args ...any) {
	l.log.Fatal().Caller(1).Msg(fmt.Sprint(args...))
}

func (l *zerologLogger) Fatalf(format string, args ...any) {
	l.log.Fatal().Caller(1).Msgf(format, args...)
}

func (l *zerologLogger) Info(args ...any) {
	l.log.Info().Caller(1).Msg(fmt.Sprint(args...))
}

func (l *zerologLogger) Infof(format string, args ...any) {
	l.log.Info().Caller(1).Msgf(format, args...)
}

func (l *zerologLogger) Warn(args ...any) {
	l.log.Warn().Caller(1).Msg(fmt.Sprint(args...))
}

func (l *zerologLogger) Warnf(format string, args ...any) {
	l.log.Warn().Caller(1).Msgf(format, args...)
}

func (l *zerologLogger) Debug(args ...any) {
	l.log.Debug().Caller(1).Msg(fmt.Sprint(args...))
}

func (l *zerologLogger) Debugf(format string, args ...any) {
	l.log.Debug().Caller(1).Msgf(format, args...)
}

func (l *zerologLogger) WithField(key string, value any) Logger {
	var log zerolog.Logger
	if err, ok := value.(error); ok {
		log = l.log.With().AnErr(key, err).Logger()
	} else {
		log = l.log.With().Any(key, value).Logger()
	}

	return &zerologLogger{
		log: log,
	}
}

func (l *zerologLogger) WithFields(fields map[string]any) Logger {
	logCtx := l.log.With()
	for k, v := range fields {
		if errs, ok := v.([]error); ok {
			logCtx = logCtx.Errs(k, errs)
		} else if err, ok := v.(error); ok {
			logCtx = logCtx.AnErr(k, err)
		} else {
			logCtx = logCtx.Any(k, v)
		}
	}

	return &zerologLogger{
		log: logCtx.Logger(),
	}
}
