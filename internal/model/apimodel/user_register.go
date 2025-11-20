package apimodel

// UserRegister запрос регистрации
// @Description Данные для регистрации нового пользователя
type UserRegister struct {
	Name     string `json:"name" example:"John Doe"`
	Email    string `json:"email" example:"john@example.com"`
	Password string `json:"password" example:"securepassword123"`
}
