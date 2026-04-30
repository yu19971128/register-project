package service

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"clinic/models"
	"clinic/repo"
)

type RegistrationService struct {
	db         *sql.DB
	patientRepo *repo.PatientRepository
	scheduleRepo *repo.ScheduleRepository
	orderRepo    *repo.OrderRepository
}

func NewRegistrationService(db *sql.DB, patientRepo *repo.PatientRepository, scheduleRepo *repo.ScheduleRepository, orderRepo *repo.OrderRepository) *RegistrationService {
	return &RegistrationService{
		db:         db,
		patientRepo: patientRepo,
		scheduleRepo: scheduleRepo,
		orderRepo:    orderRepo,
	}
}

type RegistrationResult struct {
	OrderNo    string          `json:"order_no"`
	Schedule   *ScheduleInfo   `json:"schedule"`
	Patient    *PatientInfo    `json:"patient"`
	Status     string          `json:"status"`
	CreatedAt  string          `json:"created_at"`
	TicketURL  string          `json:"ticket_url"`
}

type ScheduleInfo struct {
	ID         int64  `json:"id"`
	Department string `json:"department"`
	DoctorName string `json:"doctor_name"`
	Date       string `json:"date"`
	StartTime  string `json:"start_time"`
	EndTime    string `json:"end_time"`
}

type PatientInfo struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Gender string `json:"gender"`
	Age    int    `json:"age"`
}

type TicketResult struct {
	OrderNo       string   `json:"order_no"`
	QRCodeData    string   `json:"qrcode_data"`
	Department    string   `json:"department"`
	DoctorName    string   `json:"doctor_name"`
	Date          string   `json:"date"`
	StartTime     string   `json:"start_time"`
	EndTime       string   `json:"end_time"`
	PatientName   string   `json:"patient_name"`
	PatientGender string   `json:"patient_gender"`
	PatientAge    int      `json:"patient_age"`
	Location      string   `json:"location"`
	Status        string   `json:"status"`
	Notice        []string `json:"notice"`
}

func (s *RegistrationService) SubmitRegistration(scheduleID, patientID int64, visitorPhone string) (*RegistrationResult, error) {
	// Validate schedule
	schedule, err := s.scheduleRepo.GetByID(scheduleID)
	if err != nil {
		return nil, fmt.Errorf("get schedule: %w", err)
	}
	if schedule == nil {
		return nil, fmt.Errorf("号源不存在")
	}
	if schedule.Remaining <= 0 {
		return nil, fmt.Errorf("号源余量已为 0")
	}

	// Validate patient belongs to visitor
	patient, err := s.patientRepo.GetByIDWithVisitorPhone(patientID)
	if err != nil {
		return nil, fmt.Errorf("get patient: %w", err)
	}
	if patient == nil {
		return nil, fmt.Errorf("就诊人不存在")
	}
	// 兼容历史数据：若就诊人未绑定访客手机号，自动绑定当前访客
	if patient.VisitorPhone == "" {
		patient.VisitorPhone = visitorPhone
		_ = s.patientRepo.UpdateVisitorPhone(patientID, visitorPhone)
	}
	if patient.VisitorPhone != visitorPhone {
		return nil, fmt.Errorf("就诊人不属于当前访客")
	}

	// Check duplicate
	exists, err := s.orderRepo.ExistsByScheduleAndVisitor(scheduleID, visitorPhone)
	if err != nil {
		return nil, fmt.Errorf("check duplicate: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("同一号源不可重复提交")
	}

	// Transaction: deduct schedule + create order
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	// Deduct within transaction
	res, err := tx.Exec(`UPDATE schedules SET remaining = remaining - 1 WHERE id = ? AND remaining > 0`, scheduleID)
	if err != nil {
		return nil, fmt.Errorf("deduct schedule: %w", err)
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return nil, fmt.Errorf("号源余量已为 0")
	}

	// Generate order number
	orderNo := generateOrderNo(schedule.Date)

	// Create order within transaction
	_, err = tx.Exec(
		`INSERT INTO orders (order_no, schedule_id, patient_id, visitor_phone, status) VALUES (?, ?, ?, ?, ?)`,
		orderNo, scheduleID, patientID, visitorPhone, "confirmed",
	)
	if err != nil {
		return nil, fmt.Errorf("create order: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit tx: %w", err)
	}

	// Update status to full if needed
	updatedSchedule, _ := s.scheduleRepo.GetByID(scheduleID)
	if updatedSchedule != nil && updatedSchedule.Remaining == 0 {
		_ = s.scheduleRepo.Update(&models.Schedule{ID: scheduleID, Status: "full"})
	}

	return &RegistrationResult{
		OrderNo: orderNo,
		Schedule: &ScheduleInfo{
			ID:         scheduleID,
			Department: schedule.Department,
			DoctorName: schedule.DoctorName,
			Date:       schedule.Date,
			StartTime:  schedule.StartTime,
			EndTime:    schedule.EndTime,
		},
		Patient: &PatientInfo{
			ID:     patientID,
			Name:   patient.Name,
			Gender: patient.Gender,
			Age:    patient.Age,
		},
		Status:    "confirmed",
		CreatedAt: time.Now().Format(time.RFC3339),
		TicketURL: "/h5/register/ticket?order_no=" + orderNo,
	}, nil
}

func (s *RegistrationService) GetTicket(orderNo, visitorPhone string) (*TicketResult, error) {
	order, err := s.orderRepo.GetByOrderNo(orderNo)
	if err != nil {
		return nil, fmt.Errorf("get order: %w", err)
	}
	if order == nil {
		return nil, fmt.Errorf("挂号凭证不存在")
	}
	if order.VisitorPhone != visitorPhone {
		return nil, fmt.Errorf("无权查看该凭证")
	}

	schedule, err := s.scheduleRepo.GetByID(order.ScheduleID)
	if err != nil {
		return nil, fmt.Errorf("get schedule: %w", err)
	}

	patient, err := s.patientRepo.GetByID(order.PatientID)
	if err != nil {
		return nil, fmt.Errorf("get patient: %w", err)
	}

	return &TicketResult{
		OrderNo:       order.OrderNo,
		QRCodeData:    order.OrderNo,
		Department:    schedule.Department,
		DoctorName:    schedule.DoctorName,
		Date:          schedule.Date,
		StartTime:     schedule.StartTime,
		EndTime:       schedule.EndTime,
		PatientName:   patient.Name,
		PatientGender: patient.Gender,
		PatientAge:    patient.Age,
		Location:      "1号楼 2层 " + schedule.Department + "诊室",
		Status:        order.Status,
		Notice: []string{
			"请提前 15 分钟到院取号",
			"凭此二维码或订单号就诊",
			"如需退号请提前 30 分钟",
		},
	}, nil
}

func generateOrderNo(date string) string {
	dateStr := strings.ReplaceAll(date, "-", "")
	return fmt.Sprintf("GH%s%012d", dateStr, time.Now().UnixNano()%1000000000000)
}
