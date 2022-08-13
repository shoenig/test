SHELL = bash

default: test

.PHONY: test
test:
	@echo "--> Running Tests ..."
	@go test -v -race ./...

vet:
	@echo "--> Vet Go sources ..."
	@go vet ./...

changes:
	@echo "--> Checking for source diffs ..."
	@go generate ./...
	@go mod tidy
	@go fmt ./...
	@./scripts/changes.sh
