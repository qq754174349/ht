package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func PreflightMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.GetHeader("X-Forwarded-Method")
		if method == http.MethodOptions || c.Request.Method == http.MethodOptions {
			c.Status(http.StatusNoContent)
			c.Abort()
			return
		}
		c.Next()
	}
}
