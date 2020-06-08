package model

import (
	"fmt"
	"time"
)

// TODO: Add json tags
// TODO: Add validation tags
type Task struct {
	ID      string
	Title   string
	Note    *string
	DueDate time.Time
}

func (t Task) String() string {
	note := "<nil>"

	if t.Note != nil {
		note = *t.Note
	}

	return fmt.Sprintf("{ ID: '%v', Title: '%v', Note: '%v', DueDate: '%v'}", t.ID, t.Title, note, t.DueDate)
}
