package apimodel

import "time"

// CreateSessionResponse ответ при создании сессии
// @Description Данные созданной сессии
type CreateSessionResponse struct {
	SessionID string    `json:"session_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	UserID    string    `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440001"`
	CreatedAt time.Time `json:"created_at" example:"2024-01-15T10:30:00Z"`
}

// DeleteSessionResponse ответ при удалении сессии
// @Description Результат операции удаления сессии
type DeleteSessionResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Сессия успешно удалена"`
}

// ErrorResponse стандартный ответ с ошибкой
// @Description Сообщение об ошибке
type ErrorResponse struct {
	Error string `json:"error" example:"invalid session id"`
}

// DetailResponse ответ с детализированным сообщением
// @Description Детализированное сообщение об ошибке
type DetailResponse struct {
	Detail string `json:"detail" example:"Сессия не найдена"`
}
