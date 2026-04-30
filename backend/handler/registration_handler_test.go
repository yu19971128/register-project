package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"clinic/db"
	"clinic/models"
	"clinic/repo"
	"clinic/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupRegistrationHandler(t *testing.T) (*RegistrationHandler, *gin.Engine, *sql.DB) {
	database, err := db.Open(":memory:")
	require.NoError(t, err)
	b1, _ := os.ReadFile("../migrations/001_create_patients.sql")
	require.NoError(t, db.ExecMigration(database, string(b1)))
	b2, _ := os.ReadFile("../migrations/002_create_schedules.sql")
	require.NoError(t, db.ExecMigration(database, string(b2)))
	b3, _ := os.ReadFile("../migrations/003_create_orders.sql")
	require.NoError(t, db.ExecMigration(database, string(b3)))

	patientRepo := repo.NewPatientRepository(database)
	scheduleRepo := repo.NewScheduleRepository(database)
	orderRepo := repo.NewOrderRepository(database)
	regSvc := service.NewRegistrationService(database, patientRepo, scheduleRepo, orderRepo)
	h := NewRegistrationHandler(regSvc)

	gin.SetMode(gin.TestMode)
	g := gin.New()
	g.Use(func(c *gin.Context) {
		if vp := c.GetHeader("X-Visitor-Phone"); vp != "" {
			c.Set("visitor_phone", vp)
		}
		c.Next()
	})
	g.POST("/api/v1/registrations", h.Submit)
	g.GET("/api/v1/registrations/ticket/:order_no", h.GetTicket)
	return h, g, database
}

func seedForRegistration(t *testing.T, database *sql.DB) (int64, int64) {
	patientRepo := repo.NewPatientRepository(database)
	scheduleRepo := repo.NewScheduleRepository(database)
	pid, _ := patientRepo.Create(&models.Patient{
		Name: "张三", IDCard: "110101********1234", Phone: "138****8888",
		IDCardEncrypted: "enc", PhoneEncrypted: "enc", Gender: "male", Age: 32,
		VisitorPhone: "13800138000",
	})
	sid, _ := scheduleRepo.Create(&models.Schedule{
		Date: "2026-04-29", Department: "内科", DoctorName: "王医生",
		StartTime: "09:00", EndTime: "10:00", TotalQuota: 20, Remaining: 20, Status: "available",
	})
	return pid, sid
}

func TestRegistrationHandler_Submit_And_GetTicket(t *testing.T) {
	_, r, database := setupRegistrationHandler(t)
	seedForRegistration(t, database)

	body, _ := json.Marshal(map[string]interface{}{
		"schedule_id": 1, "patient_id": 1, "visitor_phone": "13800138000",
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/registrations", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Visitor-Phone", "13800138000")
	r.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("submit status = %d, body: %s", w.Code, w.Body.String())
	}
	var resp Response
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	resultMap, _ := resp.Data.(map[string]interface{})
	orderNo, _ := resultMap["order_no"].(string)
	assert.Contains(t, orderNo, "GH")

	// Get ticket
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/registrations/ticket/"+orderNo, nil)
	req.Header.Set("X-Visitor-Phone", "13800138000")
	r.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("ticket status = %d, body: %s", w.Code, w.Body.String())
	}
	var ticketResp Response
	_ = json.Unmarshal(w.Body.Bytes(), &ticketResp)
	ticketMap, _ := ticketResp.Data.(map[string]interface{})
	assert.Equal(t, "张三", ticketMap["patient_name"])
}

func TestRegistrationHandler_Submit_InvalidSchedule(t *testing.T) {
	_, r, _ := setupRegistrationHandler(t)
	body, _ := json.Marshal(map[string]interface{}{
		"schedule_id": 9999, "patient_id": 1, "visitor_phone": "13800138000",
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/registrations", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Visitor-Phone", "13800138000")
	r.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func toJSON(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
}
