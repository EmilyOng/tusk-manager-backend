package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type State struct {
	ID              string `gorm:"primaryKey" json:"id"`
	Name            string `gorm:"not null" json:"name"`
	CurrentPosition int    `gorm:"not null" json:"currentPosition"` // Sort key

	Tasks   []Task `gorm:"not null" json:"tasks"` // Tasks belonging to the state
	BoardID string `json:"boardId"`               // Board that the state belongs to
}

func (state *State) BeforeCreate(tx *gorm.DB) (errr error) {
	// Generates a new UUID
	state.ID = uuid.NewString()
	return
}
