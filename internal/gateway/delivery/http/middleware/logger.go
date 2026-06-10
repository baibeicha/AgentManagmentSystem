package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func RequestLogger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method

		if status >= 500 {
			logger.Error("HTTP Request Failed",
				slog.Int("status", status),
				slog.String("method", method),
				slog.String("path", path),
				slog.String("query", query),
				slog.String("ip", clientIP),
				slog.Duration("latency", latency),
				slog.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			)
			return
		}

		logger.Info("HTTP Request",
			slog.Int("status", status),
			slog.String("method", method),
			slog.String("path", path),
			slog.String("ip", clientIP),
			slog.Duration("latency", latency),
		)
	}
}
