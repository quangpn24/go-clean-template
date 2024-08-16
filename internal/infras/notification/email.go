package notification

import (
	"context"
	"fmt"
)

type EmailNotifier struct {
}

func NewEmailNotifier() EmailNotifier {
	return EmailNotifier{}
}

func (n EmailNotifier) SendNotification(ctx context.Context, msg string) {
	fmt.Println("Email have been sent: ", msg)
}
