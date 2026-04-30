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

	patientRepo := repo.NewPatientRepository(db)
	patientSvc := service.NewPatientService(patientRepo)
	patientHandler := handler.NewPatientHandler(patientSvc)

	scheduleRepo := repo.NewScheduleRepository(db)
	scheduleSvc := service.NewScheduleService(scheduleRepo)
	scheduleHandler := handler.NewScheduleHandler(scheduleSvc)

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

	return r
}
