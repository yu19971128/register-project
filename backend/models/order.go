package models

import "time"

type Order struct {
	ID           int64     `json:"id"`
	OrderNo      string    `json:"order_no"`
	ScheduleID   int64     `json:"schedule_id"`
	PatientID    int64     `json:"patient_id"`
	VisitorPhone string    `json:"visitor_phone"`
	Status       string    `json:"status"`
	CancelReason string    `json:"cancel_reason,omitempty"`
	CancelledAt  *time.Time `json:"cancelled_at,omitempty"`
	CompletedAt  *time.Time `json:"completed_at,omitempty"`
	OperatedBy   string    `json:"operated_by,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
