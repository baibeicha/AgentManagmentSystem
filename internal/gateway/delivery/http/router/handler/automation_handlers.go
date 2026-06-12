package handler

import (
	"AgentManagmentSystem/internal/gateway/delivery/http/dto"
	"AgentManagmentSystem/internal/gateway/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *GatewayHandlers) ListAlertRules(c *gin.Context) {
	rules, _ := h.automation.ListAlertRules(c.Request.Context())
	res := make([]dto.AlertRuleResponse, len(rules))
	for i, r := range rules {
		res[i] = dto.AlertRuleResponse{RuleID: r.RuleID, Name: r.Name, Metric: r.Metric, Operator: r.Operator, Threshold: r.Threshold, Duration: r.Duration, PlaybookID: r.PlaybookID}
	}
	c.JSON(http.StatusOK, res)
}

func (h *GatewayHandlers) CreateAlertRule(c *gin.Context) {
	var req dto.AlertRuleRequest
	_ = c.ShouldBindJSON(&req)
	_ = h.automation.CreateAlertRule(c.Request.Context(), domain.AlertRule{Name: req.Name, Metric: req.Metric, Operator: req.Operator, Threshold: req.Threshold, Duration: req.Duration, PlaybookID: req.PlaybookID})
	c.Status(http.StatusCreated)
}

func (h *GatewayHandlers) UpdateAlertRule(c *gin.Context) {
	var req dto.AlertRuleRequest
	_ = c.ShouldBindJSON(&req)
	_ = h.automation.UpdateAlertRule(c.Request.Context(), c.Param("id"), domain.AlertRule{Name: req.Name, Metric: req.Metric, Operator: req.Operator, Threshold: req.Threshold, Duration: req.Duration, PlaybookID: req.PlaybookID})
	c.JSON(http.StatusOK, dto.StatusResponse{Status: "ok"})
}

func (h *GatewayHandlers) DeleteAlertRule(c *gin.Context) {
	_ = h.automation.DeleteAlertRule(c.Request.Context(), c.Param("id"))
	c.JSON(http.StatusOK, dto.SuccessResponse{Success: true})
}

func (h *GatewayHandlers) ListPlaybooks(c *gin.Context) {
	playbooks, _ := h.automation.ListPlaybooks(c.Request.Context())
	res := make([]dto.PlaybookResponse, len(playbooks))
	for i, p := range playbooks {
		steps := make([]dto.PlaybookStep, len(p.Steps))
		for j, s := range p.Steps {
			steps[j] = dto.PlaybookStep{Type: s.Type, TargetID: s.TargetID, Delay: s.Delay}
		}
		res[i] = dto.PlaybookResponse{PlaybookID: p.PlaybookID, Name: p.Name, Steps: steps, CreatedAt: p.CreatedAt}
	}
	c.JSON(http.StatusOK, res)
}

func (h *GatewayHandlers) CreatePlaybook(c *gin.Context) {
	var req dto.PlaybookRequest
	_ = c.ShouldBindJSON(&req)
	c.Status(http.StatusCreated)
}

func (h *GatewayHandlers) GetPlaybook(c *gin.Context) {
	p, _ := h.automation.GetPlaybook(c.Request.Context(), c.Param("id"))
	c.JSON(http.StatusOK, dto.PlaybookResponse{PlaybookID: p.PlaybookID, Name: p.Name, CreatedAt: p.CreatedAt})
}

func (h *GatewayHandlers) UpdatePlaybook(c *gin.Context) {
	c.JSON(http.StatusOK, dto.StatusResponse{Status: "ok"})
}

func (h *GatewayHandlers) DeletePlaybook(c *gin.Context) {
	_ = h.automation.DeletePlaybook(c.Request.Context(), c.Param("id"))
	c.JSON(http.StatusOK, dto.SuccessResponse{Success: true})
}

func (h *GatewayHandlers) ListCronJobs(c *gin.Context) {
	crons, _ := h.automation.ListCronJobs(c.Request.Context())
	res := make([]dto.CronResponse, len(crons))
	for i, cr := range crons {
		res[i] = dto.CronResponse{CronID: cr.CronID, Expression: cr.Expression, ScriptID: cr.ScriptID, DeviceGroupID: cr.DeviceGroupID, IsActive: cr.IsActive, CreatedAt: cr.CreatedAt}
	}
	c.JSON(http.StatusOK, res)
}

func (h *GatewayHandlers) CreateCronJob(c *gin.Context) { c.Status(http.StatusCreated) }
func (h *GatewayHandlers) UpdateCronJob(c *gin.Context) {
	c.JSON(http.StatusOK, dto.StatusResponse{Status: "ok"})
}
func (h *GatewayHandlers) DeleteCronJob(c *gin.Context) {
	c.JSON(http.StatusOK, dto.SuccessResponse{Success: true})
}
