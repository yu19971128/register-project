package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestVisitorPhone_SetsContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/", nil)
	c.Request.Header.Set("X-Visitor-Phone", "13800138000")

	VisitorPhone()(c)

	if vp := c.GetString("visitor_phone"); vp != "13800138000" {
		t.Errorf("visitor_phone = %s, want 13800138000", vp)
	}
}

func TestVisitorPhone_EmptyHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/", nil)

	VisitorPhone()(c)

	if vp := c.GetString("visitor_phone"); vp != "" {
		t.Errorf("visitor_phone = %s, want empty", vp)
	}
}

func TestRequireVisitorPhone_RejectMissing(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/", nil)

	RequireVisitorPhone()(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want 401", w.Code)
	}
}

func TestRequireVisitorPhone_AllowPresent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/", nil)
	c.Request.Header.Set("X-Visitor-Phone", "13800138000")

	RequireVisitorPhone()(c)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", w.Code)
	}
	if vp := c.GetString("visitor_phone"); vp != "13800138000" {
		t.Errorf("visitor_phone = %s, want 13800138000", vp)
	}
}
