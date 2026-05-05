package service

import (
	"os"
	"testing"

	"clinic/db"
	"clinic/models"
	"clinic/repo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupScheduleSvc(t *testing.T) (*ScheduleService, func()) {
	database, err := db.Open(":memory:")
	require.NoError(t, err)
	b, _ := os.ReadFile("../migrations/002_create_schedules.sql")
	require.NoError(t, db.ExecMigration(database, string(b)))
	r := repo.NewScheduleRepository(database)
	return NewScheduleService(r), func() { database.Close() }
}

func TestScheduleService_Create(t *testing.T) {
	svc, cleanup := setupScheduleSvc(t)
	defer cleanup()

	s, err := svc.CreateSchedule(&models.Schedule{
		Date:       "2026-04-29",
		Department: "内科",
		DoctorName: "王医生",
		StartTime:  "09:00",
		EndTime:    "10:00",
		TotalQuota: 20,
	})
	require.NoError(t, err)
	assert.Greater(t, s.ID, int64(0))
	assert.Equal(t, 20, s.Remaining)
	assert.Equal(t, "available", s.Status)
}

func TestScheduleService_List(t *testing.T) {
	svc, cleanup := setupScheduleSvc(t)
	defer cleanup()

	_, _ = svc.CreateSchedule(&models.Schedule{Date: "2026-04-29", Department: "内科", DoctorName: "A", StartTime: "08:00", EndTime: "09:00", TotalQuota: 10})
	_, _ = svc.CreateSchedule(&models.Schedule{Date: "2026-04-29", Department: "外科", DoctorName: "B", StartTime: "09:00", EndTime: "10:00", TotalQuota: 10})

	list, total, err := svc.ListSchedules("2026-04-29", "", "", 1, 10)
	require.NoError(t, err)
	assert.Equal(t, 2, total)
	assert.Len(t, list, 2)

	list, total, err = svc.ListSchedules("2026-04-29", "外科", "", 1, 10)
	require.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Equal(t, "外科", list[0].Department)

	list, total, err = svc.ListSchedules("2026-04-29", "", "A", 1, 10)
	require.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Equal(t, "A", list[0].DoctorName)
}

func TestScheduleService_Get(t *testing.T) {
	svc, cleanup := setupScheduleSvc(t)
	defer cleanup()

	created, _ := svc.CreateSchedule(&models.Schedule{Date: "2026-04-29", Department: "内科", DoctorName: "王医生", StartTime: "09:00", EndTime: "10:00", TotalQuota: 20})
	got, err := svc.GetSchedule(created.ID)
	require.NoError(t, err)
	assert.Equal(t, "王医生", got.DoctorName)

	_, err = svc.GetSchedule(9999)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "不存在")
}

func TestScheduleService_Update(t *testing.T) {
	svc, cleanup := setupScheduleSvc(t)
	defer cleanup()

	created, _ := svc.CreateSchedule(&models.Schedule{Date: "2026-04-29", Department: "内科", DoctorName: "王医生", StartTime: "09:00", EndTime: "10:00", TotalQuota: 20})

	// deduct 2
	_, _ = svc.Deduct(created.ID)
	_, _ = svc.Deduct(created.ID)

	// update total_quota to 18 should fail (booked = 2)
	err := svc.UpdateSchedule(&models.Schedule{ID: created.ID, TotalQuota: 1})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "已预约数")

	// update total_quota to 19 should succeed
	err = svc.UpdateSchedule(&models.Schedule{ID: created.ID, TotalQuota: 19})
	require.NoError(t, err)
	got, _ := svc.GetSchedule(created.ID)
	assert.Equal(t, 19, got.TotalQuota)
}

func TestScheduleService_Delete(t *testing.T) {
	svc, cleanup := setupScheduleSvc(t)
	defer cleanup()

	created, _ := svc.CreateSchedule(&models.Schedule{Date: "2026-04-29", Department: "内科", DoctorName: "王医生", StartTime: "09:00", EndTime: "10:00", TotalQuota: 20})
	_, _ = svc.Deduct(created.ID)

	err := svc.DeleteSchedule(created.ID)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "已有预约")

	// rollback then delete
	_, _ = svc.Rollback(created.ID)
	err = svc.DeleteSchedule(created.ID)
	require.NoError(t, err)
}

func TestScheduleService_Deduct(t *testing.T) {
	svc, cleanup := setupScheduleSvc(t)
	defer cleanup()

	created, _ := svc.CreateSchedule(&models.Schedule{Date: "2026-04-29", Department: "内科", DoctorName: "王医生", StartTime: "09:00", EndTime: "10:00", TotalQuota: 2})

	ok, err := svc.Deduct(created.ID)
	require.NoError(t, err)
	assert.True(t, ok)

	got, _ := svc.GetSchedule(created.ID)
	assert.Equal(t, 1, got.Remaining)
	assert.Equal(t, "available", got.Status)

	ok, err = svc.Deduct(created.ID)
	require.NoError(t, err)
	assert.True(t, ok)

	got, _ = svc.GetSchedule(created.ID)
	assert.Equal(t, 0, got.Remaining)
	assert.Equal(t, "full", got.Status)

	ok, err = svc.Deduct(created.ID)
	require.NoError(t, err)
	assert.False(t, ok)
}

func TestScheduleService_Rollback(t *testing.T) {
	svc, cleanup := setupScheduleSvc(t)
	defer cleanup()

	created, _ := svc.CreateSchedule(&models.Schedule{Date: "2026-04-29", Department: "内科", DoctorName: "王医生", StartTime: "09:00", EndTime: "10:00", TotalQuota: 2})
	_, _ = svc.Deduct(created.ID)

	ok, err := svc.Rollback(created.ID)
	require.NoError(t, err)
	assert.True(t, ok)

	got, _ := svc.GetSchedule(created.ID)
	assert.Equal(t, 2, got.Remaining)
	assert.Equal(t, "available", got.Status)

	// rollback when full should fail
	ok, err = svc.Rollback(created.ID)
	require.NoError(t, err)
	assert.False(t, ok)
}
