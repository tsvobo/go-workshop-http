package model

import (
	"fmt"
	"time"
)

// TODO TASK-1: Add json tags - ID and note are optional
// TODO TASK-2: Add validation tags - all but note are required
type Task struct {
	ID      string    `json:"id,omitempty" validate:"required"`
	Title   string    `json:"title" validate:"required"`
	Note    *string   `json:"note,omitempty" validate:"-"`
	DueDate time.Time `json:"due_date" validate:"required"`
}

func (t Task) String() string {
	note := "<nil>"

	if t.Note != nil {
		note = *t.Note
	}

	return fmt.Sprintf("{ ID: '%v', Title: '%v', Note: '%v', DueDate: '%v'}", t.ID, t.Title, note, t.DueDate)
}
