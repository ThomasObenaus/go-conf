SHELL := /bin/bash
.DEFAULT_GOAL := all

all: test format lint finish

# This target (taken from: https://gist.github.com/prwhite/8168133) is an easy way to print out a usage/ help of all make targets.
# For all make targets the text after \#\# will be printed.
help: ## Prints the help
	@echo "$$(grep -hE '^\S+:.*##' $(MAKEFILE_LIST) | sed -e 's/:.*##\s*/:/' -e 's/^\(.\+\):\(.*\)/\1\:\2/' | column -c2 -t -s :)"

test: sep gen-mocks ## Runs all unit tests and generates a coverage report.
	@echo "--> Run the unit-tests"
	@set -o pipefail ; go test $$(go list ./...) -timeout 30s -race -covermode=atomic

test.report: sep ## Runs all unittests and generates a coverage- and a test-report.
	@echo "--> Run the unit-tests"	
	@go test $$(go list ./...) -timeout 30s -race -covermode=atomic -coverprofile=coverage.out -json | tee test-report.out

run.examples: ## Runs the examples
	@echo "--> Run the example apps"
	@echo "---> Simple Example"
	@go run ./examples/simple
	@echo ""
	@echo "---> Readme Example"
	@go run ./examples/readme
	@echo ""
	@echo "---> Multilevel Example"
	@go run ./examples/multilevel
	@echo ""
	@echo "---> MappingFunc Example"
	@go run ./examples/mapfun

lint: sep ## Runs the linter to check for coding-style issues
	@echo "--> Lint project"
	@echo "!!!!golangci-lint has to be installed. See: https://github.com/golangci/golangci-lint#install"
	@golangci-lint run --enable gofmt

gen-mocks: sep ## Generates test doubles (mocks).
	@echo "--> generate mocks (github.com/golang/mock/gomock is required for this)"
	@go get github.com/golang/mock/gomock
	@go install github.com/golang/mock/mockgen
	@mockgen -source=interfaces/provider.go -destination test/mocks/mock_provider.go

format: ## Formats the code using gofmt
	@echo "--> Formatting all sources using go fmt"
	@gofmt -w -s .

sep:
	@echo "----------------------------------------------------------------------------------"

finish:
	@echo "=================================================================================="