package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tsvobo/go-workshop-http/client/internal/client"
	"github.com/tsvobo/go-workshop-http/client/internal/logger"
	"github.com/tsvobo/go-workshop-http/client/internal/metrics"
	"github.com/tsvobo/go-workshop-http/client/internal/service"
	"github.com/tsvobo/go-workshop-http/client/internal/trace"
)

var (
	histVec = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "request_duration_seconds",
			Help:    "A histogram of request durations.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"destination", "endpoint", "code"},
	)
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	ctx := context.Background()

	// TODO TASK-6: Add tracing middleware (zipkinhttp)
	zrt, err := zipkinhttp.NewTransport(
		trace.Tracer,
		zipkinhttp.TransportTrace(true),
	)
	// TODO TASK-5.1: Add prometheus middleware and register histogramVec
	rt := metrics.InstrumentRoundTripperDuration(
		histVec.MustCurryWith(prometheus.Labels{"destination": "task-service"}),
		zrt,
	)

	// Create Task (HTTP) client
	taskClient, err := client.NewTask("http://127.0.0.1:8080", &http.Client{Timeout: 30 * time.Second, Transport: rt})
	if err != nil {
		logger.Log.Fatal(err)
	}

	// Create service for create and get tasks
	taskServ := service.Task{Creator: &taskClient, Finder: &taskClient}

	// Perform actions on task service

	n := "This task is important"
	task, err := taskServ.Create(ctx, "My task", &n, time.Now().Add(5*24*time.Hour))
	if err != nil {
		panic(err)
	}
	_, _ = taskServ.Find(ctx, task.ID)

	_, _ = taskServ.Find(ctx, "some-fake-id")

	// Expose /metrics endpoint
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		_ = http.ListenAndServe(":8079", nil)
	}()

	// Block until we receive the signal.
	<-sigs
	logger.Log.Info("Finished.")
}
