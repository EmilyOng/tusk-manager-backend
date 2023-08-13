package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tag struct {
	ID    string `gorm:"primaryKey" json:"id"`
	Name  string `gorm:"not null" json:"name"`
	Color Color  `gorm:"not null" json:"color" ts_type:"Color"`

	Tasks   []*Task `gorm:"many2many:task_tags" json:"tasks"`
	BoardID *string `json:"boardId"` // Board that the tag belongs to
}

type TagPrimitive struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color Color  `json:"color" ts_type:"Color"`

	BoardID *string `json:"boardId"`
}

// Create Tag
type CreateTagPayload struct {
	Name    string `json:"name"`
	Color   Color  `json:"color" ts_type:"Color"`
	BoardID string `json:"boardId"`
}

type CreateTagResponse struct {
	Response
	Tag TagPrimitive `json:"data"`
}

// Update Tag
type UpdateTagPayload struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	BoardID string `json:"boardId"`
	Color   Color  `json:"color" ts_type:"Color"`
}

type UpdateTagResponse struct {
	Response
	Tag TagPrimitive `json:"data"`
}

// Delete Tag
type DeleteTagPayload struct {
	ID string `json:"id"`
}

type DeleteTagResponse struct {
	Response
}

func (tag *Tag) BeforeCreate(tx *gorm.DB) (errr error) {
	// Generates a new UUID
	tag.ID = uuid.NewString()
	return
}
