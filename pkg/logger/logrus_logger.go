package logger

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

type LogrusLogger struct {
	Log *logrus.Logger
}

func NewLogrusLogger(level int) Logger {
	log := logrus.New()
	log.SetOutput(os.Stdout)
	log.SetFormatter(&logrus.TextFormatter{
		TimestampFormat:  "15:04",
		FullTimestamp:    true,
		ForceColors:      true,
		DisableTimestamp: false,
	})
	log.SetLevel(logrus.Level(level))
	log.WithFields(logrus.Fields{
		"pid": os.Getpid(),
	})

	return &LogrusLogger{Log: log}
}

func (l *LogrusLogger) GetWriter() io.Writer {
	return l.Log.Writer()
}

func (l *LogrusLogger) Printf(format string, args ...any) {
	l.Log.Infof(format, args...)
}

func (l *LogrusLogger) Debug(args ...any) {
	l.Log.Debug(args...)
}

func (l *LogrusLogger) Debugf(format string, args ...any) {
	l.Log.Debugf(format, args...)
}

func (l *LogrusLogger) Info(args ...any) {
	l.Log.Info(args...)
}

func (l *LogrusLogger) Infof(format string, args ...any) {
	l.Log.Infof(format, args...)
}

func (l *LogrusLogger) Warn(args ...any) {
	l.Log.Warn(args...)
}

func (l *LogrusLogger) Warnf(format string, args ...any) {
	l.Log.Warnf(format, args...)
}

func (l *LogrusLogger) Error(args ...any) {
	l.Log.Error(args...)
}

func (l *LogrusLogger) Errorf(format string, args ...any) {
	l.Log.Errorf(format, args...)
}

func (l *LogrusLogger) Fatal(args ...any) {
	l.Log.Fatal(args...)
}

func (l *LogrusLogger) Fatalf(format string, args ...any) {
	l.Log.Fatalf(format, args...)
}

func (l *LogrusLogger) WithField(key string, value any) Logger {
	return &LogrusEntry{
		entry: l.Log.WithField(key, value),
	}
}
func (l *LogrusLogger) WithFields(fields map[string]any) Logger {
	return &LogrusEntry{
		entry: l.Log.WithFields(fields),
	}
}

type LogrusEntry struct {
	entry *logrus.Entry
}

func (l *LogrusEntry) GetWriter() io.Writer {
	return l.entry.Writer()
}

func (l *LogrusEntry) Printf(format string, args ...any) {
	l.entry.Printf(format, args...)
}

func (l *LogrusEntry) Error(args ...any) {
	l.entry.Error(args...)
}

func (l *LogrusEntry) Errorf(format string, args ...any) {
	l.entry.Errorf(format, args...)
}

func (l *LogrusEntry) Fatal(args ...any) {
	l.entry.Fatal(args...)
}

func (l *LogrusEntry) Fatalf(format string, args ...any) {
	l.entry.Fatalf(format, args...)
}

func (l *LogrusEntry) Info(args ...any) {
	l.entry.Info(args...)
}

func (l *LogrusEntry) Infof(format string, args ...any) {
	l.entry.Infof(format, args...)
}

func (l *LogrusEntry) Warn(args ...any) {
	l.entry.Warn(args...)
}

func (l *LogrusEntry) Warnf(format string, args ...any) {
	l.entry.Warnf(format, args...)
}

func (l *LogrusEntry) Debug(args ...any) {
	l.entry.Debug(args...)
}

func (l *LogrusEntry) Debugf(format string, args ...any) {
	l.entry.Debugf(format, args...)
}

func (l *LogrusEntry) WithField(key string, value any) (entry Logger) {
	entry = &LogrusEntry{l.entry.WithField(key, value)}
	return
}

func (l *LogrusEntry) WithFields(args map[string]any) (entry Logger) {
	entry = &LogrusEntry{l.entry.WithFields(args)}
	return
}
