.PHONY: all build test run

build:
	go build ./services/...

test:
	go test ./services/...

run-job-submission:
	go run ./services/job-submission

run-scheduler:
	go run ./services/scheduler

run-job-executor:
	go run ./services/job-executor
