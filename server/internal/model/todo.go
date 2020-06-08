package model

import (
	"fmt"
	"time"
)

type Task struct {
	ID      string    `json:"id,omitempty" validate:"required"`
	Title   string    `json:"title" validate:"required"`
	Note    *string   `json:"note"`
	DueDate time.Time `json:"due_date" validate:"required"`
}

func (t Task) String() string {
	note := "<nil>"

	if t.Note != nil {
		note = *t.Note
	}

	return fmt.Sprintf("{ ID: '%v', Title: '%v', Note: '%v', DueDate: '%v'}", t.ID, t.Title, note, t.DueDate)
}
