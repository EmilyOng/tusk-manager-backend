package models

type User struct {
	ID       uint8     `gorm:"primaryKey" json:"id"`
	Name     string    `gorm:"not null" json:"name"`
	Email    string    `gorm:"not null" json:"email"`
	Password string    `gorm:"not null" json:"password"`
	Members  []*Member `json:"boardMembers"` // Boards that the user can access
	Tasks    []*Task   `json:"tasks"`        // Tasks that the user owns
}

type Profile struct {
	ID    uint8  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Auth User
type AuthUser struct {
	ID    uint8  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Token string `json:"token"`
}

type AuthUserResponse struct {
	Response
	User AuthUser `json:"data"`
}

// Login
type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Response
	User AuthUser `json:"data"`
}

// Sign Up
type SignUpPayload struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpResponse struct {
	Response
	User AuthUser `json:"data"`
}

// Get User Boards
type GetUserBoardsResponse struct {
	Response
	Boards []BoardPrimitive `json:"data"`
}
