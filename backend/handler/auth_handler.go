package handler

import (
	"net/http"
	"os"

	"clinic/middleware"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func AdminLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, "请求参数错误")
		return
	}

	adminUser := os.Getenv("ADMIN_USER")
	adminPass := os.Getenv("ADMIN_PASS")
	if adminUser == "" {
		adminUser = "admin"
	}
	if adminPass == "" {
		adminPass = "admin123"
	}

	if req.Username != adminUser || req.Password != adminPass {
		Error(c, http.StatusUnauthorized, "用户名或密码错误")
		return
	}

	token, err := middleware.GenerateToken(adminUser)
	if err != nil {
		Error(c, http.StatusInternalServerError, "生成令牌失败")
		return
	}

	OK(c, LoginResponse{Token: token})
}
