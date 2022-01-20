package models

type Color string

const (
	Turquoise Color = "Turquoise"
	Blue      Color = "Blue"
	Cyan      Color = "Cyan"
	Green     Color = "Green"
	Yellow    Color = "Yellow"
	Red       Color = "Red"
)

type Role string

const (
	Owner  Role = "Owner"
	Editor Role = "Editor"
	Viewer Role = "Viewer"
)

type ErrorCode string

const (
	NotFound          ErrorCode = "not_found"
	ServerError       ErrorCode = "server_error"
	UnauthorizedError ErrorCode = "unauthorized"
	TypeMismatch      ErrorCode = "type_mismatch"
	ConflictError     ErrorCode = "conflict"
)

type Response struct {
	Error string `json:"error" ts_type:"ErrorCode" ts_transform:"(__VALUE__).toString()"`
}
