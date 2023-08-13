package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type State struct {
	ID              string `gorm:"primaryKey" json:"id"`
	Name            string `gorm:"not null" json:"name"`
	CurrentPosition int    `gorm:"not null" json:"currentPosition"` // Sort key

	Tasks   []*Task `gorm:"not null" json:"tasks"` // Tasks belonging to the state
	BoardID *string `json:"boardId"`               // Board that the state belongs to
}

type StatePrimitive struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	CurrentPosition int    `json:"currentPosition"`

	BoardID *string `json:"boardId"`
}

// Create State
type CreateStatePayload struct {
	Name            string `json:"name"`
	BoardID         string `json:"boardId"`
	CurrentPosition int    `json:"currentPosition"`
}

type CreateStateResponse struct {
	Response
	State State `json:"data"`
}

// Update State
type UpdateStatePayload struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	BoardID         string `json:"boardId"`
	CurrentPosition int    `json:"currentPosition"`
}

type UpdateStateResponse struct {
	Response
	State StatePrimitive `json:"data"`
}

// Delete State
type DeleteStatePayload struct {
	ID string `json:"id"`
}

type DeleteStateResponse struct {
	Response
}

func (state *State) BeforeCreate(tx *gorm.DB) (errr error) {
	// Generates a new UUID
	state.ID = uuid.NewString()
	return
}
