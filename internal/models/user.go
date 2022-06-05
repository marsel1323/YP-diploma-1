package models

type User struct {
	Id       int64  `json:"id,omitempty"`
	Login    string `json:"login"`
	Password string `json:"password"`
}
