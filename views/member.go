package views

import (
	roleTypes "github.com/EmilyOng/tusk-manager/backend/types/role"
)

type MemberMinimalView struct {
	ID   string         `json:"id"`
	Role roleTypes.Role `json:"role" ts_type:"Role"`

	UserID  string `json:"userId"`  // User ID of the board member
	BoardID string `json:"boardId"` // Board that the member belongs to
}

type MemberFullView struct {
	ID   string          `json:"id"`
	Role roleTypes.Role  `json:"role" ts_type:"Role"`
	User UserMinimalView `json:"user"`
}

// Create Member
type CreateMemberPayload struct {
	Role    roleTypes.Role `json:"role" ts_type:"Role"`
	Email   string         `json:"email"` // Invitation is by email
	BoardID string         `json:"boardId"`
}

type CreateMemberResponse struct {
	Response
	Member MemberFullView `json:"data"`
}

// Update Member
type UpdateMemberPayload struct {
	ID   string         `json:"id"`
	Role roleTypes.Role `json:"role" ts_type:"Role"`
}

type UpdateMemberResponse struct {
	Response
	Member MemberFullView `json:"data"`
}

// Delete Member
type DeleteMemberPayload struct {
	ID string `json:"id"`
}

type DeleteMemberResponse struct {
	Response
}
