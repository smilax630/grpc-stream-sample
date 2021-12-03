package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Level string

func (l Level) String() string {
	return string(l)
}

const (
	LevelDebug Level = "debug"
	LevelInfo  Level = "info"
)

var dlogger = zap.NewNop()

// SetDefaultLogger sets to package scope default logger.
func SetDefaultLogger(l *zap.Logger) {
	dlogger = l
}

// WarnIf logs error if err is not nil.
func WarnIf(err error) {
	if err == nil {
		return
	}
	dlogger.Warn("error", zap.Error(err))
}

// ErrorIf logs error if err is not nil.
func ErrorIf(err error) {
	if err == nil {
		return
	}
	dlogger.Error("error", zap.Error(err))
}

// FatalIf logs error with fatal level if err is not nil.
func FatalIf(err error) {
	if err == nil {
		return
	}
	dlogger.Fatal("error", zap.Error(err))
}

// Info logs with info level.
func Info(msg string, fields ...zap.Field) {
	dlogger.Info(msg, fields...)
}

// NewLogger returns a zap logger.
// Available level is https://github.com/uber-go/zap/blob/master/zapcore/level.go/L126.
// nolint: interfacer,errcheck
func NewLogger(level Level) (*zap.Logger, func(), error) {
	zapLevel := new(zapcore.Level)
	if err := zapLevel.Set(level.String()); err != nil {
		return nil, nil, err
	}
	logger, err := newConfig(*zapLevel).Build()
	if err != nil {
		return nil, nil, err
	}
	cleanup := func() {
		logger.Sync()
	}
	return logger, cleanup, nil
}

func newConfig(level zapcore.Level) zap.Config {
	const sampling = 100
	return zap.Config{
		Level:       zap.NewAtomicLevelAt(level),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    sampling,
			Thereafter: sampling,
		},
		Encoding:         "json",
		EncoderConfig:    newEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
}

func newEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "eventTime",
		LevelKey:       "severity",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    encodeLevel,
		EncodeTime:     zapcore.EpochTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func encodeLevel(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	switch level {
	case zapcore.DebugLevel:
		enc.AppendString("DEBUG")
	case zapcore.InfoLevel:
		enc.AppendString("INFO")
	case zapcore.WarnLevel:
		enc.AppendString("WARNING")
	case zapcore.ErrorLevel:
		enc.AppendString("ERROR")
	case zapcore.DPanicLevel:
		enc.AppendString("CRITICAL")
	case zapcore.PanicLevel:
		enc.AppendString("ALERT")
	case zapcore.FatalLevel:
		enc.AppendString("EMERGENCY")
	}
}
