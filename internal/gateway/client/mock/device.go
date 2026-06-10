package mock

import (
	"AgentManagmentSystem/internal/gateway/domain"
	"bytes"
	"context"
	"io"
	"time"
)

type DeviceMock struct{}

func NewDeviceMock() *DeviceMock {
	return &DeviceMock{}
}

func (m *DeviceMock) ListDevices(ctx context.Context) ([]domain.Device, error) {
	return []domain.Device{
		{
			DeviceID:      "d-1001",
			Alias:         "prod-web-01",
			OS:            "linux",
			Arch:          "amd64",
			IPAddress:     "192.168.1.15",
			Status:        "online",
			AgentVersion:  "v1.2.0",
			GroupID:       "g-01",
			Tags:          []string{"web", "nginx", "prod"},
			LastHeartbeat: time.Now(),
		},
		{
			DeviceID:      "d-1002",
			Alias:         "prod-db-master",
			OS:            "linux",
			Arch:          "amd64",
			IPAddress:     "192.168.1.20",
			Status:        "offline",
			AgentVersion:  "v1.1.9",
			GroupID:       "g-02",
			Tags:          []string{"database", "postgres", "critical"},
			LastHeartbeat: time.Now().Add(-2 * time.Hour),
		},
	}, nil
}

func (m *DeviceMock) GetDevice(ctx context.Context, deviceID string) (*domain.Device, error) {
	return &domain.Device{
		DeviceID:      deviceID,
		Alias:         "mock-device",
		OS:            "linux",
		Arch:          "amd64",
		IPAddress:     "10.0.0.5",
		Status:        "online",
		AgentVersion:  "v1.2.0",
		LastHeartbeat: time.Now(),
	}, nil
}

func (m *DeviceMock) DeleteDevice(ctx context.Context, deviceID string) error {
	return nil
}

func (m *DeviceMock) UpdateDeviceMetadata(ctx context.Context, deviceID string, alias, groupID string, tags []string) error {
	return nil
}

func (m *DeviceMock) CreateBootstrapToken(ctx context.Context) (string, time.Time, error) {
	return "bst_mock123456789", time.Now().Add(1 * time.Hour), nil
}

func (m *DeviceMock) EnrollAgent(ctx context.Context, provisioningKey, csrPEM string) (string, string, error) {
	return "-----BEGIN CERTIFICATE-----\nMOCK_CLIENT_CERT\n-----END CERTIFICATE-----",
		"-----BEGIN CERTIFICATE-----\nMOCK_CA_CERT\n-----END CERTIFICATE-----", nil
}

func (m *DeviceMock) TriggerDeviceUpdate(ctx context.Context, deviceID, targetVersion string) error {
	return nil
}

func (m *DeviceMock) ListDeviceGroups(ctx context.Context) ([]domain.DeviceGroup, error) {
	return []domain.DeviceGroup{
		{GroupID: "g-01", Name: "Web Servers", Description: "Frontend nginx nodes"},
		{GroupID: "g-02", Name: "Databases", Description: "PostgreSQL clusters"},
	}, nil
}

func (m *DeviceMock) CreateDeviceGroup(ctx context.Context, name, description string) error {
	return nil
}

func (m *DeviceMock) ListDiscoveryCandidates(ctx context.Context) ([]domain.DiscoveryCandidate, error) {
	return []domain.DiscoveryCandidate{
		{CandidateID: "cand-01", IPAddress: "192.168.1.100", MACAddress: "00:1B:44:11:3A:B7", DiscoveredByDeviceID: "d-1001"},
	}, nil
}

func (m *DeviceMock) ApproveDiscoveryCandidate(ctx context.Context, candidateID string) error {
	return nil
}

func (m *DeviceMock) IgnoreDiscoveryCandidate(ctx context.Context, candidateID string) error {
	return nil
}

func (m *DeviceMock) GetAgentReleaseStream(ctx context.Context, os, arch string) (io.ReadCloser, error) {
	return io.NopCloser(bytes.NewReader([]byte("mock_binary_data_executable"))), nil
}

func (m *DeviceMock) GetLatestReleaseVersion(ctx context.Context) (*domain.AgentRelease, error) {
	return &domain.AgentRelease{Version: "v1.3.0", SHA256Checksum: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"}, nil
}

func (m *DeviceMock) UploadAgentRelease(ctx context.Context, version, os, arch string, fileStream io.Reader) error {
	return nil
}
