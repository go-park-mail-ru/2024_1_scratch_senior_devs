package models

type UserFormData struct {
	Username string `json:"username"`
	Password string `json:"-"`
}
