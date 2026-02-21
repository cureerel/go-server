package logger

import (
	"go.uber.org/zap"
)

type Logger interface {
	Info(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Debug(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Fatal(msg string, fields ...Field) // ADD THIS
	With(fields ...Field) Logger
}

type Field struct {
	Key   string
	Value interface{}
}

type zapLoggerAdapter struct {
	logger *zap.Logger
}

func (z *zapLoggerAdapter) Info(msg string, fields ...Field) {
	z.logger.Info(msg, convertFields(fields)...)
}

func (z *zapLoggerAdapter) Error(msg string, fields ...Field) {
	z.logger.Error(msg, convertFields(fields)...)
}

func (z *zapLoggerAdapter) Debug(msg string, fields ...Field) {
	z.logger.Debug(msg, convertFields(fields)...)
}

func (z *zapLoggerAdapter) Warn(msg string, fields ...Field) {
	z.logger.Warn(msg, convertFields(fields)...)
}

func (z *zapLoggerAdapter) Fatal(msg string, fields ...Field) { // ADD THIS
	z.logger.Fatal(msg, convertFields(fields)...)
}

func (z *zapLoggerAdapter) With(fields ...Field) Logger {
	return &zapLoggerAdapter{logger: z.logger.With(convertFields(fields)...)}
}

func convertFields(fields []Field) []zap.Field {
	zapFields := make([]zap.Field, len(fields))
	for i, f := range fields {
		zapFields[i] = zap.Any(f.Key, f.Value)
	}
	return zapFields
}

// FIX: No config needed, use defaults
func New() Logger {
	logger, _ := zap.NewProduction()
	return &zapLoggerAdapter{logger: logger}
}

// Optional: New with config
func NewWithConfig(production bool) Logger {
	if production {
		logger, _ := zap.NewProduction()
		return &zapLoggerAdapter{logger: logger}
	}
	logger, _ := zap.NewDevelopment()
	return &zapLoggerAdapter{logger: logger}
}