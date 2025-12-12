package models

// User — структура нашего пользователя
// Теги `json:"..."` нужны, чтобы Go знал, как превращать это в JSON для API
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}