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

func setupRegistrationIntegration(t *testing.T) *httptest.Server {
	t.Helper()
	database, err := db.Open(":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	b1, _ := os.ReadFile("../migrations/001_create_patients.sql")
	if err := db.ExecMigration(database, string(b1)); err != nil {
		t.Fatalf("migrate patients: %v", err)
	}
	b2, _ := os.ReadFile("../migrations/002_create_schedules.sql")
	if err := db.ExecMigration(database, string(b2)); err != nil {
		t.Fatalf("migrate schedules: %v", err)
	}
	b3, _ := os.ReadFile("../migrations/003_create_orders.sql")
	if err := db.ExecMigration(database, string(b3)); err != nil {
		t.Fatalf("migrate orders: %v", err)
	}
	b4, _ := os.ReadFile("../migrations/004_alter_orders.sql")
	if err := db.ExecMigration(database, string(b4)); err != nil {
		t.Fatalf("migrate orders alter: %v", err)
	}
	r := router.Setup(database)
	return httptest.NewServer(r)
}

func TestIntegration_Registration_Submit_And_GetTicket(t *testing.T) {
	srv := setupRegistrationIntegration(t)
	defer srv.Close()
	client := srv.Client()
	base := srv.URL
	visitorPhone := "13800138000"

	// Seed patient
	patientBody := `{"name":"张三","id_card":"110101199001011234","phone":"13800138000","gender":"male","age":32}`
	req, _ := http.NewRequest("POST", base+"/api/v1/patients", bytes.NewBufferString(patientBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Visitor-Phone", visitorPhone)
	resp, _ := client.Do(req)
	var patientResp apiResp
	json.NewDecoder(resp.Body).Decode(&patientResp)
	resp.Body.Close()
	var patientData struct{ ID int64 `json:"id"` }
	json.Unmarshal(patientResp.Data, &patientData)

	// Seed schedule via admin
	token, _ := middleware.GenerateToken("admin-1")
	scheduleBody := `{"date":"2026-04-29","department":"内科","doctor_name":"王医生","start_time":"09:00","end_time":"10:00","total_quota":2}`
	req, _ = http.NewRequest("POST", base+"/api/v1/admin/schedules", bytes.NewBufferString(scheduleBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, _ = client.Do(req)
	var scheduleResp apiResp
	json.NewDecoder(resp.Body).Decode(&scheduleResp)
	resp.Body.Close()
	var scheduleData struct{ ID int64 `json:"id"` }
	json.Unmarshal(scheduleResp.Data, &scheduleData)

	// Submit registration
	regBody := `{"schedule_id":` + strconv.FormatInt(scheduleData.ID, 10) + `,"patient_id":` + strconv.FormatInt(patientData.ID, 10) + `,"visitor_phone":"` + visitorPhone + `"}`
	req, _ = http.NewRequest("POST", base+"/api/v1/registrations", bytes.NewBufferString(regBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Visitor-Phone", visitorPhone)
	resp, _ = client.Do(req)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("submit status = %d, body: %s", resp.StatusCode, readBody(resp))
	}
	var regResp apiResp
	json.NewDecoder(resp.Body).Decode(&regResp)
	resp.Body.Close()
	var regResult struct {
		OrderNo string `json:"order_no"`
		Status  string `json:"status"`
	}
	json.Unmarshal(regResp.Data, &regResult)
	if regResult.OrderNo == "" {
		t.Fatal("expected order_no")
	}
	if regResult.Status != "confirmed" {
		t.Errorf("status = %s, want confirmed", regResult.Status)
	}

	// Get ticket
	req, _ = http.NewRequest("GET", base+"/api/v1/registrations/ticket/"+regResult.OrderNo, nil)
	req.Header.Set("X-Visitor-Phone", visitorPhone)
	resp, _ = client.Do(req)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("ticket status = %d", resp.StatusCode)
	}
	var ticketResp apiResp
	json.NewDecoder(resp.Body).Decode(&ticketResp)
	resp.Body.Close()
	var ticketData struct {
		PatientName string `json:"patient_name"`
		OrderNo     string `json:"order_no"`
	}
	json.Unmarshal(ticketResp.Data, &ticketData)
	if ticketData.PatientName != "张三" {
		t.Errorf("patient_name = %s, want 张三", ticketData.PatientName)
	}

	// Verify schedule remaining decreased
	req, _ = http.NewRequest("GET", base+"/api/v1/schedules/"+strconv.FormatInt(scheduleData.ID, 10), nil)
	resp, _ = client.Do(req)
	var getScheduleResp apiResp
	json.NewDecoder(resp.Body).Decode(&getScheduleResp)
	resp.Body.Close()
	var sData struct{ Remaining int `json:"remaining"` }
	json.Unmarshal(getScheduleResp.Data, &sData)
	if sData.Remaining != 1 {
		t.Errorf("remaining after registration = %d, want 1", sData.Remaining)
	}
}

func TestIntegration_Registration_DuplicateRejected(t *testing.T) {
	srv := setupRegistrationIntegration(t)
	defer srv.Close()
	client := srv.Client()
	base := srv.URL
	visitorPhone := "13800138000"

	patientBody := `{"name":"张三","id_card":"110101199001011234","phone":"13800138000","gender":"male","age":32}`
	req, _ := http.NewRequest("POST", base+"/api/v1/patients", bytes.NewBufferString(patientBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Visitor-Phone", visitorPhone)
	resp, _ := client.Do(req)
	var patientResp apiResp
	json.NewDecoder(resp.Body).Decode(&patientResp)
	resp.Body.Close()
	var patientData struct{ ID int64 `json:"id"` }
	json.Unmarshal(patientResp.Data, &patientData)

	token, _ := middleware.GenerateToken("admin-1")
	scheduleBody := `{"date":"2026-04-29","department":"内科","doctor_name":"王医生","start_time":"09:00","end_time":"10:00","total_quota":2}`
	req, _ = http.NewRequest("POST", base+"/api/v1/admin/schedules", bytes.NewBufferString(scheduleBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, _ = client.Do(req)
	var scheduleResp apiResp
	json.NewDecoder(resp.Body).Decode(&scheduleResp)
	resp.Body.Close()
	var scheduleData struct{ ID int64 `json:"id"` }
	json.Unmarshal(scheduleResp.Data, &scheduleData)

	regBody := `{"schedule_id":` + strconv.FormatInt(scheduleData.ID, 10) + `,"patient_id":` + strconv.FormatInt(patientData.ID, 10) + `,"visitor_phone":"` + visitorPhone + `"}`
	req, _ = http.NewRequest("POST", base+"/api/v1/registrations", bytes.NewBufferString(regBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Visitor-Phone", visitorPhone)
	resp, _ = client.Do(req)
	resp.Body.Close()

	// Duplicate
	req, _ = http.NewRequest("POST", base+"/api/v1/registrations", bytes.NewBufferString(regBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Visitor-Phone", visitorPhone)
	resp, _ = client.Do(req)
	if resp.StatusCode != http.StatusTooManyRequests {
		t.Errorf("duplicate submit status = %d, want 429", resp.StatusCode)
	}
	resp.Body.Close()
}
