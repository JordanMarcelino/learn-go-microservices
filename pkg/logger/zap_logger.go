package logger

import (
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	Log *zap.Logger
}

func NewZapLogger(level int) Logger {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("15:04"),
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		zapcore.Level(level),
	)

	logger := zap.New(core,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.Fields(
			zap.Int("pid", os.Getpid()),
		),
	)

	return &ZapLogger{Log: logger}
}

func (l *ZapLogger) GetWriter() io.Writer {
	return zap.NewStdLog(l.Log).Writer()
}

func (l *ZapLogger) Printf(format string, args ...any) {
	l.Log.Sugar().Infof(format, args...)
}

func (l *ZapLogger) Error(args ...any) {
	l.Log.Sugar().Error(args...)
}

func (l *ZapLogger) Errorf(format string, args ...any) {
	l.Log.Sugar().Errorf(format, args...)
}

func (l *ZapLogger) Fatal(args ...any) {
	l.Log.Sugar().Fatal(args...)
}

func (l *ZapLogger) Fatalf(format string, args ...any) {
	l.Log.Sugar().Fatalf(format, args...)
}

func (l *ZapLogger) Info(args ...any) {
	l.Log.Sugar().Info(args...)
}

func (l *ZapLogger) Infof(format string, args ...any) {
	l.Log.Sugar().Infof(format, args...)
}

func (l *ZapLogger) Warn(args ...any) {
	l.Log.Sugar().Warn(args...)
}

func (l *ZapLogger) Warnf(format string, args ...any) {
	l.Log.Sugar().Warnf(format, args...)
}

func (l *ZapLogger) Debug(args ...any) {
	l.Log.Sugar().Debug(args...)
}

func (l *ZapLogger) Debugf(format string, args ...any) {
	l.Log.Sugar().Debugf(format, args...)
}

func (l *ZapLogger) WithField(key string, value any) Logger {
	if err, ok := value.(error); ok {
		return &ZapLogger{Log: l.Log.With(zap.Error(err))}
	}
	return &ZapLogger{Log: l.Log.With(zap.Any(key, value))}
}

func (l *ZapLogger) WithFields(fields map[string]any) Logger {
	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		if errs, ok := v.([]error); ok {
			zapFields = append(zapFields, zap.Errors(k, errs))
		} else if err, ok := v.(error); ok {
			zapFields = append(zapFields, zap.Error(err))
		} else {
			zapFields = append(zapFields, zap.Any(k, v))
		}
	}
	return &ZapLogger{Log: l.Log.With(zapFields...)}
}
