package handler

import (
	"net/http"
	"strconv"

	"clinic/models"
	"clinic/service"
	"github.com/gin-gonic/gin"
)

type PatientHandler struct {
	svc *service.PatientService
}

func NewPatientHandler(svc *service.PatientService) *PatientHandler {
	return &PatientHandler{svc: svc}
}

func (h *PatientHandler) Create(c *gin.Context) {
	var req models.Patient
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, "参数格式错误")
		return
	}
	req.VisitorPhone = c.GetString("visitor_phone")
	p, err := h.svc.CreatePatient(&req)
	if err != nil {
		Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	OK(c, p)
}

func (h *PatientHandler) List(c *gin.Context) {
	page, pageSize := ParsePage(c)
	var list []*models.Patient
	var total int
	var err error

	if vp := c.GetString("visitor_phone"); vp != "" {
		list, total, err = h.svc.ListPatientsByVisitorPhone(vp, page, pageSize)
	} else {
		list, total, err = h.svc.ListPatients(c.Query("keyword"), page, pageSize)
	}
	if err != nil {
		Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	OK(c, gin.H{"total": total, "list": list})
}

func (h *PatientHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	p, err := h.svc.GetPatient(id)
	if err != nil {
		Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	if p == nil {
		Error(c, http.StatusNotFound, "就诊人不存在")
		return
	}
	OK(c, p)
}

func (h *PatientHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var req models.Patient
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, "参数格式错误")
		return
	}
	req.ID = id
	if err := h.svc.UpdatePatient(&req); err != nil {
		Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	OK(c, nil)
}

func (h *PatientHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.svc.DeletePatient(id); err != nil {
		Error(c, http.StatusBadRequest, err.Error())
		return
	}
	OK(c, nil)
}
