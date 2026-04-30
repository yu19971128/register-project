package repo

import (
	"os"
	"testing"

	"clinic/db"
	"clinic/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupOrderRepo(t *testing.T) (*OrderRepository, func()) {
	database, err := db.Open(":memory:")
	require.NoError(t, err)
	b1, _ := os.ReadFile("../migrations/002_create_schedules.sql")
	require.NoError(t, db.ExecMigration(database, string(b1)))
	b2, _ := os.ReadFile("../migrations/003_create_orders.sql")
	require.NoError(t, db.ExecMigration(database, string(b2)))
	return NewOrderRepository(database), func() { database.Close() }
}

func TestOrderRepository_Create_GetByOrderNo(t *testing.T) {
	r, cleanup := setupOrderRepo(t)
	defer cleanup()

	o := &models.Order{
		OrderNo:      "GH20260429001",
		ScheduleID:   1,
		PatientID:    2,
		VisitorPhone: "13800138000",
		Status:       "confirmed",
	}
	id, err := r.Create(o)
	require.NoError(t, err)
	assert.Greater(t, id, int64(0))

	got, err := r.GetByOrderNo("GH20260429001")
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, int64(1), got.ScheduleID)
	assert.Equal(t, "confirmed", got.Status)
}

func TestOrderRepository_GetByID(t *testing.T) {
	r, cleanup := setupOrderRepo(t)
	defer cleanup()

	o := &models.Order{OrderNo: "GH20260429002", ScheduleID: 1, PatientID: 2, VisitorPhone: "13800138000", Status: "confirmed"}
	id, _ := r.Create(o)

	got, err := r.GetByID(id)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, "GH20260429002", got.OrderNo)

	got, err = r.GetByID(9999)
	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestOrderRepository_ExistsByScheduleAndVisitor(t *testing.T) {
	r, cleanup := setupOrderRepo(t)
	defer cleanup()

	_, _ = r.Create(&models.Order{OrderNo: "GH20260429003", ScheduleID: 1, PatientID: 2, VisitorPhone: "13800138000", Status: "confirmed"})

	exists, err := r.ExistsByScheduleAndVisitor(1, "13800138000")
	require.NoError(t, err)
	assert.True(t, exists)

	exists, err = r.ExistsByScheduleAndVisitor(2, "13800138000")
	require.NoError(t, err)
	assert.False(t, exists)
}
