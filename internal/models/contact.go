package models

type Contact struct {
	ID     uint64 `json:"id" gorm:"primaryKey;autoIncrement"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Phone  string `json:"phone"`
	UserID uint64 `json:"user_id"`
}
