SHELL:=/usr/bin/env bash

.PHONY: frontend-build
frontend-build:
	cd frontend && rm -rf dist/ && npx parcel build

.PHONY: backend-run
backend-run: frontend-build
	go run starter.go
