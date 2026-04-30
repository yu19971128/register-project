package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message"`
}

func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{Code: 200, Data: data, Message: "ok"})
}

func Error(c *gin.Context, code int, message string) {
	c.JSON(code, Response{Code: code, Message: message})
}

func ParsePage(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.Query("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(c.Query("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	return page, pageSize
}
