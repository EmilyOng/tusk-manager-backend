package models

import (
	colorTypes "github.com/EmilyOng/cvwo/backend/types/color"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tag struct {
	ID    string           `gorm:"primaryKey" json:"id"`
	Name  string           `gorm:"not null" json:"name"`
	Color colorTypes.Color `gorm:"not null" json:"color" ts_type:"Color"`

	Tasks   []Task `gorm:"many2many:task_tags" json:"tasks"`
	BoardID string `json:"boardId"` // Board that the tag belongs to
}

func (tag *Tag) BeforeCreate(tx *gorm.DB) (errr error) {
	// Generates a new UUID
	tag.ID = uuid.NewString()
	return
}
