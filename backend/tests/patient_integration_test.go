package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"clinic/db"
	"clinic/middleware"
	"clinic/router"
)

func setupIntegration(t *testing.T) *httptest.Server {
	t.Helper()
	database, err := db.Open(":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	b, _ := os.ReadFile("../migrations/001_create_patients.sql")
	if err := db.ExecMigration(database, string(b)); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	r := router.Setup(database)
	return httptest.NewServer(r)
}

type apiResp struct {
	Code    int             `json:"code"`
	Data    json.RawMessage `json:"data"`
	Message string          `json:"message"`
}

func TestIntegration_H5_PatientLifecycle(t *testing.T) {
	srv := setupIntegration(t)
	defer srv.Close()
	client := srv.Client()
	base := srv.URL

	// Create
	body := `{"name":"张三","id_card":"110101199001011234","phone":"13800138000","gender":"male","age":32}`
	req, _ := http.NewRequest("POST", base+"/api/v1/patients", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Visitor-Phone", "13800138000")
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("create status = %d", resp.StatusCode)
	}
	var createResp apiResp
	json.NewDecoder(resp.Body).Decode(&createResp)
	resp.Body.Close()
	var created struct {
		ID    int64  `json:"id"`
		Name  string `json:"name"`
		Phone string `json:"phone"`
	}
	json.Unmarshal(createResp.Data, &created)
	if created.ID == 0 {
		t.Fatal("expected id > 0")
	}
	if created.Phone != "138****8000" {
		t.Errorf("phone mask = %s, want 138****8000", created.Phone)
	}

	// List by visitor phone
	req, _ = http.NewRequest("GET", base+"/api/v1/patients", nil)
	req.Header.Set("X-Visitor-Phone", "13800138000")
	resp, _ = client.Do(req)
	var listResp apiResp
	json.NewDecoder(resp.Body).Decode(&listResp)
	resp.Body.Close()
	var listResult struct {
		Total int `json:"total"`
		List  []struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
		} `json:"list"`
	}
	json.Unmarshal(listResp.Data, &listResult)
	if listResult.Total != 1 {
		t.Errorf("list total = %d, want 1", listResult.Total)
	}

	// Get
	req, _ = http.NewRequest("GET", base+"/api/v1/patients/"+strconv.FormatInt(created.ID, 10), nil)
	resp, _ = client.Do(req)
	var getResp apiResp
	json.NewDecoder(resp.Body).Decode(&getResp)
	resp.Body.Close()
	if getResp.Code != 200 {
		t.Errorf("get code = %d", getResp.Code)
	}

	// Update
	upd := `{"name":"张三Updated"}`
	req, _ = http.NewRequest("PUT", base+"/api/v1/patients/"+strconv.FormatInt(created.ID, 10), bytes.NewBufferString(upd))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = client.Do(req)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("update status = %d", resp.StatusCode)
	}
	resp.Body.Close()

	// Delete
	req, _ = http.NewRequest("DELETE", base+"/api/v1/patients/"+strconv.FormatInt(created.ID, 10), nil)
	resp, _ = client.Do(req)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("delete status = %d", resp.StatusCode)
	}
	resp.Body.Close()

	// Verify deletion
	req, _ = http.NewRequest("GET", base+"/api/v1/patients/"+strconv.FormatInt(created.ID, 10), nil)
	resp, _ = client.Do(req)
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("get after delete status = %d, want 404", resp.StatusCode)
	}
	resp.Body.Close()
}

func TestIntegration_Admin_RejectWithoutToken(t *testing.T) {
	srv := setupIntegration(t)
	defer srv.Close()
	client := srv.Client()

	req, _ := http.NewRequest("GET", srv.URL+"/api/v1/admin/patients", nil)
	resp, _ := client.Do(req)
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("status = %d, want 401", resp.StatusCode)
	}
	resp.Body.Close()
}

func TestIntegration_Admin_ListWithKeyword(t *testing.T) {
	srv := setupIntegration(t)
	defer srv.Close()
	client := srv.Client()
	base := srv.URL

	// Seed data via H5 endpoint
	body := `{"name":"李四","id_card":"110101199001011235","phone":"13900139000","gender":"female","age":28}`
	req, _ := http.NewRequest("POST", base+"/api/v1/patients", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Visitor-Phone", "13900139000")
	resp, _ := client.Do(req)
	resp.Body.Close()

	token, _ := middleware.GenerateToken("admin-1")
	req, _ = http.NewRequest("GET", base+"/api/v1/admin/patients?keyword=李四", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, _ = client.Do(req)
	var listResp apiResp
	json.NewDecoder(resp.Body).Decode(&listResp)
	resp.Body.Close()
	var result struct {
		Total int `json:"total"`
	}
	json.Unmarshal(listResp.Data, &result)
	if result.Total != 1 {
		t.Errorf("total = %d, want 1", result.Total)
	}
}
