package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID       string `gorm:"primaryKey" json:"id"`
	Name     string `gorm:"not null" json:"name"`
	Email    string `gorm:"not null" json:"email"`
	Password string `gorm:"not null" json:"password"`

	Members []*Member `json:"boardMembers"` // Boards that the user can access
	Tasks   []*Task   `json:"tasks"`        // Tasks that the user owns
}

func (user *User) BeforeCreate(tx *gorm.DB) (errr error) {
	// Generates a new UUID
	user.ID = uuid.NewString()
	return
}
