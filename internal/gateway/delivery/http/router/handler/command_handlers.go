package handler

import (
	"AgentManagmentSystem/internal/gateway/delivery/http/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *GatewayHandlers) ExecuteCommand(c *gin.Context) {
	var req dto.ExecuteCommandRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	res, err := h.command.ExecuteCommand(c.Request.Context(), req.DeviceID, req.Payload)
	if err != nil {
		c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.CommandExecuteResponse{ExitCode: res.ExitCode, Stdout: res.Stdout, Stderr: res.Stderr})
}

func (h *GatewayHandlers) ListScripts(c *gin.Context) {
	scripts, _ := h.command.ListScripts(c.Request.Context())
	res := make([]dto.ScriptTemplateResponse, len(scripts))
	for i, s := range scripts {
		res[i] = dto.ScriptTemplateResponse{ScriptID: s.ScriptID, Name: s.Name, Description: s.Description, Content: s.Content, Interpreter: s.Interpreter}
	}
	c.JSON(http.StatusOK, res)
}

func (h *GatewayHandlers) CreateScript(c *gin.Context) {
	var req dto.ScriptTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	_ = h.command.CreateScript(c.Request.Context(), req.Name, req.Description, req.Content, req.Interpreter)
	c.Status(http.StatusCreated)
}

func (h *GatewayHandlers) GetScript(c *gin.Context) {
	s, _ := h.command.GetScript(c.Request.Context(), c.Param("id"))
	c.JSON(http.StatusOK, dto.ScriptTemplateResponse{ScriptID: s.ScriptID, Name: s.Name, Content: s.Content, Interpreter: s.Interpreter})
}

func (h *GatewayHandlers) UpdateScript(c *gin.Context) {
	var req dto.ScriptTemplateRequest
	_ = c.ShouldBindJSON(&req)
	_ = h.command.UpdateScript(c.Request.Context(), c.Param("id"), req.Name, req.Description, req.Content, req.Interpreter)
	c.JSON(http.StatusOK, dto.StatusResponse{Status: "ok"})
}

func (h *GatewayHandlers) DeleteScript(c *gin.Context) {
	_ = h.command.DeleteScript(c.Request.Context(), c.Param("id"))
	c.JSON(http.StatusOK, dto.SuccessResponse{Success: true})
}
