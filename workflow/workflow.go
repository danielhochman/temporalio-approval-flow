package workflow

import (
	"github.com/danielhochman/temporalio-approval-flow/activities"
	"go.temporal.io/sdk/workflow"
	"time"
)

const QueueName = "twoPhaseApproval"
const CommentChannel = "twoPhaseApprovalCommentChannel"
const NotificationChannel = "twoPhaseApprovalNotificationChannel"

type Status int64

const (
	Unspecified Status = iota
	Approve
	Lock
	Unlock
)

type Comment struct {
	Timestamp time.Time
	Author string
	Message string
	Status Status
}

type State struct {
	Action string
	Comments []*Comment
}

type Notification struct {
	User string
}

func (s *State) IsApproved() bool {
	for _, comment := range s.Comments {
		if comment.Status == Approve {
			return true
		}
	}
	return false
}

func Workflow(ctx workflow.Context, state *State) error {
	err := workflow.SetQueryHandler(ctx, "getState", func() (*State, error) {
		return state, nil
	})
	if err != nil {
		return err
	}

	notifChan := workflow.GetSignalChannel(ctx, NotificationChannel)

	ch := workflow.GetSignalChannel(ctx, CommentChannel)
	for {
		selector := workflow.NewSelector(ctx)

		selector.AddReceive(notifChan, func(c workflow.ReceiveChannel, _ bool) {
			var signal Notification
			c.Receive(ctx, &signal)

			workflow.ExecuteActivity(
				workflow.WithActivityOptions(ctx, workflow.ActivityOptions{StartToCloseTimeout: time.Minute}),
				activities.SendSlackNotification, signal.User, state.Action)
		})

		selector.AddReceive(ch, func(c workflow.ReceiveChannel, _ bool) {
			var signal Comment
			c.Receive(ctx, &signal)

			state.Comments = append(state.Comments, &signal)
		})

		selector.Select(ctx)
		if state.IsApproved() {
			// Callback?
			break
		}
	}

	return nil
}
