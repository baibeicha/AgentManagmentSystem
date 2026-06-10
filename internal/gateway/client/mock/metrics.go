package mock

import (
	"AgentManagmentSystem/internal/gateway/domain"
	"context"
	"fmt"
	"time"
)

type MetricsMock struct{}

func NewMetricsMock() *MetricsMock {
	return &MetricsMock{}
}

func (m *MetricsMock) GetCurrentMetrics(ctx context.Context, deviceID string) (*domain.MetricsSnapshot, error) {
	return &domain.MetricsSnapshot{
		CPUUsagePercent:     45.2,
		RAMUsageBytes:       8589934592,  // 8 GB
		RAMTotalBytes:       16106127360, // 15 GB
		DiskIOReadBytesSec:  1024000,     // 1 MB/s
		DiskIOWriteBytesSec: 512000,      // 500 KB/s
	}, nil
}

func (m *MetricsMock) GetHistoricalMetrics(ctx context.Context, deviceID, metricType string, from, to time.Time, step string) ([][]string, error) {
	var values [][]string
	currentTime := from

	for i := 0; i < 10; i++ {
		timestamp := fmt.Sprintf("%d", currentTime.Unix())
		val := fmt.Sprintf("%.1f", 20.0+float64(i*4))
		values = append(values, []string{timestamp, val})
		currentTime = currentTime.Add(5 * time.Minute)
	}
	return values, nil
}

func (m *MetricsMock) GetActiveProcesses(ctx context.Context, deviceID string) ([]domain.Process, error) {
	return []domain.Process{
		{PID: 1, Name: "systemd", CPUUsage: 0.1, RAMUsageBytes: 15482880},
		{PID: 5432, Name: "postgres", CPUUsage: 12.4, RAMUsageBytes: 1073741824},
		{PID: 8080, Name: "omniwatch-agent", CPUUsage: 2.5, RAMUsageBytes: 52428800},
	}, nil
}
