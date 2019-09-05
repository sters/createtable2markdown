TOOL_COMMAND_NAME := createtable2markdown
TOOL_COMMAND_MAIN := ${TOOL_COMMAND_NAME}.go
TOOL_PACKAGE      := github.com/sters/${TOOL_COMMAND_NAME}

BUILD_DIR := ./build
BUILD_PATH := ${BUILD_DIR}/${TOOL_COMMAND_NAME}

GO_ENV := GO111MODULE=on CGO_ENABLED=0

.PHONY: init tidy test build install run
init:
	@${GO_ENV} go mod init
tidy:
	@${GO_ENV} go mod tidy
test:
	@${GO_ENV} CGO_ENABLED=1 go test -v -race -cover ./...
build:
	@${GO_ENV} go build -o ${BUILD_PATH} ${TOOL_COMMAND_MAIN}
install:
	@${GO_ENV} go install ${TOOL_PACKAGE}
run:
	@${GO_ENV} go run ${TOOL_PACKAGE}
