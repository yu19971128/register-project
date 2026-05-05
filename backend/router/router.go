package router

import (
	"database/sql"

	"clinic/handler"
	"clinic/middleware"
	"clinic/repo"
	"clinic/service"

	"github.com/gin-gonic/gin"
)

func Setup(db *sql.DB) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())

	patientRepo := repo.NewPatientRepository(db)
	patientSvc := service.NewPatientService(patientRepo)
	patientHandler := handler.NewPatientHandler(patientSvc)

	scheduleRepo := repo.NewScheduleRepository(db)
	scheduleSvc := service.NewScheduleService(scheduleRepo)
	scheduleHandler := handler.NewScheduleHandler(scheduleSvc)

	orderRepo := repo.NewOrderRepository(db)
	regSvc := service.NewRegistrationService(db, patientRepo, scheduleRepo, orderRepo)
	regHandler := handler.NewRegistrationHandler(regSvc)

	orderSvc := service.NewOrderService(db, orderRepo, scheduleRepo, patientRepo)
	orderHandler := handler.NewOrderHandler(orderSvc)

	// H5 API group — visitor phone auth
	h5 := r.Group("/api/v1")
	h5.Use(middleware.VisitorPhone())
	h5.POST("/patients", patientHandler.Create)
	h5.GET("/patients", patientHandler.List)
	h5.GET("/patients/:id", patientHandler.Get)
	h5.PUT("/patients/:id", patientHandler.Update)
	h5.DELETE("/patients/:id", patientHandler.Delete)

	h5.GET("/schedules", scheduleHandler.List)
	h5.GET("/schedules/:id", scheduleHandler.Get)

	h5.POST("/registrations", regHandler.Submit)
	h5.GET("/registrations/ticket/:order_no", regHandler.GetTicket)

	h5.GET("/orders", orderHandler.List)
	h5.GET("/orders/:id", orderHandler.Get)
	h5.PUT("/orders/:id/cancel", orderHandler.Cancel)
	h5.PUT("/orders/:id/change", orderHandler.Change)

	// Admin login — no auth required
	r.POST("/api/v1/admin/login", handler.AdminLogin)

	// Admin API group — JWT auth
	admin := r.Group("/api/v1/admin")
	admin.Use(middleware.JWTAuth())
	admin.GET("/patients", patientHandler.List)
	admin.GET("/patients/:id", patientHandler.Get)
	admin.PUT("/patients/:id", patientHandler.Update)
	admin.DELETE("/patients/:id", patientHandler.Delete)

	admin.POST("/schedules", scheduleHandler.Create)
	admin.GET("/schedules", scheduleHandler.List)
	admin.GET("/schedules/:id", scheduleHandler.Get)
	admin.PUT("/schedules/:id", scheduleHandler.Update)
	admin.DELETE("/schedules/:id", scheduleHandler.Delete)
	admin.POST("/schedules/:id/deduct", scheduleHandler.Deduct)
	admin.POST("/schedules/:id/rollback", scheduleHandler.Rollback)

	admin.GET("/orders", orderHandler.List)
	admin.GET("/orders/:id", orderHandler.Get)
	admin.PUT("/orders/:id/cancel", orderHandler.Cancel)
	admin.PUT("/orders/:id/change", orderHandler.Change)
	admin.PUT("/orders/:id/complete", orderHandler.Complete)

	return r
}
