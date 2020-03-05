VERSION := $(shell git describe --tags --dirty --always)
BUILD := $(shell date +%FT%T%z)

LDFLAGS = -w -s -X main.buildVersion=${VERSION} -X main.buildDate=${BUILD}

.DEFAULT_GOAL := help

.PHONY: all build test help  

build: clean ## Build binary
	go build -trimpath -ldflags "${LDFLAGS}"

coverage: test ## Create coverage report
	go tool cover -func=coverage.txt
	go tool cover -html=coverage.txt

clean: ## Delete binary
	rm -f gnucash-graphql

help: ## Show Help
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-16s\033[0m %s\n", $$1, $$2}'

test: ## Run tests
	go test -coverprofile=coverage.txt -ldflags "${LDFLAGS}" ./...
