package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/tsvobo/go-workshop-http/server/internal/logger"
	"github.com/tsvobo/go-workshop-http/server/internal/model"
	"github.com/tsvobo/go-workshop-http/server/internal/service"
	"go.uber.org/zap"
)

type todoRetrieval interface {
	Find(ctx context.Context, id string) (model.Task, error)
}

type TodoRetrievalHandler struct {
	Retrieval todoRetrieval
}

func (h *TodoRetrievalHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")

	entry, err := h.Retrieval.Find(req.Context(), id)

	if err == service.ErrNotFound {
		writer.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		logger.Log.WithContext(req.Context()).With(zap.Error(err)).Errorf("Error while getting todo by ID: '%x'.", id)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Add("Content-Type", "application/json")
	if err = json.NewEncoder(writer).Encode(entry); err != nil {
		logger.Log.WithContext(req.Context()).With(zap.Error(err)).Errorf("Error while encoding todo to JSON: '%v'.", entry)
		writer.WriteHeader(http.StatusInternalServerError)
	}
}
