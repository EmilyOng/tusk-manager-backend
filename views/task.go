package views

import (
	"time"

	"github.com/EmilyOng/cvwo/backend/models"
)

type TaskMinimalView struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	DueAt       *time.Time `json:"dueAt" ts_type:"Date" ts_transform:"new Date(__VALUE__)"`
}

type TaskFullView = models.Task

// Create Task
type CreateTaskPayload struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	DueAt       string `json:"dueAt,omitempty" ts_type:"Date" ts_transform:"new Date(__VALUE__)"`

	StateID string           `json:"stateId"`
	Tags    []TagMinimalView `json:"tags"`
	BoardID string           `json:"boardId"`
	UserID  string           `json:"userId"`
}

type CreateTaskResponse struct {
	Response
	Task TaskFullView `json:"data"`
}

// Update Task
type UpdateTaskPayload struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	DueAt       string `json:"dueAt,omitempty" ts_type:"Date" ts_transform:"new Date(__VALUE__)"`

	StateID string           `json:"stateId"`
	Tags    []TagMinimalView `json:"tags"`
	BoardID string           `json:"boardId"`
	UserID  string           `json:"userId"`
}

type UpdateTaskResponse struct {
	Response
	Task TaskFullView `json:"data"`
}

// Delete Task
type DeleteTaskPayload struct {
	ID string `json:"id"`
}

type DeleteTaskResponse struct {
	Response
}
