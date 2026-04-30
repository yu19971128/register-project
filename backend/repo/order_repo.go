package repo

import (
	"database/sql"
	"fmt"

	"clinic/models"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(o *models.Order) (int64, error) {
	res, err := r.db.Exec(
		`INSERT INTO orders (order_no, schedule_id, patient_id, visitor_phone, status)
		VALUES (?, ?, ?, ?, ?)`,
		o.OrderNo, o.ScheduleID, o.PatientID, o.VisitorPhone, o.Status,
	)
	if err != nil {
		return 0, fmt.Errorf("create order: %w", err)
	}
	return res.LastInsertId()
}

func (r *OrderRepository) GetByOrderNo(orderNo string) (*models.Order, error) {
	row := r.db.QueryRow(`SELECT id, order_no, schedule_id, patient_id, visitor_phone, status, created_at, updated_at FROM orders WHERE order_no = ?`, orderNo)
	o := &models.Order{}
	err := row.Scan(&o.ID, &o.OrderNo, &o.ScheduleID, &o.PatientID, &o.VisitorPhone, &o.Status, &o.CreatedAt, &o.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get order by no: %w", err)
	}
	return o, nil
}

func (r *OrderRepository) GetByID(id int64) (*models.Order, error) {
	row := r.db.QueryRow(`SELECT id, order_no, schedule_id, patient_id, visitor_phone, status, created_at, updated_at FROM orders WHERE id = ?`, id)
	o := &models.Order{}
	err := row.Scan(&o.ID, &o.OrderNo, &o.ScheduleID, &o.PatientID, &o.VisitorPhone, &o.Status, &o.CreatedAt, &o.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get order: %w", err)
	}
	return o, nil
}

func (r *OrderRepository) ExistsByScheduleAndVisitor(scheduleID int64, visitorPhone string) (bool, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM orders WHERE schedule_id = ? AND visitor_phone = ? AND status = 'confirmed'`, scheduleID, visitorPhone).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("check duplicate order: %w", err)
	}
	return count > 0, nil
}
