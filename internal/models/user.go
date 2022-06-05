package models

type User struct {
	ID       int64  `json:"id,omitempty"`
	Login    string `json:"login"`
	Password string `json:"password"`
}
