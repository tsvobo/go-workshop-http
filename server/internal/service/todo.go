package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/tsvobo/go-workshop-http/server/internal/logger"
	"github.com/tsvobo/go-workshop-http/server/internal/model"
)

var ErrNotFound = errors.New("not found")

type Task struct {
	data map[string]model.Task
}

func NewTask() *Task {
	d := make(map[string]model.Task)
	return &Task{data: d}
}

func (s *Task) Create(ctx context.Context, title string, note *string, dueDate time.Time) model.Task {
	id := uuid.New().String()
	entry := model.Task{
		ID:      id,
		Title:   title,
		Note:    note,
		DueDate: dueDate,
	}

	s.data[id] = entry

	logger.Log.WithContext(ctx).Infof("Task '%v' created successfully.", entry)
	return entry
}

func (s *Task) Find(ctx context.Context, id string) (model.Task, error) {
	logger.Log.WithContext(ctx).Debugf("Searching for todo with id '%v'.", id)

	if task, ok := s.data[id]; ok {
		return task, nil
	}
	return model.Task{}, ErrNotFound
}
