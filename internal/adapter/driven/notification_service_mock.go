package notification

import (
	"context"
	"log"
)

// NotificationServiceMock is a mock implementation of NotificationService
type NotificationServiceMock struct{}

// NewNotificationServiceMock creates a new instance of NotificationServiceMock
func NewMockNotificationService() *NotificationServiceMock {
	return &NotificationServiceMock{}
}

// SendTransactionNotification logs the notification details instead of sending an actual notification
func (m *NotificationServiceMock) SendTransactionNotification(ctx context.Context, userId, message string) error {
	log.Printf("Mock notification sent to user %s: %s\n", userId, message)
	return nil
}
