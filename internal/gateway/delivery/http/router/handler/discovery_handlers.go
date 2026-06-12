package handler

import (
	"AgentManagmentSystem/internal/gateway/delivery/http/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *GatewayHandlers) ListDiscoveryCandidates(c *gin.Context) {
	candidates, _ := h.discovery.ListDiscoveryCandidates(c.Request.Context())
	res := make([]dto.DiscoveryCandidateResponse, len(candidates))
	for i, cd := range candidates {
		res[i] = dto.DiscoveryCandidateResponse{
			CandidateID: cd.CandidateID, IPAddress: cd.IPAddress,
			MACAddress: cd.MACAddress, DiscoveredByDeviceID: cd.DiscoveredByDeviceID,
		}
	}
	c.JSON(http.StatusOK, res)
}

func (h *GatewayHandlers) ApproveDiscoveryCandidate(c *gin.Context) {
	_ = h.discovery.ApproveDiscoveryCandidate(c.Request.Context(), c.Param("id"))
	c.JSON(http.StatusOK, dto.StatusResponse{Status: "approved"})
}

func (h *GatewayHandlers) IgnoreDiscoveryCandidate(c *gin.Context) {
	_ = h.discovery.IgnoreDiscoveryCandidate(c.Request.Context(), c.Param("id"))
	c.JSON(http.StatusOK, dto.StatusResponse{Status: "ignored"})
}
