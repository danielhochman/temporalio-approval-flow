package main

import (
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

const queueName = "two-phase-approval"

func main() {
	c, err := client.NewClient(client.Options{})
	if err != nil {
		panic(err)
	}
	defer c.Close()

	w := worker.New(c, queueName, worker.Options{})

	err = w.Run(worker.InterruptCh())
	if err != nil {
		panic(err)
	}
}
