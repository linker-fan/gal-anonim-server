package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestJwtMiddleware(t *testing.T) {
	s := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(s)

	r.Use(authmiddleware)
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	c.Request, _ = http.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(s, c.Request)
}
