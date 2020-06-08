package logger

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap/zapcore"
)

var messagesCounter = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "log_messages",
	Help: "The count of logged messages per log level",
}, []string{"level"})

var messagesCounterHook = func(entry zapcore.Entry) error {
	messagesCounter.WithLabelValues(entry.Level.String()).Inc()
	return nil
}
