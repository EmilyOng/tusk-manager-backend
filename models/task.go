package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Task struct {
	ID          string     `gorm:"primary_key" json:"id"`
	Name        string     `gorm:"not null" json:"name"`
	Description string     `gorm:"default:''" json:"description"`
	DueAt       *time.Time `json:"dueAt" ts_type:"Date" ts_transform:"new Date(__VALUE__)"`

	Tags    []*Tag  `gorm:"many2many:task_tags" json:"tags"`
	UserID  *string `json:"userId"`                  // Owner of the task
	BoardID *string `json:"boardId"`                 // Board that the task belongs to
	StateID *string `gorm:"not null" json:"stateId"` // State that the task is at
}

type TaskPrimitive struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	DueAt       *time.Time `json:"dueAt" ts_type:"Date" ts_transform:"new Date(__VALUE__)"`

	UserID  *string `json:"userId"`
	BoardID *string `json:"boardId"`
	StateID *string `gorm:"not null" json:"stateId"`
}

// Create Task
type CreateTaskPayload struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	DueAt       string `json:"dueAt,omitempty" ts_type:"Date" ts_transform:"new Date(__VALUE__)"`

	StateID *string         `json:"stateId"`
	Tags    []*TagPrimitive `json:"tags"`
	BoardID *string         `json:"boardId"`
	UserID  *string         `json:"userId"`
}

type CreateTaskResponse struct {
	Response
	Task Task `json:"data"`
}

// Update Task
type UpdateTaskPayload struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	DueAt       string `json:"dueAt,omitempty" ts_type:"Date" ts_transform:"new Date(__VALUE__)"`

	StateID string         `json:"stateId"`
	Tags    []TagPrimitive `json:"tags"`
	BoardID string         `json:"boardId"`
	UserID  string         `json:"userId"`
}

type UpdateTaskResponse struct {
	Response
	Task Task `json:"data"`
}

// Delete Task
type DeleteTaskPayload struct {
	ID string `json:"id"`
}

type DeleteTaskResponse struct {
	Response
}

func (task *Task) BeforeCreate(tx *gorm.DB) (errr error) {
	// Generates a new UUID
	task.ID = uuid.NewString()
	return
}
