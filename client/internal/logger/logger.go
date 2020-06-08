package logger

import (
	"context"

	"github.com/openzipkin/zipkin-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.SugaredLogger
}

type Options struct {
	Level          string
	MetricsEnabled bool
}

// Global configuration for logger package.
var (
	Log            Logger
	DefaultOptions = Options{
		Level:          "debug",
		MetricsEnabled: true,
	}
)

func init() {
	MustSetup(DefaultOptions)
}

func MustSetup(opts Options) {
	l, err := New(opts)
	if err != nil {
		panic(err)
	}
	Log = *l
}

func New(opts Options) (*Logger, error) {
	c := zap.NewProductionConfig()
	c.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	c.Encoding = "console"

	var l zapcore.Level
	if err := l.Set(opts.Level); err != nil {
		return nil, err
	}
	c.Level.SetLevel(l)

	var o []zap.Option
	if opts.MetricsEnabled {
		o = append(o, zap.Hooks(messagesCounterHook))
	}
	logger, err := c.Build(o...)
	if err != nil {
		return nil, err
	}

	return &Logger{logger.Sugar()}, nil
}

func (l *Logger) WithContext(ctx context.Context) *Logger {
	var fields []interface{}

	if span := zipkin.SpanFromContext(ctx); span != nil {
		traceID := span.Context().TraceID.String()
		spanID := span.Context().ID.String()
		fields = append(fields, zap.String("traceId", traceID), zap.String("spanId", spanID))
	}

	if len(fields) > 0 {
		return &Logger{l.With(fields...)}
	}

	return l
}
