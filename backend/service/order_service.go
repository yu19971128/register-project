package service

import (
	"database/sql"
	"fmt"
	"time"

	"clinic/models"
	"clinic/repo"
)

type OrderService struct {
	db           *sql.DB
	orderRepo    *repo.OrderRepository
	scheduleRepo *repo.ScheduleRepository
	patientRepo  *repo.PatientRepository
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
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize

	vp := ""
	if !isAdmin {
		vp = visitorPhone
	}

	orders, total, err := s.orderRepo.List(date, department, doctorName, status, keyword, offset, pageSize, vp)
	if err != nil {
		return nil, err
	}

	var list []OrderListItem
	for _, o := range orders {
		item := OrderListItem{
			ID:        o.ID,
			OrderNo:   o.OrderNo,
			Status:    o.Status,
			CreatedAt: o.CreatedAt.Format(time.RFC3339),
		}

		patient, _ := s.patientRepo.GetByID(o.PatientID)
		if patient != nil {
			item.PatientName = patient.Name
		}

		schedule, _ := s.scheduleRepo.GetByID(o.ScheduleID)
		if schedule != nil {
			item.Department = schedule.Department
			item.DoctorName = schedule.DoctorName
			item.Date = schedule.Date
			item.StartTime = schedule.StartTime
			item.EndTime = schedule.EndTime
		}

		list = append(list, item)
	}

	return &OrderListResult{
		Total: total,
		List:  list,
	}, nil
}

func (s *OrderService) GetDetail(id int64, isAdmin bool, visitorPhone string) (*OrderDetailResult, error) {
	order, err := s.orderRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("get order: %w", err)
	}
	if order == nil {
		return nil, fmt.Errorf("订单不存在")
	}

	if !isAdmin && order.VisitorPhone != visitorPhone {
		return nil, fmt.Errorf("无权查看该订单")
	}

	schedule, _ := s.scheduleRepo.GetByID(order.ScheduleID)
	patient, _ := s.patientRepo.GetByID(order.PatientID)

	res := &OrderDetailResult{
		ID:           order.ID,
		OrderNo:      order.OrderNo,
		Status:       order.Status,
		VisitorPhone: order.VisitorPhone,
		CreatedAt:    order.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    order.UpdatedAt.Format(time.RFC3339),
	}

	if schedule != nil {
		res.Schedule = &ScheduleInfo{
			ID:         schedule.ID,
			Department: schedule.Department,
			DoctorName: schedule.DoctorName,
			Date:       schedule.Date,
			StartTime:  schedule.StartTime,
			EndTime:    schedule.EndTime,
		}
	}

	if patient != nil {
		res.Patient = &PatientInfo{
			ID:     patient.ID,
			Name:   patient.Name,
			Gender: patient.Gender,
			Age:    patient.Age,
		}
	}

	return res, nil
}

func (s *OrderService) Cancel(id int64, reason string, isAdmin bool, visitorPhone, operatedBy string) (*OrderResult, error) {
	order, err := s.orderRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("get order: %w", err)
	}
	if order == nil {
		return nil, fmt.Errorf("订单不存在")
	}

	if !isAdmin && order.VisitorPhone != visitorPhone {
		return nil, fmt.Errorf("无权退号")
	}

	schedule, err := s.scheduleRepo.GetByID(order.ScheduleID)
	if err != nil {
		return nil, fmt.Errorf("get schedule: %w", err)
	}

	if !canCancelOrder(order, schedule) {
		return nil, fmt.Errorf("已超过退号时限")
	}

	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	now := time.Now()
	_, err = tx.Exec(
		"UPDATE orders SET status = ?, cancel_reason = ?, cancelled_at = ?, operated_by = ? WHERE id = ?",
		"cancelled", reason, now, operatedBy, id,
	)
	if err != nil {
		return nil, fmt.Errorf("update order: %w", err)
	}

	_, err = tx.Exec(
		"UPDATE schedules SET remaining = remaining + 1 WHERE id = ? AND remaining < total_quota",
		order.ScheduleID,
	)
	if err != nil {
		return nil, fmt.Errorf("rollback schedule: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit tx: %w", err)
	}

	return &OrderResult{
		ID:          order.ID,
		OrderNo:     order.OrderNo,
		Status:      "cancelled",
		CancelledAt: now.Format(time.RFC3339),
	}, nil
}

func (s *OrderService) Change(id, newScheduleID int64, isAdmin bool, visitorPhone string) (*OrderDetailResult, error) {
	oldOrder, err := s.orderRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("get order: %w", err)
	}
	if oldOrder == nil {
		return nil, fmt.Errorf("订单不存在")
	}

	if !isAdmin && oldOrder.VisitorPhone != visitorPhone {
		return nil, fmt.Errorf("无权改号")
	}

	oldSchedule, err := s.scheduleRepo.GetByID(oldOrder.ScheduleID)
	if err != nil {
		return nil, fmt.Errorf("get schedule: %w", err)
	}

	if !canCancelOrder(oldOrder, oldSchedule) {
		return nil, fmt.Errorf("原订单已超过退号时限")
	}

	newSchedule, err := s.scheduleRepo.GetByID(newScheduleID)
	if err != nil {
		return nil, fmt.Errorf("get new schedule: %w", err)
	}
	if newSchedule == nil {
		return nil, fmt.Errorf("新号源不存在")
	}
	if newSchedule.Remaining <= 0 {
		return nil, fmt.Errorf("新号源余量已为 0")
	}

	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	res, err := tx.Exec(
		"UPDATE schedules SET remaining = remaining - 1 WHERE id = ? AND remaining > 0",
		newScheduleID,
	)
	if err != nil {
		return nil, fmt.Errorf("deduct schedule: %w", err)
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return nil, fmt.Errorf("新号源余量已为 0")
	}

	newOrderNo := generateOrderNo(newSchedule.Date)
	orderRes, err := tx.Exec(
		"INSERT INTO orders (order_no, schedule_id, patient_id, visitor_phone, status) VALUES (?, ?, ?, ?, ?)",
		newOrderNo, newScheduleID, oldOrder.PatientID, oldOrder.VisitorPhone, "confirmed",
	)
	if err != nil {
		return nil, fmt.Errorf("create order: %w", err)
	}
	newOrderID, _ := orderRes.LastInsertId()

	now := time.Now()
	_, err = tx.Exec(
		"UPDATE orders SET status = ?, cancel_reason = ?, cancelled_at = ?, operated_by = ? WHERE id = ?",
		"cancelled", "改号", now, "", id,
	)
	if err != nil {
		return nil, fmt.Errorf("update old order: %w", err)
	}

	_, err = tx.Exec(
		"UPDATE schedules SET remaining = remaining + 1 WHERE id = ? AND remaining < total_quota",
		oldOrder.ScheduleID,
	)
	if err != nil {
		return nil, fmt.Errorf("rollback old schedule: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit tx: %w", err)
	}

	patient, _ := s.patientRepo.GetByID(oldOrder.PatientID)

	result := &OrderDetailResult{
		ID:           newOrderID,
		OrderNo:      newOrderNo,
		Status:       "confirmed",
		VisitorPhone: oldOrder.VisitorPhone,
		CreatedAt:    now.Format(time.RFC3339),
		UpdatedAt:    now.Format(time.RFC3339),
		Schedule: &ScheduleInfo{
			ID:         newSchedule.ID,
			Department: newSchedule.Department,
			DoctorName: newSchedule.DoctorName,
			Date:       newSchedule.Date,
			StartTime:  newSchedule.StartTime,
			EndTime:    newSchedule.EndTime,
		},
	}
	if patient != nil {
		result.Patient = &PatientInfo{
			ID:     patient.ID,
			Name:   patient.Name,
			Gender: patient.Gender,
			Age:    patient.Age,
		}
	}

	return result, nil
}

func (s *OrderService) Complete(id int64, isAdmin bool, visitorPhone, operatedBy string) (*OrderResult, error) {
	order, err := s.orderRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("get order: %w", err)
	}
	if order == nil {
		return nil, fmt.Errorf("订单不存在")
	}

	if !isAdmin && order.VisitorPhone != visitorPhone {
		return nil, fmt.Errorf("无权完成订单")
	}

	if order.Status != "confirmed" {
		return nil, fmt.Errorf("只有待就诊订单可以完成")
	}

	now := time.Now()
	if err := s.orderRepo.Complete(id, now, operatedBy); err != nil {
		return nil, err
	}

	return &OrderResult{
		ID:      order.ID,
		OrderNo: order.OrderNo,
		Status:  "completed",
	}, nil
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
