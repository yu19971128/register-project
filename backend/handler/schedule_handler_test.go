package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"clinic/db"
	"clinic/repo"
	"clinic/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupScheduleHandler(t *testing.T) (*ScheduleHandler, *gin.Engine) {
	database, err := db.Open(":memory:")
	require.NoError(t, err)
	b, _ := os.ReadFile("../migrations/002_create_schedules.sql")
	require.NoError(t, db.ExecMigration(database, string(b)))
	r := repo.NewScheduleRepository(database)
	svc := service.NewScheduleService(r)
	h := NewScheduleHandler(svc)

	gin.SetMode(gin.TestMode)
	g := gin.New()
	g.POST("/api/v1/schedules", h.Create)
	g.GET("/api/v1/schedules", h.List)
	g.GET("/api/v1/schedules/:id", h.Get)
	g.PUT("/api/v1/schedules/:id", h.Update)
	g.DELETE("/api/v1/schedules/:id", h.Delete)
	g.POST("/api/v1/schedules/:id/deduct", h.Deduct)
	g.POST("/api/v1/schedules/:id/rollback", h.Rollback)
	return h, g
}

func TestScheduleHandler_Create_List_Get_Update_Delete(t *testing.T) {
	_, r := setupScheduleHandler(t)

	// Create
	body, _ := json.Marshal(map[string]interface{}{
		"date": "2026-04-29", "department": "内科", "doctor_name": "王医生",
		"start_time": "09:00", "end_time": "10:00", "total_quota": 20,
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/schedules", bytes.NewReader(body))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	var createResp Response
	_ = json.Unmarshal(w.Body.Bytes(), &createResp)
	id := int64(createResp.Data.(map[string]interface{})["id"].(float64))
	assert.Greater(t, id, int64(0))

	// List
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/schedules?date=2026-04-29", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	var listResp Response
	_ = json.Unmarshal(w.Body.Bytes(), &listResp)
	listData := listResp.Data.(map[string]interface{})
	assert.Equal(t, float64(1), listData["total"])

	// Get
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/schedules/"+jsonNum(id), nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	// Update total_quota
	body, _ = json.Marshal(map[string]interface{}{"total_quota": 25})
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/api/v1/schedules/"+jsonNum(id), bytes.NewReader(body))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	// Delete
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/v1/schedules/"+jsonNum(id), nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	// Get after delete
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/schedules/"+jsonNum(id), nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 404, w.Code)
}

func TestScheduleHandler_Create_Conflict_SameDoctorTime(t *testing.T) {
	_, r := setupScheduleHandler(t)
	body, _ := json.Marshal(map[string]interface{}{
		"date": "2026-04-29", "department": "内科", "doctor_name": "王医生",
		"start_time": "09:00", "end_time": "10:00", "total_quota": 20,
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/schedules", bytes.NewReader(body))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/schedules", bytes.NewReader(body))
	r.ServeHTTP(w, req)
	assert.Equal(t, 409, w.Code)
}

func TestScheduleHandler_Deduct_And_Rollback(t *testing.T) {
	_, r := setupScheduleHandler(t)
	body, _ := json.Marshal(map[string]interface{}{
		"date": "2026-04-29", "department": "内科", "doctor_name": "王医生",
		"start_time": "09:00", "end_time": "10:00", "total_quota": 1,
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/schedules", bytes.NewReader(body))
	r.ServeHTTP(w, req)
	var createResp Response
	_ = json.Unmarshal(w.Body.Bytes(), &createResp)
	id := int64(createResp.Data.(map[string]interface{})["id"].(float64))

	// Deduct
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/schedules/"+jsonNum(id)+"/deduct", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	// Deduct again should fail
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/schedules/"+jsonNum(id)+"/deduct", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 409, w.Code)

	// Rollback
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/schedules/"+jsonNum(id)+"/rollback", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestScheduleHandler_Delete_With_Booking(t *testing.T) {
	_, r := setupScheduleHandler(t)
	body, _ := json.Marshal(map[string]interface{}{
		"date": "2026-04-29", "department": "内科", "doctor_name": "王医生",
		"start_time": "09:00", "end_time": "10:00", "total_quota": 2,
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/schedules", bytes.NewReader(body))
	r.ServeHTTP(w, req)
	var createResp Response
	_ = json.Unmarshal(w.Body.Bytes(), &createResp)
	id := int64(createResp.Data.(map[string]interface{})["id"].(float64))

	// Deduct to create booking
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/schedules/"+jsonNum(id)+"/deduct", nil)
	r.ServeHTTP(w, req)

	// Delete should fail
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/v1/schedules/"+jsonNum(id), nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestScheduleHandler_Update_LowerThanBooked(t *testing.T) {
	_, r := setupScheduleHandler(t)
	body, _ := json.Marshal(map[string]interface{}{
		"date": "2026-04-29", "department": "内科", "doctor_name": "王医生",
		"start_time": "09:00", "end_time": "10:00", "total_quota": 2,
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/schedules", bytes.NewReader(body))
	r.ServeHTTP(w, req)
	var createResp Response
	_ = json.Unmarshal(w.Body.Bytes(), &createResp)
	id := int64(createResp.Data.(map[string]interface{})["id"].(float64))

	// Deduct twice so booked = 2
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/schedules/"+jsonNum(id)+"/deduct", nil)
	r.ServeHTTP(w, req)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/schedules/"+jsonNum(id)+"/deduct", nil)
	r.ServeHTTP(w, req)

	// Update to 1 should fail (booked = 2)
	body, _ = json.Marshal(map[string]interface{}{"total_quota": 1})
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/api/v1/schedules/"+jsonNum(id), bytes.NewReader(body))
	r.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestScheduleHandler_Update_RecomputeRemaining(t *testing.T) {
	_, r := setupScheduleHandler(t)

	// Create with total=10, no bookings yet → remaining=10
	body, _ := json.Marshal(map[string]interface{}{
		"date": "2026-04-29", "department": "内科", "doctor_name": "王医生",
		"start_time": "09:00", "end_time": "10:00", "total_quota": 10,
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/schedules", bytes.NewReader(body))
	r.ServeHTTP(w, req)
	var createResp Response
	_ = json.Unmarshal(w.Body.Bytes(), &createResp)
	id := int64(createResp.Data.(map[string]interface{})["id"].(float64))

	// Reduce total_quota to 2 (no bookings) → remaining must be 2, not 0
	body, _ = json.Marshal(map[string]interface{}{"total_quota": 2})
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/api/v1/schedules/"+jsonNum(id), bytes.NewReader(body))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/schedules/"+jsonNum(id), nil)
	r.ServeHTTP(w, req)
	var getResp Response
	_ = json.Unmarshal(w.Body.Bytes(), &getResp)
	data := getResp.Data.(map[string]interface{})
	assert.Equal(t, float64(2), data["total_quota"])
	assert.Equal(t, float64(2), data["remaining"])

	// Increase total_quota back to 5 → remaining must follow (5 - 0 booked = 5)
	body, _ = json.Marshal(map[string]interface{}{"total_quota": 5})
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/api/v1/schedules/"+jsonNum(id), bytes.NewReader(body))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/schedules/"+jsonNum(id), nil)
	r.ServeHTTP(w, req)
	_ = json.Unmarshal(w.Body.Bytes(), &getResp)
	data = getResp.Data.(map[string]interface{})
	assert.Equal(t, float64(5), data["remaining"])

	// Deduct 2 (booked=2, remaining=3), then reduce total to 4 → remaining = 4-2 = 2
	for i := 0; i < 2; i++ {
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("POST", "/api/v1/schedules/"+jsonNum(id)+"/deduct", nil)
		r.ServeHTTP(w, req)
	}
	body, _ = json.Marshal(map[string]interface{}{"total_quota": 4})
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/api/v1/schedules/"+jsonNum(id), bytes.NewReader(body))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/schedules/"+jsonNum(id), nil)
	r.ServeHTTP(w, req)
	_ = json.Unmarshal(w.Body.Bytes(), &getResp)
	data = getResp.Data.(map[string]interface{})
	assert.Equal(t, float64(4), data["total_quota"])
	assert.Equal(t, float64(2), data["remaining"])
}

func TestScheduleHandler_List_DoctorFilter(t *testing.T) {
	_, r := setupScheduleHandler(t)

	// Create schedules for two doctors on same date
	for _, doctor := range []string{"王医生", "李医生"} {
		body, _ := json.Marshal(map[string]interface{}{
			"date": "2026-04-29", "department": "内科", "doctor_name": doctor,
			"start_time": "09:00", "end_time": "10:00", "total_quota": 10,
		})
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/schedules", bytes.NewReader(body))
		r.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
	}

	// List without doctor filter → total = 2
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/schedules?date=2026-04-29&page=1&page_size=10", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	var resp Response
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp.Data.(map[string]interface{})
	assert.Equal(t, float64(2), data["total"])

	// List with doctor_name=王医生 → total = 1
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/schedules?date=2026-04-29&doctor_name=王医生&page=1&page_size=10", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	data = resp.Data.(map[string]interface{})
	assert.Equal(t, float64(1), data["total"])
	list := data["list"].([]interface{})
	assert.Len(t, list, 1)
	first := list[0].(map[string]interface{})
	assert.Equal(t, "王医生", first["doctor_name"])
}

func jsonNum(n int64) string {
	b, _ := json.Marshal(n)
	return string(b)
}
