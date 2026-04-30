package models

import "time"

type Patient struct {
	ID               int64     `json:"id"`
	Name             string    `json:"name"`
	IDCard           string    `json:"id_card"`
	IDCardEncrypted  string    `json:"-"`
	Phone            string    `json:"phone"`
	PhoneEncrypted   string    `json:"-"`
	Gender           string    `json:"gender,omitempty"`
	Age              int       `json:"age,omitempty"`
	Address          string    `json:"address,omitempty"`
	VisitorPhone     string    `json:"-"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
