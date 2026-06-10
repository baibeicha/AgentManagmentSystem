package middleware

import (
	"AgentManagmentSystem/internal/gateway/delivery/http/dto"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

func Recovery(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("PANIC RECOVERED",
					slog.Any("error", err),
					slog.String("stack", string(debug.Stack())),
					slog.String("path", c.Request.URL.Path),
				)

				c.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{
					Error: "Internal Server Error. Our team has been notified.",
				})
			}
		}()

		c.Next()
	}
}
