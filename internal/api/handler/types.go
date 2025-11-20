package handler

// ErrorResponse стандартный ответ с ошибкой
// @Description Сообщение об ошибке
type ErrorResponse struct {
	Error string `json:"error" example:"invalid credentials"`
}

// DetailResponse ответ с детализированным сообщением
// @Description Детализированное сообщение об ошибке
type DetailResponse struct {
	Detail string `json:"detail" example:"Пользователь с таким email уже существует"`
}
