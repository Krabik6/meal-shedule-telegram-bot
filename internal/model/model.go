package model

type SignUpCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name" binding:"required"`
}

type SignInCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

// Create recipe credentials
type CreateRecipeInput struct {
	Title         string  `json:"title"`
	Description   string  `json:"description"`
	IsPublic      bool    `json:"isPublic" db:"public"`
	Cost          float64 `json:"cost,omitempty"`
	TimeToPrepare int64   `json:"timeToPrepare,omitempty"`
	Healthy       int     `json:"healthy,omitempty"`
}

// create recipe output
type CreateRecipeOutput struct {
	Id int64 `json:"id"`
}

type LoginCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ScheduleResponse struct {
	Id int64 `json:"id"`
}

type CreateRecipeResponse struct {
	Id int64 `json:"id"`
}

type UpdateRecipeResponse struct {
	Response string `json:"response"`
}
