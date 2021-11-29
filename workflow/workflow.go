package workflow

import (
	"github.com/danielhochman/temporalio-approval-flow/activities"
	"go.temporal.io/sdk/workflow"
	"time"
)

const WorkflowName = "twoPhaseApproval"
const QueueName = "twoPhaseApproval"
const ReviewChannel = "twoPhaseApprovalReviewChannel"
const NotificationChannel = "twoPhaseApprovalNotificationChannel"

type Status int64

const (
	Unspecified Status = iota
	Approve
	Lock
	Unlock
	Comment
)

func (s Status) String() string {
	switch s {
	case Approve:
		return "approve"
	case Lock:
		return "lock"
	case Unlock:
		return "unlock"
	case Comment:
		return "comment"
	default:
		return "unspecified"
	}
}

type Review struct {
	Timestamp time.Time
	Author    string
	Message   string
	Status    Status
}

type State struct {
	User    string
	Action  string
	Reviews []*Review
}

type Notification struct {
	User string
}

func (s *State) IsLocked() bool {
	locked := false
	for _, review := range s.Reviews {
		if review.Status == Lock {
			locked = true
		} else if review.Status == Unlock {
			locked = false
		}
	}
	return locked
}

func (s *State) IsApproved() bool {
	locked := false
	for _, review := range s.Reviews {
		if review.Status == Lock {
			locked = true
		} else if review.Status == Unlock {
			locked = false
		} else if !locked && review.Status == Approve{
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

	ch := workflow.GetSignalChannel(ctx, ReviewChannel)
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
			var signal Review
			c.Receive(ctx, &signal)

			state.Reviews = append(state.Reviews, &signal)
		})

		selector.Select(ctx)
		if state.IsApproved() {
			// Callback?
			break
		}
	}

	return nil
}
