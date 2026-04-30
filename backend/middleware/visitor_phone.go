package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func VisitorPhone() gin.HandlerFunc {
	return func(c *gin.Context) {
		phone := c.GetHeader("X-Visitor-Phone")
		if phone != "" {
			c.Set("visitor_phone", phone)
		}
		c.Next()
	}
}

func RequireVisitorPhone() gin.HandlerFunc {
	return func(c *gin.Context) {
		phone := c.GetHeader("X-Visitor-Phone")
		if phone == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "缺少访客手机号"})
			return
		}
		c.Set("visitor_phone", phone)
		c.Next()
	}
}
