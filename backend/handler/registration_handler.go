package handler

import (
	"net/http"

	"clinic/service"
	"github.com/gin-gonic/gin"
)

type RegistrationHandler struct {
	svc *service.RegistrationService
}

func NewRegistrationHandler(svc *service.RegistrationService) *RegistrationHandler {
	return &RegistrationHandler{svc: svc}
}

func (h *RegistrationHandler) Submit(c *gin.Context) {
	var req struct {
		ScheduleID   int64  `json:"schedule_id"`
		PatientID    int64  `json:"patient_id"`
		VisitorPhone string `json:"visitor_phone"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, "参数格式错误")
		return
	}
	if req.ScheduleID <= 0 || req.PatientID <= 0 {
		Error(c, http.StatusBadRequest, "号源或就诊人 ID 无效")
		return
	}
	// Prefer header visitor phone if present
	if vp := c.GetString("visitor_phone"); vp != "" {
		req.VisitorPhone = vp
	}
	if req.VisitorPhone == "" {
		Error(c, http.StatusBadRequest, "缺少访客手机号")
		return
	}

	result, err := h.svc.SubmitRegistration(req.ScheduleID, req.PatientID, req.VisitorPhone)
	if err != nil {
		msg := err.Error()
		switch msg {
		case "号源不存在", "就诊人不存在", "缺少必填字段或参数无效":
			Error(c, http.StatusBadRequest, msg)
		case "就诊人不属于当前访客", "无权查看该凭证":
			Error(c, http.StatusForbidden, msg)
		case "号源余量已为 0":
			Error(c, http.StatusConflict, msg)
		case "同一号源不可重复提交":
			Error(c, http.StatusTooManyRequests, msg)
		default:
			Error(c, http.StatusInternalServerError, msg)
		}
		return
	}
	OK(c, result)
}

func (h *RegistrationHandler) GetTicket(c *gin.Context) {
	orderNo := c.Param("order_no")
	if orderNo == "" {
		Error(c, http.StatusBadRequest, "订单号不能为空")
		return
	}
	visitorPhone := c.GetString("visitor_phone")
	if visitorPhone == "" {
		Error(c, http.StatusBadRequest, "缺少访客手机号")
		return
	}

	ticket, err := h.svc.GetTicket(orderNo, visitorPhone)
	if err != nil {
		msg := err.Error()
		if msg == "挂号凭证不存在" {
			Error(c, http.StatusNotFound, msg)
			return
		}
		if msg == "无权查看该凭证" {
			Error(c, http.StatusForbidden, msg)
			return
		}
		Error(c, http.StatusInternalServerError, msg)
		return
	}
	OK(c, ticket)
}
