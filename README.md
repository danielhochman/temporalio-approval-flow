# temporalio-approval-flow

This is a quick demo of [Temporal](https://temporal.io/) and a multiple user approval flow.

To test this demo locally, start by getting Temporal running using [temporalio/docker-compose](https://github.com/temporalio/docker-compose). I've forked the repo at [danielhochman/docker-compose](https://github.com/danielhochman/docker-compose) to use a database port for Postgres that doesn't conflict with the existing Clutch database port.

### Commands
```bash
$ mkdir -p ~/go/src/github.com/danielhochman
$ cd ~/go/src/github.com/danielhochman
$ git clone git@github.com/danielhochman/docker-compose
$ git clone git@github.com/danielhochman/temporalio-approval-flow
$ cd docker-compose
$ docker-compose up -d
$ cd ..
$ cd temporalio-approval-flow
$ go run worker/worker.go

# In another window
$ cd ~/go/src/github.com/danielhochman/temporalio-approval-flow
$ go run starter.go

# Stop docker processes
$ cd ~/go/src/github.com/danielhochman/docker-compose
$ docker-compose down
```

### Example Output
```bash
# worker/worker.go
2021/11/15 12:29:47 INFO  No logger configured for temporal client. Created default one.
2021/11/15 12:29:47 INFO  Started Worker Namespace default TaskQueue twoPhaseApproval WorkerID 927986@d594@
2021/11/15 12:30:24 DEBUG ExecuteActivity Namespace default TaskQueue twoPhaseApproval WorkerID 927986@d594@ WorkflowType Workflow WorkflowID 36672567-064a-4973-bbc6-08b17c814a93 RunID cc96fb15-96ce-4ba7-b82e-0e5f22f6a85d Attempt 1 ActivityID 9 ActivityType SendSlackNotification
2021/11/15 12:30:24 INFO  Notifying 'jogan' of request to 'Terminate instance i-123456789abcdef0' Namespace default TaskQueue twoPhaseApproval WorkerID 927986@d594@ ActivityID 9 ActivityType SendSlackNotification Attempt 1 WorkflowType Workflow WorkflowID 36672567-064a-4973-bbc6-08b17c814a93 RunID cc96fb15-96ce-4ba7-b82e-0e5f22f6a85d

# starter.go
2021-11-15T12:30:20.117-0600    INFO    temporalio-approval-flow/starter.go:35  start   {"workflowID": "36672567-064a-4973-bbc6-08b17c814a93"}
2021-11-15T12:30:20.126-0600    INFO    temporalio-approval-flow/starter.go:62  polling...      {"approved": false}
2021-11-15T12:30:21.135-0600    INFO    temporalio-approval-flow/starter.go:62  polling...      {"approved": false}
2021-11-15T12:30:22.142-0600    INFO    temporalio-approval-flow/starter.go:62  polling...      {"approved": false}
2021-11-15T12:30:23.150-0600    INFO    temporalio-approval-flow/starter.go:62  polling...      {"approved": false}
2021-11-15T12:30:24.158-0600    INFO    temporalio-approval-flow/starter.go:55  we have been waiting too long, notify someone!
2021-11-15T12:30:24.173-0600    INFO    temporalio-approval-flow/starter.go:62  polling...      {"approved": false}
2021-11-15T12:30:25.181-0600    INFO    temporalio-approval-flow/starter.go:62  polling...      {"approved": false}
2021-11-15T12:30:26.188-0600    INFO    temporalio-approval-flow/starter.go:62  polling...      {"approved": false}
2021-11-15T12:30:27.196-0600    INFO    temporalio-approval-flow/starter.go:62  polling...      {"approved": false}
2021-11-15T12:30:28.203-0600    INFO    temporalio-approval-flow/starter.go:62  polling...      {"approved": false}
2021-11-15T12:30:29.211-0600    INFO    temporalio-approval-flow/starter.go:62  polling...      {"approved": false}
2021-11-15T12:30:30.117-0600    INFO    temporalio-approval-flow/starter.go:79  approver approving via signal!
2021-11-15T12:30:30.222-0600    INFO    temporalio-approval-flow/starter.go:62  polling...      {"approved": true}
2021-11-15T12:30:30.222-0600    INFO    temporalio-approval-flow/starter.go:91  workflow completed
```

## Concepts and Terminology (TODO)
- Workflow
- Activity
- Signal
- Selector

For more info on these concepts check out:
- [Temporal docs](https://docs.temporal.io/)
- [temporalio/temporal-ecommerce](https://github.com/temporalio/temporal-ecommerce)
- [temporalio/subscription-workflow-project-template-go](https://github.com/temporalio/subscription-workflow-project-template-go)
- [temporalio/samples-go](https://github.com/temporalio/samples-go)

## TODO
- [x] Implement basic approval flow with an additional notification activity.
- [ ] Expand approval business logic (handling locks).
- [ ] Write up details of the concepts demonstrated in this example.
- [ ] Write unit tests. Testability is very important since uses of `interface{}` throughout the Temporal SDK mean that the compiler won't catch a lot of issues.
