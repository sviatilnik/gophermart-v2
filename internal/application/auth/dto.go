package auth

// RegisterRequest - DTO для регистрации нового пользователя
type RegisterRequest struct {
	Login    string `json:"login" validate:"required,login"`
	Password string `json:"password" validate:"required,min=8"`
}

// LoginRequest - DTO для входа в систему
type LoginRequest struct {
	Login    string `json:"login" validate:"required,login"`
	Password string `json:"password" validate:"required"`
}

// RefreshTokenRequest - DTO для обновления токена
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// RegisterResponse - DTO ответа после регистрации
type RegisterResponse struct {
	ID    string `json:"id"`
	Login string `json:"login"`
}

// TokenResponse - DTO ответа после входа
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresIn    int64  `json:"expires_in"` // Время жизни токена (сек)
}

// ErrorResponse - Общий DTO для ошибок
type ErrorResponse struct {
	Error string `json:"error"`
}
