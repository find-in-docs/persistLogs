# This makes sure the commands are run within a BASH shell.
SHELL := /bin/bash
EXEDIR := ./bin
BIN_NAME=./bin/persistLogs

# The .PHONY target will ignore any file that exists with the same name as the target
# in your makefile, and built it regardless.
.PHONY: all init build run clean

# The all target is the default target when make is called without any arguments.
all: clean | run

init:
	- rm go.mod
	- rm go.sum
	go mod init github.com/find-in-docs/persistLogs
	go mod tidy -compat=1.17

${EXEDIR}:
	mkdir ${EXEDIR}

build: | ${EXEDIR}
	go build -o ${BIN_NAME} pkg/main/main.go

run: build
	./${BIN_NAME}

clean:
	go clean
	- rm ${BIN_NAME}
	go clean -cache -modcache -i -r
	go mod tidy -compat=1.17
