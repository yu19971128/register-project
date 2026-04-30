package handler

import (
	"net/http"
	"strconv"

	"clinic/service"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	svc *service.OrderService
}

func NewOrderHandler(svc *service.OrderService) *OrderHandler {
	return &OrderHandler{svc: svc}
}

func (h *OrderHandler) List(c *gin.Context) {
	page, pageSize := ParsePage(c)
	isAdmin := c.GetString("admin_id") != ""
	visitorPhone := c.GetString("visitor_phone")

	res, err := h.svc.List(
		c.Query("date"),
		c.Query("department"),
		c.Query("doctor_name"),
		c.Query("status"),
		c.Query("keyword"),
		page, pageSize,
		isAdmin, visitorPhone,
	)
	if err != nil {
		Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	OK(c, res)
}

func (h *OrderHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	isAdmin := c.GetString("admin_id") != ""
	visitorPhone := c.GetString("visitor_phone")

	res, err := h.svc.GetDetail(id, isAdmin, visitorPhone)
	if err != nil {
		msg := err.Error()
		switch msg {
		case "订单不存在":
			Error(c, http.StatusNotFound, msg)
		case "无权查看该订单":
			Error(c, http.StatusForbidden, msg)
		default:
			Error(c, http.StatusInternalServerError, msg)
		}
		return
	}
	OK(c, res)
}

func (h *OrderHandler) Cancel(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	isAdmin := c.GetString("admin_id") != ""
	visitorPhone := c.GetString("visitor_phone")
	operatedBy := c.GetString("admin_id")

	var req struct {
		Reason string `json:"reason"`
	}
	c.ShouldBindJSON(&req)

	res, err := h.svc.Cancel(id, req.Reason, isAdmin, visitorPhone, operatedBy)
	if err != nil {
		msg := err.Error()
		switch msg {
		case "订单不存在", "新号源不存在":
			Error(c, http.StatusNotFound, msg)
		case "无权退号":
			Error(c, http.StatusForbidden, msg)
		case "已超过退号时限", "原订单已超过退号时限", "新号源余量已为 0":
			Error(c, http.StatusBadRequest, msg)
		default:
			Error(c, http.StatusInternalServerError, msg)
		}
		return
	}
	OK(c, res)
}

func (h *OrderHandler) Change(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	isAdmin := c.GetString("admin_id") != ""
	visitorPhone := c.GetString("visitor_phone")

	var req struct {
		NewScheduleID int64 `json:"new_schedule_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, "参数格式错误")
		return
	}
	if req.NewScheduleID <= 0 {
		Error(c, http.StatusBadRequest, "新号源 ID 无效")
		return
	}

	res, err := h.svc.Change(id, req.NewScheduleID, isAdmin, visitorPhone)
	if err != nil {
		msg := err.Error()
		switch msg {
		case "订单不存在", "新号源不存在":
			Error(c, http.StatusNotFound, msg)
		case "无权改号":
			Error(c, http.StatusForbidden, msg)
		case "原订单已超过退号时限", "新号源余量已为 0":
			Error(c, http.StatusBadRequest, msg)
		default:
			Error(c, http.StatusInternalServerError, msg)
		}
		return
	}
	OK(c, res)
}
