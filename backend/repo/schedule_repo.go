package repo

import (
	"database/sql"
	"fmt"
	"strings"

	"clinic/models"
)

type ScheduleRepository struct {
	db *sql.DB
}

func NewScheduleRepository(db *sql.DB) *ScheduleRepository {
	return &ScheduleRepository{db: db}
}

func (r *ScheduleRepository) Create(s *models.Schedule) (int64, error) {
	res, err := r.db.Exec(
		`INSERT INTO schedules (date, department, doctor_name, start_time, end_time, total_quota, remaining, status)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		s.Date, s.Department, s.DoctorName, s.StartTime, s.EndTime, s.TotalQuota, s.Remaining, s.Status,
	)
	if err != nil {
		return 0, fmt.Errorf("create schedule: %w", err)
	}
	return res.LastInsertId()
}

func (r *ScheduleRepository) GetByID(id int64) (*models.Schedule, error) {
	row := r.db.QueryRow(`SELECT id, date, department, doctor_name, start_time, end_time, total_quota, remaining, status, created_at, updated_at FROM schedules WHERE id = ?`, id)
	s := &models.Schedule{}
	err := row.Scan(&s.ID, &s.Date, &s.Department, &s.DoctorName, &s.StartTime, &s.EndTime, &s.TotalQuota, &s.Remaining, &s.Status, &s.CreatedAt, &s.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get schedule: %w", err)
	}
	return s, nil
}

func (r *ScheduleRepository) List(date, department, doctorName string, offset, limit int) ([]*models.Schedule, int, error) {
	where := []string{"1 = 1"}
	args := []interface{}{}
	if date != "" {
		where = append(where, "date = ?")
		args = append(args, date)
	}
	if department != "" {
		where = append(where, "department = ?")
		args = append(args, department)
	}
	if doctorName != "" {
		where = append(where, "doctor_name = ?")
		args = append(args, doctorName)
	}
	whereStr := strings.Join(where, " AND ")

	var total int
	countArgs := append([]interface{}{}, args...)
	if err := r.db.QueryRow("SELECT COUNT(*) FROM schedules WHERE "+whereStr, countArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count schedules: %w", err)
	}

	var rows *sql.Rows
	var err error
	if limit <= 0 {
		rows, err = r.db.Query(
			"SELECT id, date, department, doctor_name, start_time, end_time, total_quota, remaining, status, created_at FROM schedules WHERE "+whereStr+" ORDER BY start_time ASC",
			args...,
		)
	} else {
		queryArgs := append([]interface{}{}, args...)
		queryArgs = append(queryArgs, limit, offset)
		rows, err = r.db.Query(
			"SELECT id, date, department, doctor_name, start_time, end_time, total_quota, remaining, status, created_at FROM schedules WHERE "+whereStr+" ORDER BY start_time ASC LIMIT ? OFFSET ?",
			queryArgs...,
		)
	}
	if err != nil {
		return nil, 0, fmt.Errorf("list schedules: %w", err)
	}
	defer rows.Close()
	return scanScheduleList(rows), total, nil
}

func scanScheduleList(rows *sql.Rows) []*models.Schedule {
	var list []*models.Schedule
	for rows.Next() {
		s := &models.Schedule{}
		_ = rows.Scan(&s.ID, &s.Date, &s.Department, &s.DoctorName, &s.StartTime, &s.EndTime, &s.TotalQuota, &s.Remaining, &s.Status, &s.CreatedAt)
		list = append(list, s)
	}
	return list
}

func (r *ScheduleRepository) Update(s *models.Schedule) error {
	fields := []string{}
	args := []interface{}{}
	if s.Date != "" {
		fields = append(fields, "date = ?")
		args = append(args, s.Date)
	}
	if s.Department != "" {
		fields = append(fields, "department = ?")
		args = append(args, s.Department)
	}
	if s.DoctorName != "" {
		fields = append(fields, "doctor_name = ?")
		args = append(args, s.DoctorName)
	}
	if s.StartTime != "" {
		fields = append(fields, "start_time = ?")
		args = append(args, s.StartTime)
	}
	if s.EndTime != "" {
		fields = append(fields, "end_time = ?")
		args = append(args, s.EndTime)
	}
	if s.TotalQuota > 0 {
		fields = append(fields, "total_quota = ?")
		args = append(args, s.TotalQuota)
	}
	if s.Remaining >= 0 {
		fields = append(fields, "remaining = ?")
		args = append(args, s.Remaining)
	}
	if s.Status != "" {
		fields = append(fields, "status = ?")
		args = append(args, s.Status)
	}
	if len(fields) == 0 {
		return nil
	}
	args = append(args, s.ID)
	_, err := r.db.Exec("UPDATE schedules SET "+strings.Join(fields, ", ")+" WHERE id = ?", args...)
	if err != nil {
		return fmt.Errorf("update schedule: %w", err)
	}
	return nil
}

func (r *ScheduleRepository) Delete(id int64) error {
	_, err := r.db.Exec(`DELETE FROM schedules WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete schedule: %w", err)
	}
	return nil
}

// Deduct atomically decreases remaining by 1 using optimistic lock.
func (r *ScheduleRepository) Deduct(id int64) (bool, error) {
	res, err := r.db.Exec(`UPDATE schedules SET remaining = remaining - 1 WHERE id = ? AND remaining > 0`, id)
	if err != nil {
		return false, fmt.Errorf("deduct schedule: %w", err)
	}
	affected, _ := res.RowsAffected()
	return affected > 0, nil
}

// Rollback atomically increases remaining by 1 (capped by total_quota).
func (r *ScheduleRepository) Rollback(id int64) (bool, error) {
	res, err := r.db.Exec(`UPDATE schedules SET remaining = remaining + 1 WHERE id = ? AND remaining < total_quota`, id)
	if err != nil {
		return false, fmt.Errorf("rollback schedule: %w", err)
	}
	affected, _ := res.RowsAffected()
	return affected > 0, nil
}
