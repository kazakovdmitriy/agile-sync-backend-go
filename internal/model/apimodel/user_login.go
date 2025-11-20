package apimodel

// UserLogin запрос авторизации
// @Description Данные для входа пользователя
type UserLogin struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"securepassword123"`
}
