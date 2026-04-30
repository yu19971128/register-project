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
