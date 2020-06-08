package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/tsvobo/go-workshop-http/server/internal/logger"
	"github.com/tsvobo/go-workshop-http/server/internal/model"
	"go.uber.org/zap"
)

type todoCreator interface {
	Create(ctx context.Context, title string, note *string, dueDate time.Time) model.Task
}

type todoRequest struct {
	Title   string    `json:"title" validate:"required"`
	Note    *string   `json:"note"`
	DueDate time.Time `json:"due_date" validate:"required"`
}

type TodoCreatorHandler struct {
	Creator todoCreator
}

func (h *TodoCreatorHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	var payload todoRequest

	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		logger.Log.WithContext(req.Context()).With(zap.Error(err)).Warn("Could not parse request body.")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := validate.Struct(payload); err != nil {
		logger.Log.WithContext(req.Context()).With(zap.Error(err)).Warnf("Request payload is not valid '%s'.", payload)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	entity := h.Creator.Create(req.Context(), payload.Title, payload.Note, payload.DueDate)

	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(writer).Encode(entity); err != nil {
		logger.Log.WithContext(req.Context()).With(zap.Error(err)).Error("Error while encoding todo to JSON: '%s'.", entity)
		writer.WriteHeader(http.StatusInternalServerError)
	}
}
