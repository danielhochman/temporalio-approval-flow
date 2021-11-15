package activities

import (
	"context"
	"fmt"
	"go.temporal.io/sdk/activity"
)

func SendSlackNotification(ctx context.Context, person, action string) error {
	logger := activity.GetLogger(ctx)
	logger.Info(fmt.Sprintf("Notifying '%s' of request to '%s'", person, action))
	return nil
}
