package repo

import (
	"os"
	"testing"
	"time"

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
	b3, _ := os.ReadFile("../migrations/004_alter_orders.sql")
	require.NoError(t, db.ExecMigration(database, string(b3)))
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

func TestOrderRepository_List(t *testing.T) {
	r, cleanup := setupOrderRepo(t)
	defer cleanup()

	_, _ = r.Create(&models.Order{OrderNo: "GH20260429001", ScheduleID: 1, PatientID: 1, VisitorPhone: "13800138000", Status: "confirmed"})
	_, _ = r.Create(&models.Order{OrderNo: "GH20260429002", ScheduleID: 2, PatientID: 2, VisitorPhone: "13800138000", Status: "cancelled"})
	_, _ = r.Create(&models.Order{OrderNo: "GH20260429003", ScheduleID: 1, PatientID: 3, VisitorPhone: "13900139000", Status: "confirmed"})

	// List all
	list, total, err := r.List("", "", "", "", 0, 10)
	require.NoError(t, err)
	assert.Equal(t, 3, total)
	assert.Len(t, list, 3)

	// Filter by status
	list, total, err = r.List("", "", "", "confirmed", 0, 10)
	require.NoError(t, err)
	assert.Equal(t, 2, total)

	// Filter by visitor phone
	list, total, err = r.List("", "", "", "", 0, 10, "13800138000")
	require.NoError(t, err)
	assert.Equal(t, 2, total)
}

func TestOrderRepository_UpdateStatus(t *testing.T) {
	r, cleanup := setupOrderRepo(t)
	defer cleanup()

	o := &models.Order{OrderNo: "GH20260429004", ScheduleID: 1, PatientID: 1, VisitorPhone: "13800138000", Status: "confirmed"}
	id, _ := r.Create(o)

	cancelledAt := time.Now()
	err := r.UpdateStatus(id, "cancelled", "个人原因", cancelledAt, "admin")
	require.NoError(t, err)

	got, _ := r.GetByID(id)
	assert.Equal(t, "cancelled", got.Status)
	assert.Equal(t, "个人原因", got.CancelReason)
	assert.Equal(t, "admin", got.OperatedBy)
}
