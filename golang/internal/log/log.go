package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LogLevel is integer enumaration, in order to allow different
// levels of logging through the code, however only the configured
// level will be logged
type LogLevel int

const (
	Debug LogLevel = 2
	Info  LogLevel = 1
	Error LogLevel = 0
)

type LogEncoder string

const (
	LogJSONEncoder    LogEncoder = "json"
	LogConsoleEncoder LogEncoder = "console"
)

// struct to confiigure the logger
type LogOpts struct {
	Name        string
	Level       LogLevel
	Encoder     LogEncoder
	TimeEncoder zapcore.TimeEncoder
}

type Log struct {
	*zap.Logger
	level LogLevel
}

// the logger interface useful while mocking
type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
}

// constructor to crate a new logger
func New(opts *LogOpts) (*Log, error) {
	conf := zap.NewProductionConfig()
	lvl := zapcore.InfoLevel
	switch opts.Level {
	case Debug:
		lvl = zapcore.DebugLevel
	case Info:
		lvl = zapcore.InfoLevel
	case Error:
		lvl = zapcore.ErrorLevel
	}
	conf.Level = zap.NewAtomicLevelAt(lvl)

	switch opts.Encoder {
	case LogJSONEncoder:
		conf.Encoding = string(LogJSONEncoder)
	case LogConsoleEncoder:
		conf.Encoding = string(LogConsoleEncoder)
	}

	switch opts.TimeEncoder {
	case nil:
		conf.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	default:
		conf.EncoderConfig.EncodeTime = opts.TimeEncoder
	}

	logger, err := conf.Build()
	if err != nil {
		return nil, err
	}

	namedLogger := logger.Named(opts.Name)
	namedLogger = namedLogger.WithOptions(zap.AddCallerSkip(1))

	l := Log{
		Logger: namedLogger,
		level:  opts.Level,
	}

	return &l, nil

}

func (l *Log) Debug(msg string, fields ...zap.Field) {
	if l.level >= Debug {
		l.Logger.Debug(msg, fields...)
	}
}

func (l *Log) Info(msg string, fields ...zap.Field) {
	if l.level >= Info {
		l.Logger.Info(msg, fields...)
	}
}
func (l *Log) Error(msg string, fields ...zap.Field) {
	if l.level >= Error {
		l.Logger.Error(msg, fields...)
	}
}
