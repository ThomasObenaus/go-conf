#.DEFAULT_GOAL := all
name := "go-base"

all: tools test finish

# This target (taken from: https://gist.github.com/prwhite/8168133) is an easy way to print out a usage/ help of all make targets.
# For all make targets the text after \#\# will be printed.
help: ## Prints the help
	@echo "$$(grep -hE '^\S+:.*##' $(MAKEFILE_LIST) | sed -e 's/:.*##\s*/:/' -e 's/^\(.\+\):\(.*\)/\1\:\2/' | column -c2 -t -s :)"


test: sep gen-mocks ## Runs all unittests and generates a coverage report.
	@echo "--> Run the unit-tests"
	@go test ./ -covermode=count -coverprofile=coverage.out

cover-upload: sep ## Uploads the unittest coverage to coveralls (for this the GO_BASE_COVERALLS_REPO_TOKEN has to be set correctly).
	# for this to get working you have to export the repo_token for your repo at coveralls.io
	# i.e. export GO_BASE_COVERALLS_REPO_TOKEN=<your token>
	@${GOPATH}/bin/goveralls -coverprofile=coverage.out -service=circleci -repotoken=${GO_BASE_COVERALLS_REPO_TOKEN}

gen-mocks: sep ## Generates test doubles (mocks).
	@echo "--> generate mocks (github.com/golang/mock/gomock is required for this)"
	@go get github.com/golang/mock/gomock
	@go install github.com/golang/mock/mockgen
	@mockgen -source=interfaces/provider.go -destination test/mocks/mock_provider.go

tools: sep ## Installs needed tools
	@echo "--> Install needed tools."
	@go get golang.org/x/tools/cmd/cover
	@go get github.com/mattn/goveralls

sep:
	@echo "----------------------------------------------------------------------------------"

finish:
	@echo "=================================================================================="