package models

import "time"

type StatusType string

const (
	New        StatusType = "NEW"
	Processing StatusType = "PROCESSING"
	Invalid    StatusType = "INVALID"
	Processed  StatusType = "PROCESSED"
	Registered StatusType = "REGISTERED"
)

type Order struct {
	ID         int        `json:"id,omitempty"`
	Number     string     `json:"number,order"`
	Status     StatusType `json:"status"`
	Accrual    int        `json:"accrual,omitempty"`
	UploadedAt time.Time  `json:"uploaded_at"`
	UserID     int        `json:"user_id,omitempty"`
}
