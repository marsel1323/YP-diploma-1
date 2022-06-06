package models

type Balance struct {
	ID        int     `json:"id,omitempty"`
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
	UserID    int     `json:"user_id,omitempty"`
}
