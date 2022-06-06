package models

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
	Number     string     `json:"number"`
	Status     StatusType `json:"status"`
	Accrual    float64    `json:"accrual,omitempty"`
	UploadedAt string     `json:"uploaded_at"`
	UserID     int        `json:"user_id,omitempty"`
}
