package driven

import (
	"context"
)

// NotificationService is an interface for sending notifications
type NotificationService interface {
	SendTransactionNotification(ctx context.Context, userId, message string) error
}
