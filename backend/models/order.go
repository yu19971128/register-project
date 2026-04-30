package models

import "time"

type Order struct {
	ID           int64     `json:"id"`
	OrderNo      string    `json:"order_no"`
	ScheduleID   int64     `json:"schedule_id"`
	PatientID    int64     `json:"patient_id"`
	VisitorPhone string    `json:"visitor_phone"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
