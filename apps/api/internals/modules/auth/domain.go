package auth

// login
type LoginRequest struct {
	Username string `json:"username" validate:"required,min=2"`
	Password string `json:"password" validate:"required,min=8"`
}

// signup
type SignupRequest struct {
	Name     string `json:"name" validate:"required,min=2"`
	Username string `json:"username" validate:"required,min=2"`
	Password string `json:"password" validate:"required,min=8"`
}

// auth response
type AuthResponse struct {
	Message string           `json:"message"`
	Data    AuthResponseData `json:"data" validate:"omitempty"`
}

type AuthResponseData struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// type for the Database
type DbUser struct {
	ID       string
	Username string
	Password string
}

type CreateDbUser struct {
	Username string
	Password string
	Name     string
	ID       string
}
