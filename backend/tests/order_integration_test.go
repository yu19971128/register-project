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

func setupOrderIntegration(t *testing.T) *httptest.Server {
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

func seedPatientAndSchedule(t *testing.T, srv *httptest.Server, client *http.Client, token, visitorPhone string) (patientID, scheduleID int64) {
	base := srv.URL

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
	scheduleBody := `{"date":"2099-04-29","department":"内科","doctor_name":"王医生","start_time":"09:00","end_time":"10:00","total_quota":10}`
	req, _ = http.NewRequest("POST", base+"/api/v1/admin/schedules", bytes.NewBufferString(scheduleBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, _ = client.Do(req)
	var scheduleResp apiResp
	json.NewDecoder(resp.Body).Decode(&scheduleResp)
	resp.Body.Close()
	var scheduleData struct{ ID int64 `json:"id"` }
	json.Unmarshal(scheduleResp.Data, &scheduleData)

	return patientData.ID, scheduleData.ID
}

func TestIntegration_Order_List_Get_Cancel(t *testing.T) {
	srv := setupOrderIntegration(t)
	defer srv.Close()
	client := srv.Client()
	base := srv.URL
	visitorPhone := "13800138000"
	token, _ := middleware.GenerateToken("admin-1")

	pid, sid := seedPatientAndSchedule(t, srv, client, token, visitorPhone)

	// Submit registration to create order
	regBody := `{"schedule_id":` + strconv.FormatInt(sid, 10) + `,"patient_id":` + strconv.FormatInt(pid, 10) + `,"visitor_phone":"` + visitorPhone + `"}`
	req, _ := http.NewRequest("POST", base+"/api/v1/registrations", bytes.NewBufferString(regBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Visitor-Phone", visitorPhone)
	resp, _ := client.Do(req)
	var regResp apiResp
	json.NewDecoder(resp.Body).Decode(&regResp)
	resp.Body.Close()
	var regResult struct{ OrderNo string `json:"order_no"` }
	json.Unmarshal(regResp.Data, &regResult)

	// H5 list orders
	req, _ = http.NewRequest("GET", base+"/api/v1/orders", nil)
	req.Header.Set("X-Visitor-Phone", visitorPhone)
	resp, _ = client.Do(req)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("h5 list orders status = %d, body: %s", resp.StatusCode, readBody(resp))
	}
	var listResp apiResp
	json.NewDecoder(resp.Body).Decode(&listResp)
	resp.Body.Close()
	var listResult struct {
		Total int `json:"total"`
		List  []struct {
			ID      int64  `json:"id"`
			OrderNo string `json:"order_no"`
			Status  string `json:"status"`
		} `json:"list"`
	}
	json.Unmarshal(listResp.Data, &listResult)
	if listResult.Total != 1 {
		t.Errorf("list total = %d, want 1", listResult.Total)
	}
	if listResult.List[0].Status != "confirmed" {
		t.Errorf("order status = %s, want confirmed", listResult.List[0].Status)
	}
	orderID := listResult.List[0].ID

	// Get order detail
	req, _ = http.NewRequest("GET", base+"/api/v1/orders/"+strconv.FormatInt(orderID, 10), nil)
	req.Header.Set("X-Visitor-Phone", visitorPhone)
	resp, _ = client.Do(req)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("get order status = %d, body: %s", resp.StatusCode, readBody(resp))
	}
	var getResp apiResp
	json.NewDecoder(resp.Body).Decode(&getResp)
	resp.Body.Close()
	var getResult struct {
		OrderNo string `json:"order_no"`
		Patient struct {
			Name string `json:"name"`
		} `json:"patient"`
	}
	json.Unmarshal(getResp.Data, &getResult)
	if getResult.Patient.Name != "张三" {
		t.Errorf("patient name = %s, want 张三", getResult.Patient.Name)
	}

	// Cancel order
	cancelBody := `{"reason":"个人原因"}`
	req, _ = http.NewRequest("PUT", base+"/api/v1/orders/"+strconv.FormatInt(orderID, 10)+"/cancel", bytes.NewBufferString(cancelBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Visitor-Phone", visitorPhone)
	resp, _ = client.Do(req)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("cancel status = %d, body: %s", resp.StatusCode, readBody(resp))
	}
	var cancelResp apiResp
	json.NewDecoder(resp.Body).Decode(&cancelResp)
	resp.Body.Close()
	var cancelResult struct{ Status string `json:"status"` }
	json.Unmarshal(cancelResp.Data, &cancelResult)
	if cancelResult.Status != "cancelled" {
		t.Errorf("cancel status = %s, want cancelled", cancelResult.Status)
	}

	// Verify schedule rolled back
	req, _ = http.NewRequest("GET", base+"/api/v1/schedules/"+strconv.FormatInt(sid, 10), nil)
	resp, _ = client.Do(req)
	var sResp apiResp
	json.NewDecoder(resp.Body).Decode(&sResp)
	resp.Body.Close()
	var sData struct{ Remaining int `json:"remaining"` }
	json.Unmarshal(sResp.Data, &sData)
	if sData.Remaining != 10 {
		t.Errorf("schedule remaining after cancel = %d, want 10", sData.Remaining)
	}
}

func TestIntegration_Order_Change(t *testing.T) {
	srv := setupOrderIntegration(t)
	defer srv.Close()
	client := srv.Client()
	base := srv.URL
	visitorPhone := "13800138000"
	token, _ := middleware.GenerateToken("admin-1")

	pid, sid := seedPatientAndSchedule(t, srv, client, token, visitorPhone)

	// Create second schedule
	scheduleBody2 := `{"date":"2099-04-29","department":"外科","doctor_name":"李医生","start_time":"14:00","end_time":"15:00","total_quota":5}`
	req, _ := http.NewRequest("POST", base+"/api/v1/admin/schedules", bytes.NewBufferString(scheduleBody2))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, _ := client.Do(req)
	var s2Resp apiResp
	json.NewDecoder(resp.Body).Decode(&s2Resp)
	resp.Body.Close()
	var s2Data struct{ ID int64 `json:"id"` }
	json.Unmarshal(s2Resp.Data, &s2Data)

	// Submit registration
	regBody := `{"schedule_id":` + strconv.FormatInt(sid, 10) + `,"patient_id":` + strconv.FormatInt(pid, 10) + `,"visitor_phone":"` + visitorPhone + `"}`
	req, _ = http.NewRequest("POST", base+"/api/v1/registrations", bytes.NewBufferString(regBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Visitor-Phone", visitorPhone)
	resp, _ = client.Do(req)
	resp.Body.Close()

	// Get order ID
	req, _ = http.NewRequest("GET", base+"/api/v1/orders", nil)
	req.Header.Set("X-Visitor-Phone", visitorPhone)
	resp, _ = client.Do(req)
	var listResp apiResp
	json.NewDecoder(resp.Body).Decode(&listResp)
	resp.Body.Close()
	var listResult struct {
		List []struct {
			ID int64 `json:"id"`
		} `json:"list"`
	}
	json.Unmarshal(listResp.Data, &listResult)
	orderID := listResult.List[0].ID

	// Change order
	changeBody := `{"new_schedule_id":` + strconv.FormatInt(s2Data.ID, 10) + `}`
	req, _ = http.NewRequest("PUT", base+"/api/v1/orders/"+strconv.FormatInt(orderID, 10)+"/change", bytes.NewBufferString(changeBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Visitor-Phone", visitorPhone)
	resp, _ = client.Do(req)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("change status = %d, body: %s", resp.StatusCode, readBody(resp))
	}
	var changeResp apiResp
	json.NewDecoder(resp.Body).Decode(&changeResp)
	resp.Body.Close()
	var changeResult struct {
		Status   string `json:"status"`
		Schedule struct {
			DoctorName string `json:"doctor_name"`
		} `json:"schedule"`
	}
	json.Unmarshal(changeResp.Data, &changeResult)
	if changeResult.Status != "confirmed" {
		t.Errorf("change result status = %s, want confirmed", changeResult.Status)
	}
	if changeResult.Schedule.DoctorName != "李医生" {
		t.Errorf("change result doctor = %s, want 李医生", changeResult.Schedule.DoctorName)
	}

	// Verify old schedule rolled back
	req, _ = http.NewRequest("GET", base+"/api/v1/schedules/"+strconv.FormatInt(sid, 10), nil)
	resp, _ = client.Do(req)
	var sResp apiResp
	json.NewDecoder(resp.Body).Decode(&sResp)
	resp.Body.Close()
	var sData struct{ Remaining int `json:"remaining"` }
	json.Unmarshal(sResp.Data, &sData)
	if sData.Remaining != 10 {
		t.Errorf("old schedule remaining after change = %d, want 10", sData.Remaining)
	}

	// Verify new schedule deducted
	req, _ = http.NewRequest("GET", base+"/api/v1/schedules/"+strconv.FormatInt(s2Data.ID, 10), nil)
	resp, _ = client.Do(req)
	var nsResp apiResp
	json.NewDecoder(resp.Body).Decode(&nsResp)
	resp.Body.Close()
	var nsData struct{ Remaining int `json:"remaining"` }
	json.Unmarshal(nsResp.Data, &nsData)
	if nsData.Remaining != 4 {
		t.Errorf("new schedule remaining after change = %d, want 4", nsData.Remaining)
	}
}

func TestIntegration_Order_AdminListAndCancel(t *testing.T) {
	srv := setupOrderIntegration(t)
	defer srv.Close()
	client := srv.Client()
	base := srv.URL
	visitorPhone := "13800138000"
	token, _ := middleware.GenerateToken("admin-1")

	pid, sid := seedPatientAndSchedule(t, srv, client, token, visitorPhone)

	// Submit registration
	regBody := `{"schedule_id":` + strconv.FormatInt(sid, 10) + `,"patient_id":` + strconv.FormatInt(pid, 10) + `,"visitor_phone":"` + visitorPhone + `"}`
	req, _ := http.NewRequest("POST", base+"/api/v1/registrations", bytes.NewBufferString(regBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Visitor-Phone", visitorPhone)
	resp, _ := client.Do(req)
	resp.Body.Close()

	// Admin list orders
	req, _ = http.NewRequest("GET", base+"/api/v1/admin/orders", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, _ = client.Do(req)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("admin list status = %d, body: %s", resp.StatusCode, readBody(resp))
	}
	var listResp apiResp
	json.NewDecoder(resp.Body).Decode(&listResp)
	resp.Body.Close()
	var listResult struct {
		Total int `json:"total"`
	}
	json.Unmarshal(listResp.Data, &listResult)
	if listResult.Total != 1 {
		t.Errorf("admin list total = %d, want 1", listResult.Total)
	}

	// Admin cancel
	req, _ = http.NewRequest("GET", base+"/api/v1/orders", nil)
	req.Header.Set("X-Visitor-Phone", visitorPhone)
	resp, _ = client.Do(req)
	var h5ListResp apiResp
	json.NewDecoder(resp.Body).Decode(&h5ListResp)
	resp.Body.Close()
	var h5ListResult struct {
		List []struct {
			ID int64 `json:"id"`
		} `json:"list"`
	}
	json.Unmarshal(h5ListResp.Data, &h5ListResult)
	orderID := h5ListResult.List[0].ID

	cancelBody := `{"reason":"医生停诊"}`
	req, _ = http.NewRequest("PUT", base+"/api/v1/admin/orders/"+strconv.FormatInt(orderID, 10)+"/cancel", bytes.NewBufferString(cancelBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, _ = client.Do(req)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("admin cancel status = %d, body: %s", resp.StatusCode, readBody(resp))
	}
	resp.Body.Close()
}
