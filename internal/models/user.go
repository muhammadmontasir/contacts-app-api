package models

type User struct {
	ID       uint64 `json:"id" gorm:"primaryKey;autoIncrement"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"-"`
	Active   bool   `json:"active"`
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
