package mock

import (
	"AgentManagmentSystem/internal/gateway/domain"
	"context"
	"time"
)

type AutomationMock struct{}

func NewAutomationMock() *AutomationMock {
	return &AutomationMock{}
}

func (m *AutomationMock) ListAlertRules(ctx context.Context) ([]domain.AlertRule, error) {
	return []domain.AlertRule{
		{
			RuleID:     "rule-01",
			Name:       "High CPU Usage Alert",
			Metric:     "cpu",
			Operator:   ">",
			Threshold:  90.0,
			Duration:   "5m",
			PlaybookID: "playbook-01",
		},
		{
			RuleID:    "rule-02",
			Name:      "Low Disk Space",
			Metric:    "disk",
			Operator:  ">",
			Threshold: 95.0,
			Duration:  "10m",
		},
	}, nil
}

func (m *AutomationMock) CreateAlertRule(ctx context.Context, rule domain.AlertRule) error {
	return nil
}

func (m *AutomationMock) UpdateAlertRule(ctx context.Context, ruleID string, rule domain.AlertRule) error {
	return nil
}

func (m *AutomationMock) DeleteAlertRule(ctx context.Context, ruleID string) error {
	return nil
}

func (m *AutomationMock) ListPlaybooks(ctx context.Context) ([]domain.Playbook, error) {
	return []domain.Playbook{
		{
			PlaybookID: "playbook-01",
			Name:       "Restart Nginx & Clean Cache",
			Steps: []domain.PlaybookStep{
				{Type: "EXECUTE_SCRIPT", TargetID: "script-01", Delay: "0s"},
				{Type: "NOTIFY_CHANNEL", TargetID: "channel-01", Delay: "5s"},
			},
			CreatedAt: time.Now().Add(-24 * time.Hour),
		},
	}, nil
}

func (m *AutomationMock) CreatePlaybook(ctx context.Context, name string, steps []domain.PlaybookStep) error {
	return nil
}

func (m *AutomationMock) GetPlaybook(ctx context.Context, playbookID string) (*domain.Playbook, error) {
	return &domain.Playbook{
		PlaybookID: playbookID,
		Name:       "Mock Playbook",
		Steps: []domain.PlaybookStep{
			{Type: "EXECUTE_SCRIPT", TargetID: "mock-script", Delay: "0s"},
		},
		CreatedAt: time.Now(),
	}, nil
}

func (m *AutomationMock) UpdatePlaybook(ctx context.Context, playbookID string, name string, steps []domain.PlaybookStep) error {
	return nil
}

func (m *AutomationMock) DeletePlaybook(ctx context.Context, playbookID string) error {
	return nil
}

func (m *AutomationMock) ListCronJobs(ctx context.Context) ([]domain.CronJob, error) {
	return []domain.CronJob{
		{
			CronID:        "cron-01",
			Expression:    "0 3 * * 5", // Каждую пятницу в 3 ночи
			ScriptID:      "script-01",
			DeviceGroupID: "g-01",
			IsActive:      true,
			CreatedAt:     time.Now().Add(-48 * time.Hour),
		},
	}, nil
}

func (m *AutomationMock) CreateCronJob(ctx context.Context, expression, scriptID, deviceGroupID string) error {
	return nil
}

func (m *AutomationMock) UpdateCronJob(ctx context.Context, cronID string, expression, scriptID, deviceGroupID string, isActive bool) error {
	return nil
}

func (m *AutomationMock) DeleteCronJob(ctx context.Context, cronID string) error {
	return nil
}
