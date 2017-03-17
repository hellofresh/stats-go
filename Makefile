REPO=github.com/hellofresh/stats-go

# Other config
NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m

GO_PROJECT_FILES=`go list -f '{{.Dir}}' ./... | grep -v /vendor/ | sed -n '1!p'`
GO_PROJECT_PACKAGES=`go list ./... | grep -v /vendor/`

.PHONY: all deps fmt unit-tests

all: deps unit-tests

deps:
	@echo "$(OK_COLOR)==> Installing glide with dependencies$(NO_COLOR)"
	@go get -u github.com/Masterminds/glide
	@glide install

# Format the source code
fmt:
	@gofmt -s=true -w $(GO_PROJECT_FILES)

# Run unit-tests
unit-tests:
	@go test ${GO_PROJECT_PACKAGES} -v
