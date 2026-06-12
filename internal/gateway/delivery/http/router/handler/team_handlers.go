package handler

import (
	"AgentManagmentSystem/internal/gateway/delivery/http/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *GatewayHandlers) ListUsers(c *gin.Context) {
	users, _ := h.team.ListUsers(c.Request.Context())
	res := make([]dto.UserResponse, len(users))
	for i, u := range users {
		res[i] = dto.UserResponse{UserID: u.UserID, Email: u.Email, Role: u.Role, Status: u.Status}
	}
	c.JSON(http.StatusOK, res)
}

func (h *GatewayHandlers) InviteUser(c *gin.Context) {
	c.JSON(http.StatusOK, dto.StatusResponse{Status: "invited"})
}
func (h *GatewayHandlers) UpdateUserRoles(c *gin.Context) {
	c.JSON(http.StatusOK, dto.StatusResponse{Status: "ok"})
}
func (h *GatewayHandlers) DeleteUser(c *gin.Context) {
	c.JSON(http.StatusOK, dto.SuccessResponse{Success: true})
}
