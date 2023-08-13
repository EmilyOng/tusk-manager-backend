package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Member struct {
	ID   string `gorm:"primary_key" json:"id"`
	Role Role   `gorm:"not null" json:"role" ts_type:"Role"`

	UserID  *string `json:"userId"` // User ID of the board member
	User    *User   `json:"user"`
	BoardID *string `json:"boardId"` // Board that the member belongs to
}

type MemberPrimitive struct {
	ID   string `json:"id"`
	Role Role   `json:"role" ts_type:"Role"`

	UserID  *string `json:"userId"`  // User ID of the board member
	BoardID *string `json:"boardId"` // Board that the member belongs to
}

type MemberProfile struct {
	ID      string  `json:"id"`
	Role    Role    `json:"role" ts_type:"Role"`
	Profile Profile `json:"profile"`
}

// Create Member
type CreateMemberPayload struct {
	Role    Role   `json:"role" ts_type:"Role"`
	Email   string `json:"email"` // Invitation is by email
	BoardID string `json:"boardId"`
}

type CreateMemberResponse struct {
	Response
	MemberProfile MemberProfile `json:"data"`
}

// Update Member
type UpdateMemberPayload struct {
	ID   string `json:"id"`
	Role Role   `json:"role" ts_type:"Role"`
}

type UpdateMemberResponse struct {
	Response
	MemberProfile MemberProfile `json:"data"`
}

// Delete Member
type DeleteMemberPayload struct {
	ID string `json:"id"`
}

type DeleteMemberResponse struct {
	Response
}

func (member *Member) BeforeCreate(tx *gorm.DB) (errr error) {
	// Generates a new UUID
	member.ID = uuid.NewString()
	return
}
