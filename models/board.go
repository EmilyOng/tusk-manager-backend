package models

import (
	colorTypes "github.com/EmilyOng/cvwo/backend/types/color"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Board struct {
	ID    string           `gorm:"primaryKey" json:"id"`
	Name  string           `gorm:"not null" json:"name"`
	Color colorTypes.Color `gorm:"not null" json:"color" ts_type:"Color"`

	Tasks   []*Task   `gorm:"not null" json:"tasks"`  // Tasks belonging to the board
	Tags    []*Tag    `gorm:"not null" json:"tags"`   // Tags belonging to the board
	States  []*State  `gorm:"not null" json:"states"` // States belonging to the board
	Members []*Member `json:"boardMembers"`           // Members belonging to the board
}

func (board *Board) BeforeCreate(tx *gorm.DB) (err error) {
	if len(board.ID) > 0 {
		return
	}
	// Generates a new UUID
	board.ID = uuid.NewString()
	return
}
