package mock

import (
	"AgentManagmentSystem/internal/gateway/domain"
	"context"
)

type CommandMock struct{}

func NewCommandMock() *CommandMock {
	return &CommandMock{}
}

func (m *CommandMock) ExecuteCommand(ctx context.Context, deviceID, payload string) (*domain.CommandResult, error) {
	return &domain.CommandResult{
		ExitCode: 0,
		Stdout:   "mock execution successful for: " + payload + "\n",
		Stderr:   "",
	}, nil
}

func (m *CommandMock) ListScripts(ctx context.Context) ([]domain.ScriptTemplate, error) {
	return []domain.ScriptTemplate{
		{
			ScriptID:    "script-01",
			Name:        "Clear Docker Cache",
			Description: "Removes unused docker images and volumes",
			Content:     "#!/bin/bash\ndocker system prune -af",
			Interpreter: "bash",
		},
	}, nil
}

func (m *CommandMock) CreateScript(ctx context.Context, name, description, content, interpreter string) error {
	return nil
}

func (m *CommandMock) GetScript(ctx context.Context, scriptID string) (*domain.ScriptTemplate, error) {
	return &domain.ScriptTemplate{
		ScriptID:    scriptID,
		Name:        "Mock Script",
		Content:     "echo 'hello world'",
		Interpreter: "bash",
	}, nil
}

func (m *CommandMock) UpdateScript(ctx context.Context, scriptID string, name, description, content, interpreter string) error {
	return nil
}

func (m *CommandMock) DeleteScript(ctx context.Context, scriptID string) error {
	return nil
}
