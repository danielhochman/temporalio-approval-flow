package main

import (
	"github.com/danielhochman/temporalio-approval-flow/activities"
	"github.com/danielhochman/temporalio-approval-flow/workflow"
	temporalclient "go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	temporalworkflow "go.temporal.io/sdk/workflow"
)

func main() {
	c, err := temporalclient.NewClient(temporalclient.Options{})
	if err != nil {
		panic(err)
	}
	defer c.Close()

	w := worker.New(c, workflow.QueueName, worker.Options{})

	w.RegisterActivity(activities.SendSlackNotification)

	w.RegisterWorkflowWithOptions(workflow.Workflow, temporalworkflow.RegisterOptions{Name: workflow.WorkflowName})

	err = w.Run(worker.InterruptCh())
	if err != nil {
		panic(err)
	}
}
