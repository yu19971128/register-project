package repo

import (
	"database/sql"
	"fmt"
	"strings"

	"clinic/models"
)

type PatientRepository struct {
	db *sql.DB
}

func NewPatientRepository(db *sql.DB) *PatientRepository {
	return &PatientRepository{db: db}
}

func (r *PatientRepository) Create(p *models.Patient) (int64, error) {
	res, err := r.db.Exec(
		`INSERT INTO patients (name, id_card, id_card_encrypted, phone, phone_encrypted, gender, age, address, visitor_phone)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		p.Name, p.IDCard, p.IDCardEncrypted, p.Phone, p.PhoneEncrypted, p.Gender, p.Age, p.Address, p.VisitorPhone,
	)
	if err != nil {
		return 0, fmt.Errorf("create patient: %w", err)
	}
	return res.LastInsertId()
}

func (r *PatientRepository) GetByID(id int64) (*models.Patient, error) {
	row := r.db.QueryRow(`SELECT id, name, id_card, phone, gender, age, address, created_at, updated_at FROM patients WHERE id = ?`, id)
	p := &models.Patient{}
	err := row.Scan(&p.ID, &p.Name, &p.IDCard, &p.Phone, &p.Gender, &p.Age, &p.Address, &p.CreatedAt, &p.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get patient: %w", err)
	}
	return p, nil
}

func (r *PatientRepository) ListByVisitorPhone(vp string, offset, limit int) ([]*models.Patient, int, error) {
	var total int
	if err := r.db.QueryRow(`SELECT COUNT(*) FROM patients WHERE visitor_phone = ?`, vp).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count patients: %w", err)
	}
	rows, err := r.db.Query(
		`SELECT id, name, phone, gender, age, created_at FROM patients WHERE visitor_phone = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`,
		vp, limit, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("list patients: %w", err)
	}
	defer rows.Close()
	return scanPatientList(rows), total, nil
}

func (r *PatientRepository) List(keyword string, offset, limit int) ([]*models.Patient, int, error) {
	var total int
	var err error
	if keyword != "" {
		k := "%" + keyword + "%"
		err = r.db.QueryRow(`SELECT COUNT(*) FROM patients WHERE name LIKE ? OR phone LIKE ? OR id_card LIKE ?`, k, k, k).Scan(&total)
	} else {
		err = r.db.QueryRow(`SELECT COUNT(*) FROM patients`).Scan(&total)
	}
	if err != nil {
		return nil, 0, fmt.Errorf("count patients: %w", err)
	}

	var rows *sql.Rows
	if keyword != "" {
		k := "%" + keyword + "%"
		rows, err = r.db.Query(
			`SELECT id, name, phone, gender, age, created_at FROM patients WHERE name LIKE ? OR phone LIKE ? OR id_card LIKE ? ORDER BY created_at DESC LIMIT ? OFFSET ?`,
			k, k, k, limit, offset,
		)
	} else {
		rows, err = r.db.Query(
			`SELECT id, name, phone, gender, age, created_at FROM patients ORDER BY created_at DESC LIMIT ? OFFSET ?`,
			limit, offset,
		)
	}
	if err != nil {
		return nil, 0, fmt.Errorf("list patients: %w", err)
	}
	defer rows.Close()
	return scanPatientList(rows), total, nil
}

func scanPatientList(rows *sql.Rows) []*models.Patient {
	var list []*models.Patient
	for rows.Next() {
		p := &models.Patient{}
		_ = rows.Scan(&p.ID, &p.Name, &p.Phone, &p.Gender, &p.Age, &p.CreatedAt)
		list = append(list, p)
	}
	return list
}

func (r *PatientRepository) Update(p *models.Patient) error {
	fields := []string{}
	args := []interface{}{}
	if p.Name != "" {
		fields = append(fields, "name = ?")
		args = append(args, p.Name)
	}
	if p.IDCard != "" {
		fields = append(fields, "id_card = ?")
		args = append(args, p.IDCard)
		fields = append(fields, "id_card_encrypted = ?")
		args = append(args, p.IDCardEncrypted)
	}
	if p.Phone != "" {
		fields = append(fields, "phone = ?")
		args = append(args, p.Phone)
		fields = append(fields, "phone_encrypted = ?")
		args = append(args, p.PhoneEncrypted)
	}
	if p.Gender != "" {
		fields = append(fields, "gender = ?")
		args = append(args, p.Gender)
	}
	if p.Age >= 0 {
		fields = append(fields, "age = ?")
		args = append(args, p.Age)
	}
	if p.Address != "" {
		fields = append(fields, "address = ?")
		args = append(args, p.Address)
	}
	if len(fields) == 0 {
		return nil
	}
	args = append(args, p.ID)
	_, err := r.db.Exec("UPDATE patients SET "+strings.Join(fields, ", ")+" WHERE id = ?", args...)
	if err != nil {
		return fmt.Errorf("update patient: %w", err)
	}
	return nil
}

func (r *PatientRepository) Delete(id int64) error {
	_, err := r.db.Exec(`DELETE FROM patients WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete patient: %w", err)
	}
	return nil
}

func (r *PatientRepository) HasActiveOrders(patientID int64) (bool, error) {
	var count int
	// 依赖 orders 表，若不存在则视为无关联订单
	err := r.db.QueryRow(`SELECT COUNT(*) FROM orders WHERE patient_id = ? AND status = 'confirmed'`, patientID).Scan(&count)
	if err != nil {
		// orders 表可能尚未创建，视为无关联订单
		return false, nil
	}
	return count > 0, nil
}
