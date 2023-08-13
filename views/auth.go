package views

type AuthUserView struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Token string `json:"token"`
}

// Authentication
type AuthUserResponse struct {
	Response
	User AuthUserView `json:"data"`
}

// Login
type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Response
	User AuthUserView `json:"data"`
}

// Sign Up
type SignUpPayload struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpResponse struct {
	Response
	User AuthUserView `json:"data"`
}
