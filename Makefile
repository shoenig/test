SHELL = bash

default: test

.PHONY: test
test:
	@echo "--> Running Tests ..."
	@go test -v -race ./...

.PHONY: vet
vet:
	@echo "--> Vet Go sources ..."
	@go vet ./...

.PHONY: generate
generate:
	@echo "--> Go generate ..."
	@go generate ./...

.PHONY: changes
changes: generate
	@echo "--> Checking for source diffs ..."
	@go mod tidy
	@go fmt ./...
	@./scripts/changes.sh
