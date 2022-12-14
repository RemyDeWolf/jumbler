TARGETS = darwin/amd64 darwin/arm64 linux/amd64 linux/386 windows/amd64 windows/386

VERSION ?= dev
GITHUB_SHA ?= $(shell git rev-parse HEAD)
BUILD_TIME = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ" | tr -d '\n')
GO_VERSION = $(shell go version | awk {'print $$3'})
LDFLAGS = -s -w
PKG = github.com/remydewolf/jumbler

THIS_FILE := $(lastword $(MAKEFILE_LIST))

usage:
	@echo ""
	@echo "Task                 : Description"
	@echo "-----------------    : -------------------"
	@echo "make clean           : Remove all build files and reset assets"
	@echo "make build           : Generate build for current OS"
	@echo "make format      	: Format code"
	@echo "make load-test       : Execute load test suite"
	@echo "make run           	: Run using local code"
	@echo "make test            : Execute test suite"
	@echo "make version         : Show version"
	@echo ""

format:
	go fmt github.com/remydewolf/...

lint:
	golangci-lint run ./...

test:
	go test -race ./pkg/... ./cmd/... 

pre-commit:
	go mod tidy
	@$(MAKE) -f $(THIS_FILE) format
	@$(MAKE) -f $(THIS_FILE) lint
	@$(MAKE) -f $(THIS_FILE) test

run:
	go run -race main.go

version:
	@go run -race main.go version

build: LDFLAGS += -X $(PKG)/pkg/version.GitCommit=$(GITHUB_SHA)
build: LDFLAGS += -X $(PKG)/pkg/version.BuildTime=$(BUILD_TIME)
build: LDFLAGS += -X $(PKG)/pkg/version.GoVersion=$(GO_VERSION)
build: LDFLAGS += -X $(PKG)/pkg/version.Version=$(VERSION)
build:
	go build -race -ldflags "$(LDFLAGS)"
	@echo "You can now execute ./jumbler"

clean:
	@rm -f ./jumbler
	@rm -rf ./bin/*

release: LDFLAGS += -X $(PKG)/pkg/api.GitCommit=$(GITHUB_SHA)
release: LDFLAGS += -X $(PKG)/pkg/api.BuildTime=$(BUILD_TIME)
release: LDFLAGS += -X $(PKG)/pkg/api.GoVersion=$(GO_VERSION)
release: LDFLAGS += -X $(PKG)/pkg/api.Version=$(VERSION)
release:
	@echo "Building binaries..."
	@gox \
		-osarch "$(TARGETS)" \
		-ldflags "$(LDFLAGS)" \
		-output "./bin/jumbler_{{.OS}}_{{.Arch}}"

	@echo "Building ARM binaries..."
	GOOS=linux GOARCH=arm GOARM=5 go build -ldflags "$(LDFLAGS)" -o "./bin/jumbler_linux_arm_v5"

	@echo "Building ARM64 binaries..."
	GOOS=linux GOARCH=arm64 GOARM=7 go build -ldflags "$(LDFLAGS)" -o "./bin/jumbler_linux_arm64_v7"

	@echo "\nPackaging binaries...\n"
	@./script/package.sh

setup:
	go install github.com/mitchellh/gox@v1.0.1
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.46.2
