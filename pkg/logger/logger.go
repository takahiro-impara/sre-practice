package logger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger provides structured logging capabilities
type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Warn(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
	Debug(ctx context.Context, msg string, fields ...zap.Field)
	WithFields(fields ...zap.Field) Logger
}

type logger struct {
	*zap.Logger
}

// NewLogger creates a new logger instance based on environment
func NewLogger(env string) (Logger, error) {
	var config zap.Config
	if env == "production" {
		config = NewProductionConfig()
	} else {
		config = NewDevelopmentConfig()
	}

	zapLogger, err := config.Build(
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zap.ErrorLevel),
	)
	if err != nil {
		return nil, err
	}

	return &logger{Logger: zapLogger}, nil
}

// NewProductionConfig returns production logging configuration
func NewProductionConfig() zap.Config {
	config := zap.NewProductionConfig()
	config.Sampling = &zap.SamplingConfig{
		Initial:    100,
		Thereafter: 100,
	}
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	return config
}

// NewDevelopmentConfig returns development logging configuration
func NewDevelopmentConfig() zap.Config {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return config
}

func (l *logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	l.Logger.Info(msg, append(fields, extractTraceID(ctx)...)...)
}

func (l *logger) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	l.Logger.Warn(msg, append(fields, extractTraceID(ctx)...)...)
}

func (l *logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	l.Logger.Error(msg, append(fields, extractTraceID(ctx)...)...)
}

func (l *logger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	l.Logger.Debug(msg, append(fields, extractTraceID(ctx)...)...)
}

func (l *logger) WithFields(fields ...zap.Field) Logger {
	return &logger{Logger: l.Logger.With(fields...)}
}

// extractTraceID extracts trace ID from context if present
func extractTraceID(ctx context.Context) []zap.Field {
	// TODO: Extract trace ID from context when OpenTelemetry is integrated
	return nil
}