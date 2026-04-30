package repo

import (
	"os"
	"testing"

	"clinic/db"
	"clinic/models"
)

func setupTestDB(t *testing.T) *PatientRepository {
	t.Helper()
	database, err := db.Open(":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	b, _ := os.ReadFile("../migrations/001_create_patients.sql")
	if err := db.ExecMigration(database, string(b)); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return NewPatientRepository(database)
}

func TestPatientRepository_Create_Get(t *testing.T) {
	r := setupTestDB(t)
	p := &models.Patient{
		Name:            "张三",
		IDCard:          "110101********1234",
		IDCardEncrypted: "encrypted_1234",
		Phone:           "138****8888",
		PhoneEncrypted:  "encrypted_phone",
		Gender:          "male",
		Age:             32,
		Address:         "北京市朝阳区",
		VisitorPhone:    "13800138000",
	}
	id, err := r.Create(p)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if id == 0 {
		t.Fatal("expected id > 0")
	}

	got, err := r.GetByID(id)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got == nil {
		t.Fatal("expected patient, got nil")
	}
	if got.Name != "张三" {
		t.Errorf("name = %s, want 张三", got.Name)
	}
}

func TestPatientRepository_ListByVisitorPhone(t *testing.T) {
	r := setupTestDB(t)
	for i := 0; i < 3; i++ {
		_, err := r.Create(&models.Patient{Name: "P" + string(rune('0'+i)), VisitorPhone: "13800138000", IDCardEncrypted: "e" + string(rune('0'+i)), PhoneEncrypted: "p" + string(rune('0'+i)), Gender: "male"})
		if err != nil {
			t.Fatalf("create %d: %v", i, err)
		}
	}

	list, total, err := r.ListByVisitorPhone("13800138000", 0, 10)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if total != 3 {
		t.Errorf("total = %d, want 3", total)
	}
	if len(list) != 3 {
		t.Errorf("len(list) = %d, want 3", len(list))
	}
}

func TestPatientRepository_Update_Delete(t *testing.T) {
	r := setupTestDB(t)
	p := &models.Patient{Name: "张三", VisitorPhone: "13800138000", IDCardEncrypted: "e1", PhoneEncrypted: "p1", Gender: "male"}
	id, err := r.Create(p)
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	if err := r.Update(&models.Patient{ID: id, Name: "张三Updated"}); err != nil {
		t.Fatalf("update: %v", err)
	}

	got, err := r.GetByID(id)
	if err != nil {
		t.Fatalf("get after update: %v", err)
	}
	if got == nil {
		t.Fatal("expected patient after update, got nil")
	}
	if got.Name != "张三Updated" {
		t.Errorf("name = %s, want 张三Updated", got.Name)
	}

	if err := r.Delete(id); err != nil {
		t.Fatalf("delete: %v", err)
	}

	got, _ = r.GetByID(id)
	if got != nil {
		t.Error("expected nil after delete")
	}
}
