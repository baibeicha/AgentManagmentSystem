package mock

import (
	"AgentManagmentSystem/internal/gateway/domain"
	"context"
	"time"
)

type AuditMock struct{}

func NewAuditMock() *AuditMock {
	return &AuditMock{}
}

func (m *AuditMock) ListAuditLogs(ctx context.Context) ([]domain.AuditLog, error) {
	return []domain.AuditLog{
		{
			LogID:      "log-001",
			Timestamp:  time.Now().Add(-5 * time.Minute),
			UserID:     "u-1111-2222-3333-4444",
			ActionType: "COMMAND_EXECUTE",
			DeviceID:   "d-1001",
			Payload:    "{\"command\": \"systemctl restart nginx\"}",
			Status:     "SUCCESS",
		},
		{
			LogID:      "log-002",
			Timestamp:  time.Now().Add(-1 * time.Hour),
			UserID:     "u-5555-6666-7777-8888",
			ActionType: "DEVICE_DELETE",
			DeviceID:   "d-1002",
			Payload:    "{}",
			Status:     "BLOCKED_BY_POLICY",
		},
	}, nil
}

func (m *AuditMock) VerifyAuditIntegrity(ctx context.Context) (string, string, error) {
	return "VALID", "", nil
}

func (m *AuditMock) ExportConfig(ctx context.Context) (string, error) {
	mockYaml := `
version: 1
groups:
  - name: Web Servers
rules:
  - name: High CPU
    metric: cpu
    operator: ">"
    threshold: 90
`
	return mockYaml, nil
}

func (m *AuditMock) ImportConfig(ctx context.Context, yamlManifest string) error {
	return nil
}
