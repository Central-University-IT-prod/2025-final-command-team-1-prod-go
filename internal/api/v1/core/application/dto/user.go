package dto

type BindTokenDto struct {
	Token string `json:"token" binding:"required"`
}

type CreateUserDto struct {
	Username         string `json:"username" db:"username" binding:"required,max=32,min=6"`
	Email            string `json:"email" db:"email" binding:"required,email,max=64,min=6"`
	Password         string `json:"password" db:"password" binding:"required,max=64,min=6"`
	TelegramUsername string `json:"telegram_username" db:"telegram_username" binding:"omitempty,max=32,min=4"`
}

type UserDto struct {
	Username         string `json:"username" db:"username"`
	Email            string `json:"email" db:"email"`
	Password         string `json:"password" db:"password"`
	CreatedAt        string `json:"created_at" db:"created_at"`
	UpdatedAt        string `json:"updated_at" db:"updated_at"`
	TelegramUsername string `json:"telegram_username" db:"telegram_username" binding:"omitempty,max=32,min=4"`
	IsAdmin          bool   `json:"is_admin" db:"is_admin" binding:"omitempty"`
}

type UserWithRatingDto struct {
	Username         string  `json:"username" db:"username"`
	Email            string  `json:"email" db:"email"`
	CreatedAt        string  `json:"created_at" db:"created_at"`
	UpdatedAt        string  `json:"updated_at" db:"updated_at"`
	TelegramUsername string  `json:"telegram_username" db:"telegram_username" binding:"omitempty,max=32,min=4"`
	IsAdmin          bool    `json:"is_admin" db:"is_admin" binding:"omitempty"`
	Rating           float64 `json:"rating"`
}

type GetUserDto struct {
	Username         string `json:"username" db:"username"`
	Email            string `json:"email" db:"email"`
	CreatedAt        string `json:"created_at" db:"created_at"`
	UpdatedAt        string `json:"updated_at" db:"updated_at"`
	TelegramUsername string `json:"telegram_username" db:"telegram_username" binding:"omitempty,max=32,min=4"`
	IsAdmin          bool   `json:"is_admin" db:"is_admin" binding:"omitempty"`
}

type UpdateUserDto struct {
	Username string `json:"username,omitempty" db:"username" binding:"omitempty,max=32,min=6"`
	Password string `json:"password,omitempty" db:"password"`
	// sets in service automatically
	UpdatedAt        string `json:"updated_at,omitempty" db:"updated_at"`
	TelegramUsername string `json:"telegram_username" db:"telegram_username" binding:"omitempty,max=32,min=4"`
}
