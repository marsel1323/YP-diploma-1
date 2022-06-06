package models

type User struct {
	ID       int    `json:"id,omitempty"`
	Login    string `json:"login"`
	Password string `json:"password"`
}
