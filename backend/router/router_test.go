package router

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"clinic/db"
	"clinic/middleware"

	"github.com/gin-gonic/gin"
)

func setupRouter(t *testing.T) *gin.Engine {
	t.Helper()
	database, err := db.Open(":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	b, _ := os.ReadFile("../migrations/001_create_patients.sql")
	if err := db.ExecMigration(database, string(b)); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	b2, _ := os.ReadFile("../migrations/002_create_schedules.sql")
	if err := db.ExecMigration(database, string(b2)); err != nil {
		t.Fatalf("migrate schedules: %v", err)
	}
	b3, _ := os.ReadFile("../migrations/003_create_orders.sql")
	if err := db.ExecMigration(database, string(b3)); err != nil {
		t.Fatalf("migrate orders: %v", err)
	}
	return Setup(database)
}

func TestRouter_H5_CreatePatient(t *testing.T) {
	r := setupRouter(t)
	body := `{"name":"张三","id_card":"110101199001011234","phone":"13800138000","gender":"male","age":32}`
	req, _ := http.NewRequest("POST", "/api/v1/patients", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Visitor-Phone", "13800138000")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200, body: %s", w.Code, w.Body.String())
	}
	var resp struct {
		Code int `json:"code"`
		Data struct {
			ID int64 `json:"id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp.Data.ID == 0 {
		t.Fatal("expected id > 0")
	}
}

func TestRouter_H5_ListPatientsByVisitorPhone(t *testing.T) {
	r := setupRouter(t)
	body := `{"name":"张三","id_card":"110101199001011234","phone":"13800138000","gender":"male","age":32}`
	req1, _ := http.NewRequest("POST", "/api/v1/patients", bytes.NewBufferString(body))
	req1.Header.Set("Content-Type", "application/json")
	req1.Header.Set("X-Visitor-Phone", "13800138000")
	r.ServeHTTP(httptest.NewRecorder(), req1)

	req2, _ := http.NewRequest("GET", "/api/v1/patients", nil)
	req2.Header.Set("X-Visitor-Phone", "13800138000")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req2)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	var resp struct {
		Code int `json:"code"`
		Data struct {
			Total int `json:"total"`
		} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Data.Total != 1 {
		t.Errorf("total = %d, want 1", resp.Data.Total)
	}
}

func TestRouter_Admin_RejectWithoutToken(t *testing.T) {
	r := setupRouter(t)
	req, _ := http.NewRequest("GET", "/api/v1/admin/patients", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want 401", w.Code)
	}
}

func TestRouter_Admin_ListWithToken(t *testing.T) {
	r := setupRouter(t)
	token, _ := middleware.GenerateToken("admin-1")
	req, _ := http.NewRequest("GET", "/api/v1/admin/patients?page=1&page_size=10", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200, body: %s", w.Code, w.Body.String())
	}
}

func TestRouter_ScheduleRoutes_Registered(t *testing.T) {
	r := setupRouter(t)
	body := `{"date":"2026-04-29","department":"内科","doctor_name":"王医生","start_time":"09:00","end_time":"10:00","total_quota":20}`

	// H5 list schedules
	req, _ := http.NewRequest("GET", "/api/v1/schedules", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("list schedules status = %d, want 200", w.Code)
	}

	// Admin create schedule
	token, _ := middleware.GenerateToken("admin-1")
	req, _ = http.NewRequest("POST", "/api/v1/admin/schedules", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("create schedule status = %d, want 200, body: %s", w.Code, w.Body.String())
	}

	// Admin list schedules
	req, _ = http.NewRequest("GET", "/api/v1/admin/schedules", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("admin list schedules status = %d, want 200", w.Code)
	}
}

func TestRouter_RegistrationRoutes_Registered(t *testing.T) {
	r := setupRouter(t)

	// Submit without visitor phone should fail (middleware rejects)
	body := `{"schedule_id":1,"patient_id":1,"visitor_phone":"13800138000"}`
	req, _ := http.NewRequest("POST", "/api/v1/registrations", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	// Expect 400 because no X-Visitor-Phone header (VisitorPhone middleware will set empty)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("submit without visitor phone status = %d, want 400", w.Code)
	}

	// Submit with visitor phone
	req, _ = http.NewRequest("POST", "/api/v1/registrations", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Visitor-Phone", "13800138000")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	// Will fail due to missing seed data but verifies route is registered
	if w.Code != http.StatusBadRequest {
		t.Fatalf("submit with visitor phone status = %d, want 400 (patient not found)", w.Code)
	}
}
