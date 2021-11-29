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
$ make worker-run

# In another window
$ cd ~/go/src/github.com/danielhochman/temporalio-approval-flow
$ make backend-run

# Stop docker processes
$ cd ~/go/src/github.com/danielhochman/docker-compose
$ docker-compose down

# Visit localhost:9000 for the application and localhost:8088 for the Temporal dashboard.
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
