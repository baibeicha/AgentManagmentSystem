package handler

import (
	"AgentManagmentSystem/internal/gateway/delivery/http/dto"
	"AgentManagmentSystem/internal/gateway/service"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *GatewayHandlers) ListDevices(c *gin.Context) {
	devices, _ := h.device.ListDevices(c.Request.Context())
	res := make([]dto.DeviceResponse, len(devices))
	for i, d := range devices {
		res[i] = dto.DeviceResponse{
			DeviceID: d.DeviceID, Alias: d.Alias, OS: d.OS, Arch: d.Arch,
			IPAddress: d.IPAddress, Status: d.Status, AgentVersion: d.AgentVersion,
			GroupID: d.GroupID, Tags: d.Tags, LastHeartbeat: d.LastHeartbeat,
		}
	}
	c.JSON(http.StatusOK, res)
}

func (h *GatewayHandlers) GetDevice(c *gin.Context) {
	d, _ := h.device.GetDevice(c.Request.Context(), c.Param("id"))
	c.JSON(http.StatusOK, dto.DeviceResponse{
		DeviceID: d.DeviceID, Alias: d.Alias, OS: d.OS, Arch: d.Arch,
		IPAddress: d.IPAddress, Status: d.Status, AgentVersion: d.AgentVersion,
		LastHeartbeat: d.LastHeartbeat,
	})
}

func (h *GatewayHandlers) DeleteDevice(c *gin.Context) {
	_ = h.device.DeleteDevice(c.Request.Context(), c.Param("id"))
	c.JSON(http.StatusOK, dto.SuccessResponse{Success: true})
}

func (h *GatewayHandlers) UpdateDeviceMetadata(c *gin.Context) {
	var req dto.UpdateDeviceMetadataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	_ = h.device.UpdateDeviceMetadata(c.Request.Context(), c.Param("id"), req.Alias, req.GroupID, req.Tags)
	c.JSON(http.StatusOK, dto.StatusResponse{Status: "ok"})
}

func (h *GatewayHandlers) CreateBootstrapToken(c *gin.Context) {
	token, exp, _ := h.device.CreateBootstrapToken(c.Request.Context())
	c.JSON(http.StatusOK, dto.BootstrapResponse{BootstrapToken: token, ExpiresAt: exp})
}

func (h *GatewayHandlers) EnrollAgent(c *gin.Context) {
	var req dto.AgentEnrollRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	client, ca, err := h.device.EnrollAgent(c.Request.Context(), req.ProvisioningKey, req.CSRPem)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.AgentEnrollResponse{ClientCertificatePEM: client, CACertificatePEM: ca})
}

func (h *GatewayHandlers) TriggerDeviceUpdate(c *gin.Context) {
	var req dto.UpdateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	_ = h.device.TriggerDeviceUpdate(c.Request.Context(), req.DeviceID, req.TargetVersion)
	c.Status(http.StatusAccepted)
}

func (h *GatewayHandlers) ListDeviceGroups(c *gin.Context) {
	groups, _ := h.device.ListDeviceGroups(c.Request.Context())
	res := make([]dto.DeviceGroupResponse, len(groups))
	for i, g := range groups {
		res[i] = dto.DeviceGroupResponse{GroupID: g.GroupID, Name: g.Name, Description: g.Description}
	}
	c.JSON(http.StatusOK, res)
}

func (h *GatewayHandlers) CreateDeviceGroup(c *gin.Context) {
	var req dto.CreateDeviceGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	_ = h.device.CreateDeviceGroup(c.Request.Context(), req.Name, req.Description)
	c.Status(http.StatusCreated)
}

func (h *GatewayHandlers) DownloadAgent(c *gin.Context) {
	stream, err := h.device.(service.ReleaseService).GetAgentReleaseStream(c.Request.Context(), c.Param("os"), c.Param("arch"))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "release not found"})
		return
	}
	defer stream.Close()
	c.Header("Content-Disposition", "attachment; filename=agent.bin")
	c.Header("Content-Type", "application/octet-stream")
	_, _ = io.Copy(c.Writer, stream)
}

func (h *GatewayHandlers) LatestAgentVersion(c *gin.Context) {
	rel, _ := h.device.(service.ReleaseService).GetLatestReleaseVersion(c.Request.Context())
	c.JSON(http.StatusOK, dto.LatestReleaseResponse{Version: rel.Version, Sha256Checksum: rel.SHA256Checksum})
}

func (h *GatewayHandlers) UploadAgentRelease(c *gin.Context) {
	file, _ := c.FormFile("file")
	opened, _ := file.Open()
	defer opened.Close()
	_ = h.device.(service.ReleaseService).UploadAgentRelease(c.Request.Context(), c.PostForm("version"), c.PostForm("os"), c.PostForm("arch"), opened)
	c.Status(http.StatusCreated)
}
