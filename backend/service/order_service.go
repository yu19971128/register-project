package service

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	"clinic/models"
	"clinic/repo"
)

type OrderService struct {
	db           *sql.DB
	orderRepo    *repo.OrderRepository
	scheduleRepo *repo.ScheduleRepository
	patientRepo  *repo.PatientRepository
	mu           sync.Mutex
	seq          int
}

func NewOrderService(db *sql.DB, orderRepo *repo.OrderRepository, scheduleRepo *repo.ScheduleRepository, patientRepo *repo.PatientRepository) *OrderService {
	return &OrderService{
		db:           db,
		orderRepo:    orderRepo,
		scheduleRepo: scheduleRepo,
		patientRepo:  patientRepo,
	}
}

type OrderListResult struct {
	Total int             `json:"total"`
	List  []OrderListItem `json:"list"`
}

type OrderListItem struct {
	ID          int64  `json:"id"`
	OrderNo     string `json:"order_no"`
	PatientName string `json:"patient_name"`
	Department  string `json:"department"`
	DoctorName  string `json:"doctor_name"`
	Date        string `json:"date"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
}

type OrderDetailResult struct {
	ID           int64         `json:"id"`
	OrderNo      string        `json:"order_no"`
	Status       string        `json:"status"`
	Schedule     *ScheduleInfo `json:"schedule"`
	Patient      *PatientInfo  `json:"patient"`
	VisitorPhone string        `json:"visitor_phone"`
	CreatedAt    string        `json:"created_at"`
	UpdatedAt    string        `json:"updated_at"`
}

type OrderResult struct {
	ID          int64  `json:"id"`
	OrderNo     string `json:"order_no"`
	Status      string `json:"status"`
	CancelledAt string `json:"cancelled_at,omitempty"`
}

func (s *OrderService) List(date, department, doctorName, status, keyword string, page, pageSize int, isAdmin bool, visitorPhone string) (*OrderListResult, error) {
	return nil, nil
}

func (s *OrderService) GetDetail(id int64, isAdmin bool, visitorPhone string) (*OrderDetailResult, error) {
	return nil, nil
}

func (s *OrderService) Cancel(id int64, reason string, isAdmin bool, visitorPhone, operatedBy string) (*OrderResult, error) {
	return nil, nil
}

func (s *OrderService) Change(id, newScheduleID int64, isAdmin bool, visitorPhone string) (*OrderDetailResult, error) {
	return nil, nil
}

func (s *OrderService) generateOrderNo(date string) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.seq++
	dateStr := strings.ReplaceAll(date, "-", "")
	return fmt.Sprintf("GH%s%04d", dateStr, s.seq)
}

func canCancelOrder(order *models.Order, schedule *models.Schedule) bool {
	if order.Status != "confirmed" {
		return false
	}
	appointmentTime, err := time.Parse("2006-01-02 15:04", schedule.Date+" "+schedule.StartTime)
	if err != nil {
		return false
	}
	return time.Now().Add(30 * time.Minute).Before(appointmentTime)
}
