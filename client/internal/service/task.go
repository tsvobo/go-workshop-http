package service

import (
	"context"
	"time"

	"github.com/tsvobo/go-workshop-http/client/internal/logger"
	"github.com/tsvobo/go-workshop-http/client/internal/model"
	"github.com/tsvobo/go-workshop-http/client/internal/trace"
	"go.uber.org/zap"
)

type taskCreator interface {
	Create(ctx context.Context, task model.Task) (model.Task, error)
}

type taskFinder interface {
	Find(ctx context.Context, id string) (model.Task, error)
}

type Task struct {
	Creator taskCreator
	Finder  taskFinder
}

func (s *Task) Create(ctx context.Context, title string, note *string, dueDate time.Time) (model.Task, error) {
	span, ctx := trace.Tracer.StartSpanFromContext(ctx, "create task")
	defer span.Finish()

	entry := model.Task{Title: title, Note: note, DueDate: dueDate}

	t, err := s.Creator.Create(ctx, entry)
	if err != nil {
		logger.Log.WithContext(ctx).With(zap.Error(err)).Errorf("Failed to create task '%s'.", entry)
	} else {
		logger.Log.WithContext(ctx).Infof("Task '%v' created successfully.", t)
	}

	return t, err
}

func (s *Task) Find(ctx context.Context, id string) (model.Task, error) {
	span, ctx := trace.Tracer.StartSpanFromContext(ctx, "find task")
	defer span.Finish()

	logger.Log.WithContext(ctx).Debugf("Searching for task with id '%v'.", id)

	return s.Finder.Find(ctx, id)
}
