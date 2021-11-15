package main

import (
	"context"
	"github.com/danielhochman/temporalio-approval-flow/workflow"
	"go.temporal.io/sdk/client"
	"go.uber.org/zap"
	"sync"
	"time"
)

func main() {
	c, err := client.NewClient(client.Options{})
	if err != nil {
		panic(err)
	}
	defer c.Close()

	logger, _ := zap.NewDevelopmentConfig().Build()

	// Initialize the state.
	state := &workflow.State{
		Action: "Terminate instance i-123456789abcdef0",
	}

	// Execute the workflow with the initial state.
	opts := client.StartWorkflowOptions{TaskQueue: workflow.QueueName}
	we, err := c.ExecuteWorkflow(context.Background(), opts, workflow.Workflow, state)
	if err != nil {
		panic(err)
	}

	// Start workers, one poller, one approver.
	var wg sync.WaitGroup
	logger.Info("start", zap.String("workflowID", we.GetID()))

	// This is the "frontend" polling for approval.
	wg.Add(1)
	go func() {
		defer wg.Done()

		i := 0
		for {
			resp, err := c.QueryWorkflow(context.Background(), we.GetID(), we.GetRunID(), "getState")
			if err != nil {
				panic(err)
			}

			var result workflow.State
			if err := resp.Get(&result); err != nil {
				panic(err)
			}

			logger.Info("polling...", zap.Bool("approved", result.IsApproved()))
			if result.IsApproved() {
				break
			}

			if i == 4 {
				logger.Info("we have been waiting too long, notify someone!")
				err := c.SignalWorkflow(context.Background(), we.GetID(), we.GetRunID(), workflow.NotificationChannel, &workflow.Notification{User: "jogan"})
				if err != nil {
					panic(err)
				}
			}

			time.Sleep(1 * time.Second)
			i += 1
		}
	}()

	// This is the approver.
	wg.Add(1)
	go func() {
		defer wg.Done()

		time.Sleep(10 * time.Second)

		logger.Info("approver approving via signal!")
		err := c.SignalWorkflow(
			context.Background(),
			we.GetID(), we.GetRunID(),
			workflow.CommentChannel,
			&workflow.Comment{Timestamp: time.Now(), Author: "dhochman@lyft.com", Message: "LGTM", Status: workflow.Approve})
		if err != nil {
			panic(err)
		}
	}()

	wg.Wait()
	logger.Info("workflow completed")
}
