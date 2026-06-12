package handler

import (
	"AgentManagmentSystem/internal/gateway/delivery/http/dto"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *GatewayHandlers) GetCurrentMetrics(c *gin.Context) {
	m, _ := h.metrics.GetCurrentMetrics(c.Request.Context(), c.Param("id"))
	c.JSON(http.StatusOK, dto.MetricsCurrentResponse{
		CPUUsagePercent: m.CPUUsagePercent, RAMUsageBytes: m.RAMUsageBytes, RAMTotalBytes: m.RAMTotalBytes,
		DiskIOReadBytesSec: m.DiskIOReadBytesSec, DiskIOWriteBytesSec: m.DiskIOWriteBytesSec,
	})
}

func (h *GatewayHandlers) GetHistoricalMetrics(c *gin.Context) {
	metrics, _ := h.metrics.GetHistoricalMetrics(c.Request.Context(), c.Param("id"), c.Query("metric_type"), time.Now().Add(-1*time.Hour), time.Now(), "1m")
	c.JSON(http.StatusOK, dto.TimeseriesResponse{Metric: c.Query("metric_type"), Values: metrics})
}

func (h *GatewayHandlers) GetActiveProcesses(c *gin.Context) {
	processes, _ := h.metrics.GetActiveProcesses(c.Request.Context(), c.Param("id"))
	res := make([]dto.ProcessResponse, len(processes))
	for i, p := range processes {
		res[i] = dto.ProcessResponse{PID: p.PID, Name: p.Name, CPUUsage: p.CPUUsage, RAMUsageBytes: p.RAMUsageBytes}
	}
	c.JSON(http.StatusOK, res)
}
