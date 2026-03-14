// pkg/logger/logger.go
package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Field is a structured log key-value pair.
type Field struct {
	Key   string
	Value any
}

// Logger is the application-wide logging interface.
type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
}

// New returns a coloured, human-readable dev logger.
//
// Example line:
//
//	15:04:05.000  INF  starting  env=development  port=8080
func New() Logger { return build(false) }

// NewProduction / NewProd return a structured JSON logger for production.
// Both names exist so callers may use either convention.
func NewProduction() Logger { return build(true) }
func NewProd() Logger       { return build(true) }

// ── builder ──────────────────────────────────────────────────

func build(prod bool) Logger {
	level := zapcore.DebugLevel

	var enc zapcore.Encoder
	if prod {
		cfg := zap.NewProductionEncoderConfig()
		cfg.TimeKey = "ts"
		cfg.MessageKey = "msg"
		cfg.EncodeTime = zapcore.ISO8601TimeEncoder
		enc = zapcore.NewJSONEncoder(cfg)
		level = zapcore.InfoLevel
	} else {
		enc = zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
			TimeKey:        "T",
			LevelKey:       "L",
			MessageKey:     "M",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    colorLevel,
			EncodeTime:     compactTime,
			EncodeDuration: zapcore.StringDurationEncoder,
			// CallerKey / EncodeCaller intentionally omitted to keep dev output
			// uncluttered. Re-add if you want file:line in output.
		})
	}

	core := zapcore.NewCore(enc, zapcore.AddSync(os.Stdout), level)
	return &zapLogger{z: zap.New(core, zap.WithCaller(false))}
}

// compactTime prints HH:MM:SS.mmm — no date, no timezone.
func compactTime(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("15:04:05.000"))
}

// colorLevel prints a fixed-width, ANSI-coloured level tag.
//
//	DBG → cyan   INF → green   WRN → yellow   ERR → red   FTL → magenta
func colorLevel(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	const reset = "\033[0m"
	var badge string
	switch l {
	case zapcore.DebugLevel:
		badge = "\033[36mDBG" + reset // cyan
	case zapcore.InfoLevel:
		badge = "\033[32mINF" + reset // green
	case zapcore.WarnLevel:
		badge = "\033[33mWRN" + reset // yellow
	case zapcore.ErrorLevel:
		badge = "\033[31mERR" + reset // red
	case zapcore.FatalLevel:
		badge = "\033[35mFTL" + reset // magenta
	default:
		badge = l.CapitalString()
	}
	enc.AppendString(badge)
}

// ── zapLogger ────────────────────────────────────────────────

type zapLogger struct{ z *zap.Logger }

func toZap(ff []Field) []zap.Field {
	out := make([]zap.Field, len(ff))
	for i, f := range ff {
		out[i] = zap.Any(f.Key, f.Value)
	}
	return out
}

func (l *zapLogger) Debug(msg string, ff ...Field) { l.z.Debug(msg, toZap(ff)...) }
func (l *zapLogger) Info(msg string, ff ...Field)  { l.z.Info(msg, toZap(ff)...) }
func (l *zapLogger) Warn(msg string, ff ...Field)  { l.z.Warn(msg, toZap(ff)...) }
func (l *zapLogger) Error(msg string, ff ...Field) { l.z.Error(msg, toZap(ff)...) }
func (l *zapLogger) Fatal(msg string, ff ...Field) { l.z.Fatal(msg, toZap(ff)...) }