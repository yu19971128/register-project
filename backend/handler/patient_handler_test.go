package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"clinic/db"
	"clinic/models"
	"clinic/repo"
	"clinic/service"

	"github.com/gin-gonic/gin"
)

type testResponse struct {
	Code    int             `json:"code"`
	Data    json.RawMessage `json:"data"`
	Message string          `json:"message"`
}

func setupPatientHandler(t *testing.T) *PatientHandler {
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
	svc := service.NewPatientService(r)
	return NewPatientHandler(svc)
}

func performRequest(h func(*gin.Context), c *gin.Context) {
	h(c)
}

func parseResp(t *testing.T, b *bytes.Buffer) testResponse {
	t.Helper()
	var resp testResponse
	if err := json.Unmarshal(b.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	return resp
}

func TestPatientHandler_Create(t *testing.T) {
	h := setupPatientHandler(t)
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	body := `{"name":"张三","id_card":"110101199001011234","phone":"13800138000","gender":"male","age":32}`
	c.Request, _ = http.NewRequest("POST", "/patients", bytes.NewBufferString(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("visitor_phone", "13800138000")

	performRequest(h.Create, c)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	resp := parseResp(t, w.Body)
	if resp.Code != 200 {
		t.Fatalf("resp.code = %d, want 200", resp.Code)
	}
	var p models.Patient
	if err := json.Unmarshal(resp.Data, &p); err != nil {
		t.Fatalf("unmarshal patient: %v", err)
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
}

func TestPatientHandler_ListByVisitorPhone(t *testing.T) {
	h := setupPatientHandler(t)
	gin.SetMode(gin.TestMode)
	_, _ = h.svc.CreatePatient(&models.Patient{Name: "张三", IDCard: "110101199001011234", Phone: "13800138000", Gender: "male", VisitorPhone: "13800138000"})
	_, _ = h.svc.CreatePatient(&models.Patient{Name: "李四", IDCard: "110101199001011235", Phone: "13900139000", Gender: "female", VisitorPhone: "13800138000"})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/patients", nil)
	c.Request.URL.RawQuery = "page=1&page_size=10"
	c.Set("visitor_phone", "13800138000")

	performRequest(h.List, c)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	resp := parseResp(t, w.Body)
	var result struct {
		Total int                `json:"total"`
		List  []*models.Patient `json:"list"`
	}
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		t.Fatalf("unmarshal list: %v", err)
	}
	if result.Total != 2 {
		t.Errorf("total = %d, want 2", result.Total)
	}
	if len(result.List) != 2 {
		t.Errorf("len(list) = %d, want 2", len(result.List))
	}
}

func TestPatientHandler_ListByKeyword(t *testing.T) {
	h := setupPatientHandler(t)
	gin.SetMode(gin.TestMode)
	_, _ = h.svc.CreatePatient(&models.Patient{Name: "张三", IDCard: "110101199001011234", Phone: "13800138000", Gender: "male", VisitorPhone: "13800138000"})
	_, _ = h.svc.CreatePatient(&models.Patient{Name: "李四", IDCard: "110101199001011235", Phone: "13900139000", Gender: "female", VisitorPhone: "13800138000"})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/patients", nil)
	c.Request.URL.RawQuery = "keyword=张三&page=1&page_size=10"

	performRequest(h.List, c)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	resp := parseResp(t, w.Body)
	var result struct {
		Total int                `json:"total"`
		List  []*models.Patient `json:"list"`
	}
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		t.Fatalf("unmarshal list: %v", err)
	}
	if result.Total != 1 {
		t.Errorf("total = %d, want 1", result.Total)
	}
	if len(result.List) != 1 || result.List[0].Name != "张三" {
		t.Errorf("unexpected list result")
	}
}

func TestPatientHandler_Get(t *testing.T) {
	h := setupPatientHandler(t)
	gin.SetMode(gin.TestMode)
	p, _ := h.svc.CreatePatient(&models.Patient{Name: "张三", IDCard: "110101199001011234", Phone: "13800138000", Gender: "male", VisitorPhone: "13800138000"})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/patients/"+strconv.FormatInt(p.ID, 10), nil)
	c.Params = gin.Params{{Key: "id", Value: strconv.FormatInt(p.ID, 10)}}

	performRequest(h.Get, c)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	resp := parseResp(t, w.Body)
	var got models.Patient
	if err := json.Unmarshal(resp.Data, &got); err != nil {
		t.Fatalf("unmarshal patient: %v", err)
	}
	if got.Name != "张三" {
		t.Errorf("name = %s, want 张三", got.Name)
	}
}

func TestPatientHandler_Get_NotFound(t *testing.T) {
	h := setupPatientHandler(t)
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/patients/9999", nil)
	c.Params = gin.Params{{Key: "id", Value: "9999"}}

	performRequest(h.Get, c)
	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", w.Code)
	}
	resp := parseResp(t, w.Body)
	if resp.Code != http.StatusNotFound {
		t.Errorf("resp.code = %d, want 404", resp.Code)
	}
}

func TestPatientHandler_Update(t *testing.T) {
	h := setupPatientHandler(t)
	gin.SetMode(gin.TestMode)
	p, _ := h.svc.CreatePatient(&models.Patient{Name: "张三", IDCard: "110101199001011234", Phone: "13800138000", Gender: "male", VisitorPhone: "13800138000"})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	body := `{"name":"张三Updated"}`
	c.Request, _ = http.NewRequest("PUT", "/patients/"+strconv.FormatInt(p.ID, 10), bytes.NewBufferString(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: strconv.FormatInt(p.ID, 10)}}

	performRequest(h.Update, c)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}

	// verify
	got, _ := h.svc.GetPatient(p.ID)
	if got == nil || got.Name != "张三Updated" {
		t.Errorf("name after update = %s, want 张三Updated", got.Name)
	}
}

func TestPatientHandler_Delete(t *testing.T) {
	h := setupPatientHandler(t)
	gin.SetMode(gin.TestMode)
	p, _ := h.svc.CreatePatient(&models.Patient{Name: "张三", IDCard: "110101199001011234", Phone: "13800138000", Gender: "male", VisitorPhone: "13800138000"})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("DELETE", "/patients/"+strconv.FormatInt(p.ID, 10), nil)
	c.Params = gin.Params{{Key: "id", Value: strconv.FormatInt(p.ID, 10)}}

	performRequest(h.Delete, c)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}

	got, _ := h.svc.GetPatient(p.ID)
	if got != nil {
		t.Error("expected patient deleted")
	}
}
