package handler

import (
	"AgentManagmentSystem/internal/gateway/delivery/http/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *GatewayHandlers) Register(c *gin.Context) {
	var req dto.RegistrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	success, err := h.auth.Register(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.SuccessResponse{Success: success})
}

func (h *GatewayHandlers) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	access, refresh, err := h.auth.Login(c.Request.Context(), req.Login, req.Password, req.DeviceID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.TokenResponse{AccessToken: access, RefreshToken: refresh})
}

func (h *GatewayHandlers) Refresh(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	access, refresh, err := h.auth.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.TokenResponse{AccessToken: access, RefreshToken: refresh})
}

func (h *GatewayHandlers) Logout(c *gin.Context) {
	var req dto.LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	_ = h.auth.Logout(c.Request.Context(), req.RefreshToken)
	c.JSON(http.StatusOK, dto.SuccessResponse{Success: true})
}

func (h *GatewayHandlers) GetMe(c *gin.Context) {
	u, err := h.auth.GetMe(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.UserProfileResponse{
		UserID: u.UserID, Login: u.Email, TenantID: u.Status, Role: u.Role,
	})
}

func (h *GatewayHandlers) GetPermissions(c *gin.Context) {
	perms, _ := h.auth.GetPermissions(c.Request.Context())
	c.JSON(http.StatusOK, dto.PermissionsResponse{Permissions: perms})
}

func (h *GatewayHandlers) Setup2FA(c *gin.Context) {
	secret, url, _ := h.auth.Setup2FA(c.Request.Context())
	c.JSON(http.StatusOK, dto.Setup2FAResponse{Secret: secret, OTPAuthURL: url})
}

func (h *GatewayHandlers) Verify2FA(c *gin.Context) {
	var req dto.Verify2FARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	status, err := h.auth.Verify2FA(c.Request.Context(), req.Code, req.ActionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.StatusResponse{Status: status})
}

func (h *GatewayHandlers) Disable2FA(c *gin.Context) {
	var req dto.Disable2FARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	if err := h.auth.Disable2FA(c.Request.Context(), req.Code); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.StatusResponse{Status: "disabled"})
}
