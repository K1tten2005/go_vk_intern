package csp

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const CSP = "default-src 'none'; script-src 'self'; connect-src 'self'; img-src 'self'; style-src 'self'; base-uri 'self'; form-action 'self'"

func CspMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Security-Policy", CSP)
		if c.Request.Method == http.MethodOptions {
			return
		}
		c.Next()
	}
}
