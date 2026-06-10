package mock

import (
	"AgentManagmentSystem/internal/gateway/domain"
	"context"
)

type NotificationMock struct{}

func NewNotificationMock() *NotificationMock {
	return &NotificationMock{}
}

func (m *NotificationMock) ListNotificationChannels(ctx context.Context) ([]domain.NotificationChannel, error) {
	return []domain.NotificationChannel{
		{
			ChannelID:   "channel-01",
			Name:        "DevOps Telegram Bot",
			Type:        "telegram",
			Destination: "https://api.telegram.org/bot123/sendMessage?chat_id=-1005555",
			IsActive:    true,
		},
		{
			ChannelID:   "channel-02",
			Name:        "Incident Discord Webhook",
			Type:        "discord",
			Destination: "https://discord.com/api/webhooks/123/abc",
			IsActive:    true,
		},
	}, nil
}

func (m *NotificationMock) CreateNotificationChannel(ctx context.Context, channel domain.NotificationChannel) error {
	return nil
}

func (m *NotificationMock) UpdateNotificationChannel(ctx context.Context, channelID string, channel domain.NotificationChannel) error {
	return nil
}

func (m *NotificationMock) DeleteNotificationChannel(ctx context.Context, channelID string) error {
	return nil
}
