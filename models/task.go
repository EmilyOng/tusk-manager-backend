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

	Tags    []Tag  `gorm:"many2many:task_tags" json:"tags"`
	UserID  string `json:"userId"`                  // Owner of the task
	BoardID string `json:"boardId"`                 // Board that the task belongs to
	StateID string `gorm:"not null" json:"stateId"` // State that the task is at
}

func (task *Task) BeforeCreate(tx *gorm.DB) (errr error) {
	// Generates a new UUID
	task.ID = uuid.NewString()
	return
}
