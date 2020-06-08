package main

import (
	"net/http"

	"github.com/go-chi/chi"
	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tsvobo/go-workshop-http/server/internal/handler"
	"github.com/tsvobo/go-workshop-http/server/internal/logger"
	"github.com/tsvobo/go-workshop-http/server/internal/service"
	"github.com/tsvobo/go-workshop-http/server/internal/trace"
	"go.uber.org/zap"
)

func main() {
	tracer, _ := trace.NewTracer()

	mux := chi.NewRouter()
	mux.Use(zipkinhttp.NewServerMiddleware(tracer))

	s := service.NewTask()

	mux.Route("/v1/tasks", func(r chi.Router) {
		r.Get("/{id}", (&handler.TodoRetrievalHandler{Retrieval: s}).ServeHTTP)
		r.Post("/", (&handler.TodoCreatorHandler{Creator: s}).ServeHTTP)
	})
	mux.Handle("/metrics", promhttp.Handler())

	srv := &http.Server{Addr: ":8080", Handler: mux}

	if err := srv.ListenAndServe(); err != nil {
		logger.Log.With(zap.Error(err)).Fatal("Failed to start server.")
	}
}
