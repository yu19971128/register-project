package service

import (
	"os"
	"testing"

	"clinic/db"
	"clinic/models"
	"clinic/repo"
)

func setupPatientService(t *testing.T) *PatientService {
	t.Helper()
	database, err := db.Open(":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	b, _ := os.ReadFile("../migrations/001_create_patients.sql")
	if err := db.ExecMigration(database, string(b)); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	r := repo.NewPatientRepository(database)
	return NewPatientService(r)
}

func TestPatientService_Create_With_Encryption(t *testing.T) {
	svc := setupPatientService(t)
	p, err := svc.CreatePatient(&models.Patient{
		Name:         "张三",
		IDCard:       "110101199001011234",
		Phone:        "13800138000",
		Gender:       "male",
		Age:          32,
		VisitorPhone: "13800138000",
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if p.ID == 0 {
		t.Fatal("expected id > 0")
	}
	if p.IDCard != "110101********1234" {
		t.Errorf("id_card mask = %s, want 110101********1234", p.IDCard)
	}
	if p.Phone != "138****8000" {
		t.Errorf("phone mask = %s, want 138****8000", p.Phone)
	}
	if p.IDCardEncrypted == "" {
		t.Error("expected encrypted id_card")
	}
}

func TestPatientService_List_Get_Update_Delete(t *testing.T) {
	svc := setupPatientService(t)
	p1, _ := svc.CreatePatient(&models.Patient{Name: "张三", IDCard: "110101199001011234", Phone: "13800138000", Gender: "male", VisitorPhone: "13800138000"})
	p2, _ := svc.CreatePatient(&models.Patient{Name: "李四", IDCard: "110101199001011235", Phone: "13900139000", Gender: "female", VisitorPhone: "13800138000"})

	list, total, err := svc.ListPatientsByVisitorPhone("13800138000", 1, 10)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if total != 2 {
		t.Errorf("total = %d, want 2", total)
	}
	if len(list) != 2 {
		t.Errorf("len(list) = %d, want 2", len(list))
	}

	got, err := svc.GetPatient(p1.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got == nil {
		t.Fatal("expected patient, got nil")
	}

	if err := svc.UpdatePatient(&models.Patient{ID: p1.ID, Name: "张三Updated"}); err != nil {
		t.Fatalf("update: %v", err)
	}

	if err := svc.DeletePatient(p2.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}

	_, err = svc.GetPatient(p2.ID)
	if err != nil {
		t.Fatalf("get after delete: %v", err)
	}
}
