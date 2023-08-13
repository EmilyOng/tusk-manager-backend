package models

import (
	roleTypes "github.com/EmilyOng/cvwo/backend/types/role"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Member struct {
	ID   string         `gorm:"primary_key" json:"id"`
	Role roleTypes.Role `gorm:"not null" json:"role" ts_type:"Role"`

	UserID  string `json:"userId"` // User ID of the board member
	User    *User  `json:"user"`
	BoardID string `json:"boardId"` // Board that the member belongs to
}

func (member *Member) BeforeCreate(tx *gorm.DB) (err error) {
	if len(member.ID) > 0 {
		return
	}
	// Generates a new UUID
	member.ID = uuid.NewString()
	return
}
