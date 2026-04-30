package models

import "time"

type Schedule struct {
	ID          int64     `json:"id"`
	Date        string    `json:"date"`
	Department  string    `json:"department"`
	DoctorName  string    `json:"doctor_name"`
	StartTime   string    `json:"start_time"`
	EndTime     string    `json:"end_time"`
	TotalQuota  int       `json:"total_quota"`
	Remaining   int       `json:"remaining"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
