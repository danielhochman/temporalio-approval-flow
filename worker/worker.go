package main

import (
	"github.com/danielhochman/temporalio-approval-flow/activities"
	"github.com/danielhochman/temporalio-approval-flow/workflow"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	c, err := client.NewClient(client.Options{})
	if err != nil {
		panic(err)
	}
	defer c.Close()

	w := worker.New(c, workflow.QueueName, worker.Options{})

	w.RegisterActivity(activities.SendSlackNotification)

	w.RegisterWorkflow(workflow.Workflow)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		panic(err)
	}
}
