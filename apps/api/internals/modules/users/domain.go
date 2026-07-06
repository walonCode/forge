package users

import "time"

// user profile returned to clients (never includes the password)
type UserProfile struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

// update profile request (both fields optional; at least one required)
type UpdateProfileRequest struct {
	Name     string `json:"name" validate:"omitempty,min=2"`
	Username string `json:"username" validate:"omitempty,min=2"`
}

// update password request
type UpdatePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required,min=8"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
}

// generic user response envelope
type UserResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}
