package repo

import (
	"os"
	"testing"

	"clinic/db"
	"clinic/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupScheduleRepo(t *testing.T) (*ScheduleRepository, func()) {
	database, err := db.Open(":memory:")
	require.NoError(t, err)
	b, _ := os.ReadFile("../migrations/002_create_schedules.sql")
	require.NoError(t, db.ExecMigration(database, string(b)))
	return NewScheduleRepository(database), func() { database.Close() }
}

func TestScheduleRepository_Create_Get_List_Update_Delete(t *testing.T) {
	r, cleanup := setupScheduleRepo(t)
	defer cleanup()

	// Create
	s := &models.Schedule{
		Date:       "2026-04-29",
		Department: "内科",
		DoctorName: "王医生",
		StartTime:  "09:00",
		EndTime:    "10:00",
		TotalQuota: 20,
		Remaining:  20,
		Status:     "available",
	}
	id, err := r.Create(s)
	require.NoError(t, err)
	assert.Greater(t, id, int64(0))

	// Get
	got, err := r.GetByID(id)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, "内科", got.Department)
	assert.Equal(t, "王医生", got.DoctorName)
	assert.Equal(t, 20, got.Remaining)

	// List
	list, total, err := r.List("2026-04-29", "", 0, 10)
	require.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Len(t, list, 1)

	// List with department filter
	list, total, err = r.List("2026-04-29", "外科", 0, 10)
	require.NoError(t, err)
	assert.Equal(t, 0, total)

	// Update
	s.ID = id
	s.TotalQuota = 25
	err = r.Update(s)
	require.NoError(t, err)
	got, err = r.GetByID(id)
	require.NoError(t, err)
	assert.Equal(t, 25, got.TotalQuota)

	// Delete
	err = r.Delete(id)
	require.NoError(t, err)
	got, err = r.GetByID(id)
	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestScheduleRepository_Deduct_Success(t *testing.T) {
	r, cleanup := setupScheduleRepo(t)
	defer cleanup()

	s := &models.Schedule{Date: "2026-04-29", Department: "内科", DoctorName: "王医生", StartTime: "09:00", EndTime: "10:00", TotalQuota: 2, Remaining: 2, Status: "available"}
	id, _ := r.Create(s)

	ok, err := r.Deduct(id)
	require.NoError(t, err)
	assert.True(t, ok)

	got, _ := r.GetByID(id)
	assert.Equal(t, 1, got.Remaining)
}

func TestScheduleRepository_Deduct_Fail_When_Zero(t *testing.T) {
	r, cleanup := setupScheduleRepo(t)
	defer cleanup()

	s := &models.Schedule{Date: "2026-04-29", Department: "内科", DoctorName: "王医生", StartTime: "09:00", EndTime: "10:00", TotalQuota: 1, Remaining: 0, Status: "full"}
	id, _ := r.Create(s)

	ok, err := r.Deduct(id)
	require.NoError(t, err)
	assert.False(t, ok)
}

func TestScheduleRepository_Rollback_Success(t *testing.T) {
	r, cleanup := setupScheduleRepo(t)
	defer cleanup()

	s := &models.Schedule{Date: "2026-04-29", Department: "内科", DoctorName: "王医生", StartTime: "09:00", EndTime: "10:00", TotalQuota: 2, Remaining: 1, Status: "available"}
	id, _ := r.Create(s)

	ok, err := r.Rollback(id)
	require.NoError(t, err)
	assert.True(t, ok)

	got, _ := r.GetByID(id)
	assert.Equal(t, 2, got.Remaining)
}

func TestScheduleRepository_Rollback_Capped_By_TotalQuota(t *testing.T) {
	r, cleanup := setupScheduleRepo(t)
	defer cleanup()

	s := &models.Schedule{Date: "2026-04-29", Department: "内科", DoctorName: "王医生", StartTime: "09:00", EndTime: "10:00", TotalQuota: 2, Remaining: 2, Status: "full"}
	id, _ := r.Create(s)

	ok, err := r.Rollback(id)
	require.NoError(t, err)
	assert.False(t, ok)

	got, _ := r.GetByID(id)
	assert.Equal(t, 2, got.Remaining)
}
