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

func setupOrderSvc(t *testing.T) (*OrderService, func()) {
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

	return NewOrderService(database, orderRepo, scheduleRepo, patientRepo), func() { database.Close() }
}

func seedPatientScheduleOrder(t *testing.T, svc *OrderService) (patientID, scheduleID, orderID int64) {
	pRepo := svc.patientRepo
	sRepo := svc.scheduleRepo
	oRepo := svc.orderRepo

	pid, _ := pRepo.Create(&models.Patient{
		Name: "张三", IDCard: "110101********1234", Phone: "138****8888",
		IDCardEncrypted: "enc", PhoneEncrypted: "enc", Gender: "male", Age: 32,
		VisitorPhone: "13800138000",
	})

	sid, _ := sRepo.Create(&models.Schedule{
		Date: "2099-04-29", Department: "内科", DoctorName: "王医生",
		StartTime: "09:00", EndTime: "10:00", TotalQuota: 20, Remaining: 19, Status: "available",
	})

	oid, _ := oRepo.Create(&models.Order{
		OrderNo: "GH20990429001", ScheduleID: sid, PatientID: pid,
		VisitorPhone: "13800138000", Status: "confirmed",
	})
	return pid, sid, oid
}

func seedPastOrder(t *testing.T, svc *OrderService) (patientID, scheduleID, orderID int64) {
	pRepo := svc.patientRepo
	sRepo := svc.scheduleRepo
	oRepo := svc.orderRepo

	pid, _ := pRepo.Create(&models.Patient{
		Name: "李四", IDCard: "110101********5678", Phone: "139****9999",
		IDCardEncrypted: "enc", PhoneEncrypted: "enc", Gender: "female", Age: 28,
		VisitorPhone: "13800138000",
	})

	sid, _ := sRepo.Create(&models.Schedule{
		Date: "2000-04-29", Department: "外科", DoctorName: "李医生",
		StartTime: "09:00", EndTime: "10:00", TotalQuota: 20, Remaining: 20, Status: "available",
	})

	oid, _ := oRepo.Create(&models.Order{
		OrderNo: "GH20000429001", ScheduleID: sid, PatientID: pid,
		VisitorPhone: "13800138000", Status: "confirmed",
	})
	return pid, sid, oid
}

func TestOrderService_List_Admin(t *testing.T) {
	svc, cleanup := setupOrderSvc(t)
	defer cleanup()
	seedPatientScheduleOrder(t, svc)

	res, err := svc.List("", "", "", "", "", 1, 10, true, "")
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, 1, res.Total)
	assert.Len(t, res.List, 1)
	assert.Equal(t, "张三", res.List[0].PatientName)
	assert.Equal(t, "内科", res.List[0].Department)
}

func TestOrderService_List_Patient(t *testing.T) {
	svc, cleanup := setupOrderSvc(t)
	defer cleanup()
	seedPatientScheduleOrder(t, svc)

	// another patient's order
	pid2, sid2, _ := seedPatientScheduleOrder(t, svc)
	_ = pid2
	_ = sid2
	// Actually seedPatientScheduleOrder always uses visitor_phone 13800138000, so both belong to same visitor.
	// Create a different visitor order manually
	oRepo := svc.orderRepo
	oRepo.Create(&models.Order{OrderNo: "GH20990429002", ScheduleID: sid2, PatientID: pid2, VisitorPhone: "13900139000", Status: "confirmed"})

	res, err := svc.List("", "", "", "", "", 1, 10, false, "13800138000")
	require.NoError(t, err)
	assert.Equal(t, 1, res.Total) // only the seeded one (GH20990429001) because the second seed reuses same phone
}

func TestOrderService_GetDetail_Success(t *testing.T) {
	svc, cleanup := setupOrderSvc(t)
	defer cleanup()
	_, _, oid := seedPatientScheduleOrder(t, svc)

	res, err := svc.GetDetail(oid, false, "13800138000")
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, "GH20990429001", res.OrderNo)
	assert.Equal(t, "张三", res.Patient.Name)
	assert.Equal(t, "内科", res.Schedule.Department)
}

func TestOrderService_GetDetail_Forbidden(t *testing.T) {
	svc, cleanup := setupOrderSvc(t)
	defer cleanup()
	_, _, oid := seedPatientScheduleOrder(t, svc)

	_, err := svc.GetDetail(oid, false, "13900139000")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "无权")
}

func TestOrderService_Cancel_Success(t *testing.T) {
	svc, cleanup := setupOrderSvc(t)
	defer cleanup()
	_, sid, oid := seedPatientScheduleOrder(t, svc)

	res, err := svc.Cancel(oid, "个人原因", false, "13800138000", "")
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, "cancelled", res.Status)

	// schedule rolled back
	s, _ := svc.scheduleRepo.GetByID(sid)
	assert.Equal(t, 20, s.Remaining)
}

func TestOrderService_Cancel_TimeLimit(t *testing.T) {
	svc, cleanup := setupOrderSvc(t)
	defer cleanup()
	_, _, oid := seedPastOrder(t, svc)

	_, err := svc.Cancel(oid, "个人原因", false, "13800138000", "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "退号时限")
}

func TestOrderService_Cancel_Forbidden(t *testing.T) {
	svc, cleanup := setupOrderSvc(t)
	defer cleanup()
	_, _, oid := seedPatientScheduleOrder(t, svc)

	_, err := svc.Cancel(oid, "个人原因", false, "13900139000", "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "无权")
}

func TestOrderService_Change_Success(t *testing.T) {
	svc, cleanup := setupOrderSvc(t)
	defer cleanup()
	_, sid, oid := seedPatientScheduleOrder(t, svc)

	// create new schedule
	nsid, _ := svc.scheduleRepo.Create(&models.Schedule{
		Date: "2099-04-29", Department: "外科", DoctorName: "李医生",
		StartTime: "14:00", EndTime: "15:00", TotalQuota: 10, Remaining: 10, Status: "available",
	})

	res, err := svc.Change(oid, nsid, false, "13800138000")
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, "confirmed", res.Status)
	assert.Equal(t, "李医生", res.Schedule.DoctorName)

	// old schedule rolled back
	oldS, _ := svc.scheduleRepo.GetByID(sid)
	assert.Equal(t, 20, oldS.Remaining)

	// new schedule deducted
	newS, _ := svc.scheduleRepo.GetByID(nsid)
	assert.Equal(t, 9, newS.Remaining)

	// old order cancelled
	oldOrder, _ := svc.orderRepo.GetByID(oid)
	assert.Equal(t, "cancelled", oldOrder.Status)
}

func TestOrderService_Change_NewScheduleFull(t *testing.T) {
	svc, cleanup := setupOrderSvc(t)
	defer cleanup()
	_, _, oid := seedPatientScheduleOrder(t, svc)

	// create full schedule
	nsid, _ := svc.scheduleRepo.Create(&models.Schedule{
		Date: "2099-04-29", Department: "外科", DoctorName: "李医生",
		StartTime: "14:00", EndTime: "15:00", TotalQuota: 10, Remaining: 0, Status: "full",
	})

	_, err := svc.Change(oid, nsid, false, "13800138000")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "余量")

	// old order unchanged
	oldOrder, _ := svc.orderRepo.GetByID(oid)
	assert.Equal(t, "confirmed", oldOrder.Status)
}

func TestOrderService_Change_Forbidden(t *testing.T) {
	svc, cleanup := setupOrderSvc(t)
	defer cleanup()
	_, _, oid := seedPatientScheduleOrder(t, svc)

	nsid, _ := svc.scheduleRepo.Create(&models.Schedule{
		Date: "2099-04-29", Department: "外科", DoctorName: "李医生",
		StartTime: "14:00", EndTime: "15:00", TotalQuota: 10, Remaining: 10, Status: "available",
	})

	_, err := svc.Change(oid, nsid, false, "13900139000")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "无权")
}
