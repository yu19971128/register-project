package handler

import (
	"net/http"
	"strconv"
	"strings"

	"clinic/models"
	"clinic/service"
	"github.com/gin-gonic/gin"
)

type ScheduleHandler struct {
	svc *service.ScheduleService
}

func NewScheduleHandler(svc *service.ScheduleService) *ScheduleHandler {
	return &ScheduleHandler{svc: svc}
}

func (h *ScheduleHandler) Create(c *gin.Context) {
	var req models.Schedule
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, "参数格式错误")
		return
	}
	if req.Date == "" || req.Department == "" || req.DoctorName == "" || req.StartTime == "" || req.EndTime == "" || req.TotalQuota <= 0 {
		Error(c, http.StatusBadRequest, "缺少必填字段或参数无效")
		return
	}
	if req.StartTime >= req.EndTime {
		Error(c, http.StatusBadRequest, "开始时间必须早于结束时间")
		return
	}
	s, err := h.svc.CreateSchedule(&req)
	if err != nil {
		if isConflict(err) {
			Error(c, http.StatusConflict, "同一医生同一时间段号源已存在")
			return
		}
		Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	OK(c, s)
}

func (h *ScheduleHandler) List(c *gin.Context) {
	page, pageSize := ParsePage(c)
	date := c.Query("date")
	if date == "" {
		date = c.GetString("schedule_date")
	}
	department := c.Query("department")
	list, total, err := h.svc.ListSchedules(date, department, page, pageSize)
	if err != nil {
		Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	OK(c, gin.H{"total": total, "list": list})
}

func (h *ScheduleHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	s, err := h.svc.GetSchedule(id)
	if err != nil {
		if err.Error() == "号源不存在" {
			Error(c, http.StatusNotFound, err.Error())
			return
		}
		Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	OK(c, s)
}

func (h *ScheduleHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var req models.Schedule
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, "参数格式错误")
		return
	}
	req.ID = id
	if err := h.svc.UpdateSchedule(&req); err != nil {
		if err.Error() == "号源不存在" {
			Error(c, http.StatusNotFound, err.Error())
			return
		}
		Error(c, http.StatusBadRequest, err.Error())
		return
	}
	OK(c, nil)
}

func (h *ScheduleHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.svc.DeleteSchedule(id); err != nil {
		if err.Error() == "号源不存在" {
			Error(c, http.StatusNotFound, err.Error())
			return
		}
		Error(c, http.StatusBadRequest, err.Error())
		return
	}
	OK(c, nil)
}

func (h *ScheduleHandler) Deduct(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	ok, err := h.svc.Deduct(id)
	if err != nil {
		Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	if !ok {
		Error(c, http.StatusConflict, "号源余量已为 0")
		return
	}
	OK(c, gin.H{"id": id, "deducted": true})
}

func (h *ScheduleHandler) Rollback(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	ok, err := h.svc.Rollback(id)
	if err != nil {
		Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	if !ok {
		Error(c, http.StatusBadRequest, "号源余量已达上限，无需回滚")
		return
	}
	OK(c, gin.H{"id": id, "rolled_back": true})
}

func isConflict(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "UNIQUE constraint failed")
}
