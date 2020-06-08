package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tsvobo/go-workshop-http/client/internal/client"
	"github.com/tsvobo/go-workshop-http/client/internal/logger"
	"github.com/tsvobo/go-workshop-http/client/internal/service"
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	ctx := context.Background()

	// TODO: Add tracing middleware (zipkinhttp)
	// TODO: Add prometheus middleware and register request_duration_seconds histogramVec
	rt := http.DefaultTransport

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
