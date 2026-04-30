package repo

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

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
	row := r.db.QueryRow(`SELECT id, order_no, schedule_id, patient_id, visitor_phone, status, cancel_reason, cancelled_at, completed_at, operated_by, created_at, updated_at FROM orders WHERE order_no = ?`, orderNo)
	o := &models.Order{}
	var cancelledAt, completedAt sql.NullTime
	var cancelReason, operatedBy sql.NullString
	err := row.Scan(&o.ID, &o.OrderNo, &o.ScheduleID, &o.PatientID, &o.VisitorPhone, &o.Status, &cancelReason, &cancelledAt, &completedAt, &operatedBy, &o.CreatedAt, &o.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get order by no: %w", err)
	}
	if cancelReason.Valid {
		o.CancelReason = cancelReason.String
	}
	if cancelledAt.Valid {
		o.CancelledAt = &cancelledAt.Time
	}
	if completedAt.Valid {
		o.CompletedAt = &completedAt.Time
	}
	if operatedBy.Valid {
		o.OperatedBy = operatedBy.String
	}
	return o, nil
}

func (r *OrderRepository) GetByID(id int64) (*models.Order, error) {
	row := r.db.QueryRow(`SELECT id, order_no, schedule_id, patient_id, visitor_phone, status, cancel_reason, cancelled_at, completed_at, operated_by, created_at, updated_at FROM orders WHERE id = ?`, id)
	o := &models.Order{}
	var cancelledAt, completedAt sql.NullTime
	var cancelReason, operatedBy sql.NullString
	err := row.Scan(&o.ID, &o.OrderNo, &o.ScheduleID, &o.PatientID, &o.VisitorPhone, &o.Status, &cancelReason, &cancelledAt, &completedAt, &operatedBy, &o.CreatedAt, &o.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get order: %w", err)
	}
	if cancelReason.Valid {
		o.CancelReason = cancelReason.String
	}
	if cancelledAt.Valid {
		o.CancelledAt = &cancelledAt.Time
	}
	if completedAt.Valid {
		o.CompletedAt = &completedAt.Time
	}
	if operatedBy.Valid {
		o.OperatedBy = operatedBy.String
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

func (r *OrderRepository) List(date, department, doctorName, status, keyword string, offset, limit int, visitorPhone ...string) ([]*models.Order, int, error) {
	joins := "LEFT JOIN schedules s ON o.schedule_id = s.id LEFT JOIN patients p ON o.patient_id = p.id"
	where := []string{"1 = 1"}
	args := []interface{}{}
	if date != "" {
		where = append(where, "o.created_at >= ? AND o.created_at < ?")
		nextDate, _ := time.Parse("2006-01-02", date)
		nextDate = nextDate.Add(24 * time.Hour)
		args = append(args, date, nextDate.Format("2006-01-02"))
	}
	if department != "" {
		where = append(where, "s.department = ?")
		args = append(args, department)
	}
	if doctorName != "" {
		where = append(where, "s.doctor_name = ?")
		args = append(args, doctorName)
	}
	if status != "" {
		where = append(where, "o.status = ?")
		args = append(args, status)
	}
	if keyword != "" {
		where = append(where, "(o.order_no LIKE ? OR p.name LIKE ?)")
		k := "%" + keyword + "%"
		args = append(args, k, k)
	}
	if len(visitorPhone) > 0 && visitorPhone[0] != "" {
		where = append(where, "o.visitor_phone = ?")
		args = append(args, visitorPhone[0])
	}
	whereStr := strings.Join(where, " AND ")

	var total int
	countArgs := append([]interface{}{}, args...)
	if err := r.db.QueryRow("SELECT COUNT(*) FROM orders o "+joins+" WHERE "+whereStr, countArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count orders: %w", err)
	}

	queryArgs := append([]interface{}{}, args...)
	queryArgs = append(queryArgs, limit, offset)
	rows, err := r.db.Query(
		"SELECT o.id, o.order_no, o.schedule_id, o.patient_id, o.visitor_phone, o.status, o.cancel_reason, o.cancelled_at, o.completed_at, o.operated_by, o.created_at FROM orders o "+joins+" WHERE "+whereStr+" ORDER BY o.created_at DESC LIMIT ? OFFSET ?",
		queryArgs...,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("list orders: %w", err)
	}
	defer rows.Close()
	return scanOrderList(rows), total, nil
}

func scanOrderList(rows *sql.Rows) []*models.Order {
	var list []*models.Order
	for rows.Next() {
		o := &models.Order{}
		var cancelledAt, completedAt sql.NullTime
		_ = rows.Scan(&o.ID, &o.OrderNo, &o.ScheduleID, &o.PatientID, &o.VisitorPhone, &o.Status, &o.CancelReason, &cancelledAt, &completedAt, &o.OperatedBy, &o.CreatedAt)
		if cancelledAt.Valid {
			o.CancelledAt = &cancelledAt.Time
		}
		if completedAt.Valid {
			o.CompletedAt = &completedAt.Time
		}
		list = append(list, o)
	}
	return list
}

func (r *OrderRepository) UpdateStatus(id int64, status, reason string, cancelledAt time.Time, operatedBy string) error {
	_, err := r.db.Exec(
		"UPDATE orders SET status = ?, cancel_reason = ?, cancelled_at = ?, operated_by = ? WHERE id = ?",
		status, reason, cancelledAt, operatedBy, id,
	)
	if err != nil {
		return fmt.Errorf("update order status: %w", err)
	}
	return nil
}
