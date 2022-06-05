package models

import "time"

type StatusType string

const (
	New        StatusType = "NEW"
	Processing StatusType = "PROCESSING"
	Invalid    StatusType = "INVALID"
	Processed  StatusType = "PROCESSED"
)

type Order struct {
	Number     string     `json:"number"`
	Status     StatusType `json:"status"`
	Accrual    int        `json:"accrual,omitempty"`
	UploadedAt time.Time  `json:"uploaded_at"`
}
