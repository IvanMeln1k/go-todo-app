package domain

type User struct {
	Id       int    `json:"-"`
	Name     string `json:"name" validate:"required"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required" db:"password_hash"`
}
