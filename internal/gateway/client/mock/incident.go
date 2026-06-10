package mock

import (
	"AgentManagmentSystem/internal/gateway/domain"
	"context"
	"time"
)

type IncidentMock struct{}

func NewIncidentMock() *IncidentMock {
	return &IncidentMock{}
}

func (m *IncidentMock) ListActiveIncidents(ctx context.Context) ([]domain.Incident, error) {
	return []domain.Incident{
		{
			IncidentID: "inc-001",
			DeviceID:   "d-1002",
			Type:       "DEAD_MAN_SWITCH_TRIGGERED",
			Severity:   "critical",
			Status:     "OPEN",
			CreatedAt:  time.Now().Add(-30 * time.Minute),
		},
		{
			IncidentID: "inc-002",
			DeviceID:   "d-1001",
			Type:       "ALERT_RULE_TRIGGERED",
			Severity:   "warning",
			Status:     "ACKNOWLEDGED",
			CreatedAt:  time.Now().Add(-1 * time.Hour),
		},
	}, nil
}

func (m *IncidentMock) AcknowledgeIncident(ctx context.Context, incidentID string) error {
	return nil
}

func (m *IncidentMock) ResolveIncident(ctx context.Context, incidentID string) error {
	return nil
}
