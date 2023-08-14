package models

import (
	colorTypes "github.com/EmilyOng/tusk-manager/backend/types/color"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tag struct {
	ID    string           `gorm:"primaryKey" json:"id"`
	Name  string           `gorm:"not null" json:"name"`
	Color colorTypes.Color `gorm:"not null" json:"color" ts_type:"Color"`

	Tasks   []*Task `gorm:"many2many:task_tags" json:"tasks"`
	BoardID string  `json:"boardId"` // Board that the tag belongs to
}

func (tag *Tag) BeforeCreate(tx *gorm.DB) (err error) {
	if len(tag.ID) > 0 {
		return
	}
	// Generates a new UUID
	tag.ID = uuid.NewString()
	return
}
