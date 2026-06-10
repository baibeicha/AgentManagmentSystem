package router

import (
	"AgentManagmentSystem/internal/gateway/delivery/http/middleware"
	"AgentManagmentSystem/pkg/config"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Config, log *slog.Logger) *gin.Engine {
	r := gin.New()

	r.Use(middleware.RequestLogger(log), middleware.Recovery(log))

	r.GET("/ping", func(c *gin.Context) {
		time.Sleep(3 * time.Second)
		c.String(http.StatusOK, "pong")
	})

	return r
}
