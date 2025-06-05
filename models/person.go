package models

// Person представляет модель человека в базе данных
type Person struct {
	ID          uint   `gorm:"primaryKey" json:"id"`    // Уникальный идентификатор
	Name        string `gorm:"not null" json:"name"`    // Имя (обязательное)
	Surname     string `gorm:"not null" json:"surname"` // Фамилия (обязательная)
	Patronymic  string `json:"patronymic,omitempty"`    // Отчество (опционально)
	Age         *int   `json:"age,omitempty"`           // Возраст (опционально)
	Gender      string `json:"gender,omitempty"`        // Пол (опционально)
	Nationality string `json:"nationality,omitempty"`   // Национальность (опционально)
}

// ErrorResponse представляет структуру ошибки для ответа API
type ErrorResponse struct {
	Error string `json:"error" example:"внутренняя ошибка сервера"`
}

// MessageResponse представляет успешное сообщение от API
type MessageResponse struct {
	Message string `json:"message" example:"Удалён"`
}
