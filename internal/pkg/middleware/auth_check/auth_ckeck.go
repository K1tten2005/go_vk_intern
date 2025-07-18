package authcheck

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/K1tten2005/go_vk_intern/internal/pkg/utils/logger"
	"github.com/K1tten2005/go_vk_intern/internal/pkg/utils/send_err"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		loggerVar := logger.GetLoggerFromContext(c.Request.Context()).With(slog.String("func", logger.GetFuncName()))

		_, err := c.Cookie("VKJWT")
		if err != nil {
			if err == http.ErrNoCookie {
				logger.LogHandlerError(loggerVar, fmt.Errorf("no token: %w", err), http.StatusBadRequest)
				send_err.SendError(c.Writer, "no token", http.StatusBadRequest)
				return
			}
			logger.LogHandlerError(loggerVar, fmt.Errorf("error while parsing cookie: %w", err), http.StatusBadRequest)
			send_err.SendError(c.Writer, "error while parsing cookie", http.StatusBadRequest)
			c.Abort()
			return
		}
		c.Next()
	}
}
