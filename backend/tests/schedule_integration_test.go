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

func setupScheduleIntegration(t *testing.T) *httptest.Server {
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
	r := router.Setup(database)
	return httptest.NewServer(r)
}

func TestIntegration_Schedule_Lifecycle(t *testing.T) {
	srv := setupScheduleIntegration(t)
	defer srv.Close()
	client := srv.Client()
	base := srv.URL
	token, _ := middleware.GenerateToken("admin-1")

	// Create schedule
	body := `{"date":"2026-04-29","department":"内科","doctor_name":"王医生","start_time":"09:00","end_time":"10:00","total_quota":20}`
	req, _ := http.NewRequest("POST", base+"/api/v1/admin/schedules", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("create status = %d, body: %s", resp.StatusCode, readBody(resp))
	}
	var createResp apiResp
	json.NewDecoder(resp.Body).Decode(&createResp)
	resp.Body.Close()
	var created struct {
		ID         int64  `json:"id"`
		Remaining  int    `json:"remaining"`
		Status     string `json:"status"`
		TotalQuota int    `json:"total_quota"`
	}
	json.Unmarshal(createResp.Data, &created)
	if created.ID == 0 {
		t.Fatal("expected id > 0")
	}
	if created.Remaining != 20 || created.Status != "available" {
		t.Fatalf("remaining=%d status=%s, want 20 available", created.Remaining, created.Status)
	}

	// List schedules
	req, _ = http.NewRequest("GET", base+"/api/v1/schedules?date=2026-04-29", nil)
	resp, _ = client.Do(req)
	var listResp apiResp
	json.NewDecoder(resp.Body).Decode(&listResp)
	resp.Body.Close()
	var listResult struct {
		Total int `json:"total"`
		List  []struct {
			ID int64 `json:"id"`
		} `json:"list"`
	}
	json.Unmarshal(listResp.Data, &listResult)
	if listResult.Total != 1 {
		t.Errorf("list total = %d, want 1", listResult.Total)
	}

	// Deduct
	req, _ = http.NewRequest("POST", base+"/api/v1/admin/schedules/"+strconv.FormatInt(created.ID, 10)+"/deduct", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, _ = client.Do(req)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("deduct status = %d", resp.StatusCode)
	}
	resp.Body.Close()

	// Verify remaining decreased
	req, _ = http.NewRequest("GET", base+"/api/v1/schedules/"+strconv.FormatInt(created.ID, 10), nil)
	resp, _ = client.Do(req)
	var getResp apiResp
	json.NewDecoder(resp.Body).Decode(&getResp)
	resp.Body.Close()
	var got struct {
		Remaining int    `json:"remaining"`
		Status    string `json:"status"`
	}
	json.Unmarshal(getResp.Data, &got)
	if got.Remaining != 19 {
		t.Errorf("remaining after deduct = %d, want 19", got.Remaining)
	}

	// Rollback
	req, _ = http.NewRequest("POST", base+"/api/v1/admin/schedules/"+strconv.FormatInt(created.ID, 10)+"/rollback", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, _ = client.Do(req)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("rollback status = %d", resp.StatusCode)
	}
	resp.Body.Close()

	// Update total_quota
	upd := `{"total_quota":25}`
	req, _ = http.NewRequest("PUT", base+"/api/v1/admin/schedules/"+strconv.FormatInt(created.ID, 10), bytes.NewBufferString(upd))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, _ = client.Do(req)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("update status = %d", resp.StatusCode)
	}
	resp.Body.Close()

	// Delete
	req, _ = http.NewRequest("DELETE", base+"/api/v1/admin/schedules/"+strconv.FormatInt(created.ID, 10), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, _ = client.Do(req)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("delete status = %d", resp.StatusCode)
	}
	resp.Body.Close()

	// Verify deletion
	req, _ = http.NewRequest("GET", base+"/api/v1/schedules/"+strconv.FormatInt(created.ID, 10), nil)
	resp, _ = client.Do(req)
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("get after delete status = %d, want 404", resp.StatusCode)
	}
	resp.Body.Close()
}

func TestIntegration_Schedule_DeductToFull(t *testing.T) {
	srv := setupScheduleIntegration(t)
	defer srv.Close()
	client := srv.Client()
	base := srv.URL
	token, _ := middleware.GenerateToken("admin-1")

	body := `{"date":"2026-04-29","department":"内科","doctor_name":"王医生","start_time":"09:00","end_time":"10:00","total_quota":1}`
	req, _ := http.NewRequest("POST", base+"/api/v1/admin/schedules", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, _ := client.Do(req)
	var createResp apiResp
	json.NewDecoder(resp.Body).Decode(&createResp)
	resp.Body.Close()
	var created struct{ ID int64 `json:"id"` }
	json.Unmarshal(createResp.Data, &created)

	// Deduct last one
	req, _ = http.NewRequest("POST", base+"/api/v1/admin/schedules/"+strconv.FormatInt(created.ID, 10)+"/deduct", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, _ = client.Do(req)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("deduct status = %d", resp.StatusCode)
	}
	resp.Body.Close()

	// Verify full
	req, _ = http.NewRequest("GET", base+"/api/v1/schedules/"+strconv.FormatInt(created.ID, 10), nil)
	resp, _ = client.Do(req)
	var getResp apiResp
	json.NewDecoder(resp.Body).Decode(&getResp)
	resp.Body.Close()
	var got struct {
		Remaining int    `json:"remaining"`
		Status    string `json:"status"`
	}
	json.Unmarshal(getResp.Data, &got)
	if got.Remaining != 0 || got.Status != "full" {
		t.Errorf("remaining=%d status=%s, want 0 full", got.Remaining, got.Status)
	}

	// Deduct again should fail
	req, _ = http.NewRequest("POST", base+"/api/v1/admin/schedules/"+strconv.FormatInt(created.ID, 10)+"/deduct", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, _ = client.Do(req)
	if resp.StatusCode != http.StatusConflict {
		t.Errorf("deduct again status = %d, want 409", resp.StatusCode)
	}
	resp.Body.Close()
}

func TestIntegration_Schedule_DeleteWithBookingRejected(t *testing.T) {
	srv := setupScheduleIntegration(t)
	defer srv.Close()
	client := srv.Client()
	base := srv.URL
	token, _ := middleware.GenerateToken("admin-1")

	body := `{"date":"2026-04-29","department":"内科","doctor_name":"王医生","start_time":"09:00","end_time":"10:00","total_quota":2}`
	req, _ := http.NewRequest("POST", base+"/api/v1/admin/schedules", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, _ := client.Do(req)
	var createResp apiResp
	json.NewDecoder(resp.Body).Decode(&createResp)
	resp.Body.Close()
	var created struct{ ID int64 `json:"id"` }
	json.Unmarshal(createResp.Data, &created)

	// Deduct to create booking
	req, _ = http.NewRequest("POST", base+"/api/v1/admin/schedules/"+strconv.FormatInt(created.ID, 10)+"/deduct", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, _ = client.Do(req)
	resp.Body.Close()

	// Delete should fail
	req, _ = http.NewRequest("DELETE", base+"/api/v1/admin/schedules/"+strconv.FormatInt(created.ID, 10), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, _ = client.Do(req)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("delete with booking status = %d, want 400", resp.StatusCode)
	}
	resp.Body.Close()
}

func readBody(resp *http.Response) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	resp.Body.Close()
	return buf.String()
}
