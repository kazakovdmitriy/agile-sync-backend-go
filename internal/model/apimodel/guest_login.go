package apimodel

// GuestLogin запрос гостевого входа
// @Description Данные для создания гостевого аккаунта
type GuestLogin struct {
	Name string `json:"name" example:"Guest User"`
}
