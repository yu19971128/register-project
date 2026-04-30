package service

import (
	"os"
	"testing"
	"time"

	"clinic/db"
	"clinic/models"
	"clinic/repo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupRegistrationSvc(t *testing.T) (*RegistrationService, func()) {
	database, err := db.Open(":memory:")
	require.NoError(t, err)
	b1, _ := os.ReadFile("../migrations/001_create_patients.sql")
	require.NoError(t, db.ExecMigration(database, string(b1)))
	b2, _ := os.ReadFile("../migrations/002_create_schedules.sql")
	require.NoError(t, db.ExecMigration(database, string(b2)))
	b3, _ := os.ReadFile("../migrations/003_create_orders.sql")
	require.NoError(t, db.ExecMigration(database, string(b3)))
	b4, _ := os.ReadFile("../migrations/004_alter_orders.sql")
	require.NoError(t, db.ExecMigration(database, string(b4)))

	patientRepo := repo.NewPatientRepository(database)
	scheduleRepo := repo.NewScheduleRepository(database)
	orderRepo := repo.NewOrderRepository(database)

	return NewRegistrationService(database, patientRepo, scheduleRepo, orderRepo), func() { database.Close() }
}

func seedPatientAndSchedule(t *testing.T, r *RegistrationService) (patientID, scheduleID int64) {
	pRepo := r.patientRepo
	sRepo := r.scheduleRepo

	pid, _ := pRepo.Create(&models.Patient{
		Name: "张三", IDCard: "110101********1234", Phone: "138****8888",
		IDCardEncrypted: "enc", PhoneEncrypted: "enc", Gender: "male", Age: 32,
		VisitorPhone: "13800138000",
	})

	sid, _ := sRepo.Create(&models.Schedule{
		Date: "2026-04-29", Department: "内科", DoctorName: "王医生",
		StartTime: "09:00", EndTime: "10:00", TotalQuota: 20, Remaining: 20, Status: "available",
	})
	return pid, sid
}

func TestRegistrationService_Submit_Success(t *testing.T) {
	svc, cleanup := setupRegistrationSvc(t)
	defer cleanup()
	pid, sid := seedPatientAndSchedule(t, svc)

	result, err := svc.SubmitRegistration(sid, pid, "13800138000")
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "confirmed", result.Status)
	assert.Contains(t, result.OrderNo, "GH20260429")

	// verify schedule deducted
	s, _ := svc.scheduleRepo.GetByID(sid)
	assert.Equal(t, 19, s.Remaining)
}

func TestRegistrationService_Submit_PatientNotFound(t *testing.T) {
	svc, cleanup := setupRegistrationSvc(t)
	defer cleanup()
	_, sid := seedPatientAndSchedule(t, svc)

	_, err := svc.SubmitRegistration(sid, 9999, "13800138000")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "就诊人不存在")
}

func TestRegistrationService_Submit_PatientNotBelong(t *testing.T) {
	svc, cleanup := setupRegistrationSvc(t)
	defer cleanup()
	pid, sid := seedPatientAndSchedule(t, svc)

	_, err := svc.SubmitRegistration(sid, pid, "13900139000")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "不属于当前访客")
}

func TestRegistrationService_Submit_ScheduleNotFound(t *testing.T) {
	svc, cleanup := setupRegistrationSvc(t)
	defer cleanup()
	pid, _ := seedPatientAndSchedule(t, svc)

	_, err := svc.SubmitRegistration(9999, pid, "13800138000")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "号源不存在")
}

func TestRegistrationService_Submit_Duplicate(t *testing.T) {
	svc, cleanup := setupRegistrationSvc(t)
	defer cleanup()
	pid, sid := seedPatientAndSchedule(t, svc)

	_, err := svc.SubmitRegistration(sid, pid, "13800138000")
	require.NoError(t, err)

	_, err = svc.SubmitRegistration(sid, pid, "13800138000")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "重复提交")
}

func TestRegistrationService_GetTicket_Success(t *testing.T) {
	svc, cleanup := setupRegistrationSvc(t)
	defer cleanup()
	pid, sid := seedPatientAndSchedule(t, svc)

	result, _ := svc.SubmitRegistration(sid, pid, "13800138000")
	ticket, err := svc.GetTicket(result.OrderNo, "13800138000")
	require.NoError(t, err)
	require.NotNil(t, ticket)
	assert.Equal(t, result.OrderNo, ticket.OrderNo)
	assert.Equal(t, "张三", ticket.PatientName)
}

func TestRegistrationService_GetTicket_Forbidden(t *testing.T) {
	svc, cleanup := setupRegistrationSvc(t)
	defer cleanup()
	pid, sid := seedPatientAndSchedule(t, svc)

	result, _ := svc.SubmitRegistration(sid, pid, "13800138000")
	_, err := svc.GetTicket(result.OrderNo, "13900139000")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "无权查看")
}

func TestRegistrationService_GetTicket_NotFound(t *testing.T) {
	svc, cleanup := setupRegistrationSvc(t)
	defer cleanup()

	_, err := svc.GetTicket("GH20260429099", "13800138000")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "不存在")
}

func TestRegistrationService_GenerateOrderNo(t *testing.T) {
	no1 := generateOrderNo("2026-04-29")
	time.Sleep(time.Millisecond)
	no2 := generateOrderNo("2026-04-29")
	assert.NotEqual(t, no1, no2)
	assert.Contains(t, no1, "20260429")
}
