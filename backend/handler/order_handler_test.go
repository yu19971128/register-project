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

func setupOrderHandler(t *testing.T) (*OrderHandler, *gin.Engine, *sql.DB) {
	database, err := db.Open(":memory:")
	require.NoError(t, err)
	b0, _ := os.ReadFile("../migrations/001_create_patients.sql")
	require.NoError(t, db.ExecMigration(database, string(b0)))
	b1, _ := os.ReadFile("../migrations/002_create_schedules.sql")
	require.NoError(t, db.ExecMigration(database, string(b1)))
	b2, _ := os.ReadFile("../migrations/003_create_orders.sql")
	require.NoError(t, db.ExecMigration(database, string(b2)))
	b3, _ := os.ReadFile("../migrations/004_alter_orders.sql")
	require.NoError(t, db.ExecMigration(database, string(b3)))

	patientRepo := repo.NewPatientRepository(database)
	scheduleRepo := repo.NewScheduleRepository(database)
	orderRepo := repo.NewOrderRepository(database)
	orderSvc := service.NewOrderService(database, orderRepo, scheduleRepo, patientRepo)
	h := NewOrderHandler(orderSvc)

	gin.SetMode(gin.TestMode)
	g := gin.New()
	g.Use(func(c *gin.Context) {
		if vp := c.GetHeader("X-Visitor-Phone"); vp != "" {
			c.Set("visitor_phone", vp)
		}
		if adminID := c.GetHeader("X-Admin-ID"); adminID != "" {
			c.Set("admin_id", adminID)
		}
		c.Next()
	})

	// H5 routes
	g.GET("/api/v1/orders", h.List)
	g.GET("/api/v1/orders/:id", h.Get)
	g.PUT("/api/v1/orders/:id/cancel", h.Cancel)
	g.PUT("/api/v1/orders/:id/change", h.Change)

	// Admin routes
	g.GET("/api/v1/admin/orders", h.List)
	g.GET("/api/v1/admin/orders/:id", h.Get)
	g.PUT("/api/v1/admin/orders/:id/cancel", h.Cancel)
	g.PUT("/api/v1/admin/orders/:id/change", h.Change)

	return h, g, database
}

func seedOrderData(t *testing.T, db *sql.DB) (orderID int64) {
	patientRepo := repo.NewPatientRepository(db)
	scheduleRepo := repo.NewScheduleRepository(db)
	orderRepo := repo.NewOrderRepository(db)

	pid, _ := patientRepo.Create(&models.Patient{
		Name: "张三", IDCard: "110101********1234", Phone: "138****8888",
		IDCardEncrypted: "enc", PhoneEncrypted: "enc", Gender: "male", Age: 32,
		VisitorPhone: "13800138000",
	})

	sid, _ := scheduleRepo.Create(&models.Schedule{
		Date: "2099-04-29", Department: "内科", DoctorName: "王医生",
		StartTime: "09:00", EndTime: "10:00", TotalQuota: 20, Remaining: 19, Status: "available",
	})

	oid, _ := orderRepo.Create(&models.Order{
		OrderNo: "GH20990429001", ScheduleID: sid, PatientID: pid,
		VisitorPhone: "13800138000", Status: "confirmed",
	})
	return oid
}

func TestOrderHandler_List_H5(t *testing.T) {
	_, r, db := setupOrderHandler(t)
	seedOrderData(t, db)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/orders", nil)
	req.Header.Set("X-Visitor-Phone", "13800138000")
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	var resp Response
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp.Data.(map[string]interface{})
	assert.Equal(t, float64(1), data["total"])
}

func TestOrderHandler_List_Admin(t *testing.T) {
	_, r, db := setupOrderHandler(t)
	seedOrderData(t, db)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/admin/orders?keyword=张三", nil)
	req.Header.Set("X-Admin-ID", "admin1")
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	var resp Response
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp.Data.(map[string]interface{})
	assert.Equal(t, float64(1), data["total"])
}

func TestOrderHandler_Get_H5(t *testing.T) {
	_, r, db := setupOrderHandler(t)
	oid := seedOrderData(t, db)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/orders/"+jsonNum(oid), nil)
	req.Header.Set("X-Visitor-Phone", "13800138000")
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	var resp Response
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp.Data.(map[string]interface{})
	assert.Equal(t, "GH20990429001", data["order_no"])
}

func TestOrderHandler_Get_Admin(t *testing.T) {
	_, r, db := setupOrderHandler(t)
	oid := seedOrderData(t, db)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/admin/orders/"+jsonNum(oid), nil)
	req.Header.Set("X-Admin-ID", "admin1")
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestOrderHandler_Get_Forbidden(t *testing.T) {
	_, r, db := setupOrderHandler(t)
	oid := seedOrderData(t, db)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/orders/"+jsonNum(oid), nil)
	req.Header.Set("X-Visitor-Phone", "13900139000")
	r.ServeHTTP(w, req)
	assert.Equal(t, 403, w.Code)
}

func TestOrderHandler_Cancel_H5(t *testing.T) {
	_, r, db := setupOrderHandler(t)
	oid := seedOrderData(t, db)

	body, _ := json.Marshal(map[string]interface{}{"reason": "个人原因"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/orders/"+jsonNum(oid)+"/cancel", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Visitor-Phone", "13800138000")
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	var resp Response
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp.Data.(map[string]interface{})
	assert.Equal(t, "cancelled", data["status"])
}

func TestOrderHandler_Cancel_Admin(t *testing.T) {
	_, r, db := setupOrderHandler(t)
	oid := seedOrderData(t, db)

	body, _ := json.Marshal(map[string]interface{}{"reason": "医生停诊"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/admin/orders/"+jsonNum(oid)+"/cancel", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Admin-ID", "admin1")
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestOrderHandler_Cancel_TimeLimit(t *testing.T) {
	_, r, db := setupOrderHandler(t)
	patientRepo := repo.NewPatientRepository(db)
	scheduleRepo := repo.NewScheduleRepository(db)
	orderRepo := repo.NewOrderRepository(db)

	pid, _ := patientRepo.Create(&models.Patient{
		Name: "李四", IDCard: "110101********5678", Phone: "139****9999",
		IDCardEncrypted: "enc", PhoneEncrypted: "enc", Gender: "female", Age: 28,
		VisitorPhone: "13800138000",
	})
	sid, _ := scheduleRepo.Create(&models.Schedule{
		Date: "2000-04-29", Department: "外科", DoctorName: "李医生",
		StartTime: "09:00", EndTime: "10:00", TotalQuota: 20, Remaining: 20, Status: "available",
	})
	oid, _ := orderRepo.Create(&models.Order{
		OrderNo: "GH20000429001", ScheduleID: sid, PatientID: pid,
		VisitorPhone: "13800138000", Status: "confirmed",
	})

	body, _ := json.Marshal(map[string]interface{}{"reason": "个人原因"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/orders/"+jsonNum(oid)+"/cancel", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Visitor-Phone", "13800138000")
	r.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestOrderHandler_Change_H5(t *testing.T) {
	_, r, db := setupOrderHandler(t)
	oid := seedOrderData(t, db)
	scheduleRepo := repo.NewScheduleRepository(db)
	nsid, _ := scheduleRepo.Create(&models.Schedule{
		Date: "2099-04-29", Department: "外科", DoctorName: "李医生",
		StartTime: "14:00", EndTime: "15:00", TotalQuota: 10, Remaining: 10, Status: "available",
	})

	body, _ := json.Marshal(map[string]interface{}{"new_schedule_id": nsid})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/orders/"+jsonNum(oid)+"/change", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Visitor-Phone", "13800138000")
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	var resp Response
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp.Data.(map[string]interface{})
	assert.Equal(t, "confirmed", data["status"])
}

func TestOrderHandler_Change_Admin(t *testing.T) {
	_, r, db := setupOrderHandler(t)
	oid := seedOrderData(t, db)
	scheduleRepo := repo.NewScheduleRepository(db)
	nsid, _ := scheduleRepo.Create(&models.Schedule{
		Date: "2099-04-29", Department: "外科", DoctorName: "李医生",
		StartTime: "14:00", EndTime: "15:00", TotalQuota: 10, Remaining: 10, Status: "available",
	})

	body, _ := json.Marshal(map[string]interface{}{"new_schedule_id": nsid})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/admin/orders/"+jsonNum(oid)+"/change", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Admin-ID", "admin1")
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestOrderHandler_Change_NewScheduleFull(t *testing.T) {
	_, r, db := setupOrderHandler(t)
	oid := seedOrderData(t, db)
	scheduleRepo := repo.NewScheduleRepository(db)
	nsid, _ := scheduleRepo.Create(&models.Schedule{
		Date: "2099-04-29", Department: "外科", DoctorName: "李医生",
		StartTime: "14:00", EndTime: "15:00", TotalQuota: 10, Remaining: 0, Status: "full",
	})

	body, _ := json.Marshal(map[string]interface{}{"new_schedule_id": nsid})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/orders/"+jsonNum(oid)+"/change", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Visitor-Phone", "13800138000")
	r.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}
