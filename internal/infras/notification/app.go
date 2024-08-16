package notification

import (
	"context"
	"fmt"
)

type AppNotifier struct {
}

func NewAppNotifier() AppNotifier {
	return AppNotifier{}
}

func (n AppNotifier) SendNotification(ctx context.Context, msg string) {
	fmt.Println("Notification sent to app: ", msg)
}
