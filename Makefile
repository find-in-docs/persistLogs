# This makes sure the commands are run within a BASH shell.
SHELL := /bin/bash
EXEDIR := ./bin
BIN_NAME=./bin/persistLogs

LATESTVER_WITH_UPDATE := "$(shell go list -m -u github.com/samirgadkari/sidecar | rg '.*?\s+.*?\s+\[(v.*?)\]$$' --replace '$$1')"
UPDATED_LATESTVER := "$(shell go list -m -u github.com/samirgadkari/sidecar | rg '.*?\s+(v.*?)$$' --replace '$$1')"

ifeq ($(strip $(LATESTVER_WITH_UPDATE)), "")
	LATESTVER := $(UPDATED_LATESTVER)
else
	LATESTVER := $(LATESTVER_WITH_UPDATE)
endif

# The .PHONY target will ignore any file that exists with the same name as the target
# in your makefile, and built it regardless.
.PHONY: all init build run clean

# The all target is the default target when make is called without any arguments.
all: clean | run

printvars:
	@echo "$(LATESTVER_WITH_UPDATE)"
	@echo $(UPDATED_LATESTVER)
	@echo $(LATESTVER)

init:
	go mod init github.com/samirgadkari/persistLogs
	go get github.com/samirgadkari/sidecar
	mkdir cli
	cd cli && cobra init
	cd cli && cobra add serve

${EXEDIR}:
	mkdir ${EXEDIR}

build: | ${EXEDIR}
	go get github.com/samirgadkari/sidecar@$(LATESTVER)
	go build -o ${BIN_NAME} cli/main.go

run: build
	./${BIN_NAME} serve

clean:
	go clean
	rm ${BIN_NAME}
	go clean -cache -modcache -i -r
	go get github.com/samirgadkari/sidecar@$(LATESTVER)
	go mod tidy
