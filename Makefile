# This makes sure the commands are run within a BASH shell.
SHELL := /bin/bash
EXEDIR := ./bin
BIN_NAME=./bin/persistLogs

# The .PHONY target will ignore any file that exists with the same name as the target
# in your makefile, and built it regardless.
.PHONY: all init build run clean upload

# The all target is the default target when make is called without any arguments.
all: clean | run

init:
	echo "Setting up local ..."
	go env -w GOPROXY=direct 
	- rm go.mod
	- rm go.sum
	go mod init github.com/find-in-docs/persistLogs
	go mod tidy
	echo "Setting up minikube ..."
	docker build -t persistlogs -f ./Dockerfile_downloadPkgs .
	  # go env -w GOPROXY="https://proxy.golang.org,direct"

${EXEDIR}:
	echo "Building exe directory ..."
	mkdir ${EXEDIR}

build: | ${EXEDIR}
	echo ">>>>>>>>>>>>>>>>>"
	echo "  Get latest tagged version for your code (ex. sidecar) that you depend on."
	echo "  Use: go get github.com/find-in-docs/sidecar@v0.0.0-beta.10-lw (for example)."
	echo "  This will ensure your packages from github are synced with Google's servers"
	echo "  (https://proxy.golang.org). You can change this value using:"
	echo "  go env -w GOPROXY=direct"
	echo "  to get it from the github repo directly."
	echo "  Google's server might take 30 minutes to sync up with github after you request"
	echo "  your package from them. So the first time you request it after the package change occurs"
	echo "  in github, it will get it from github directly, then add your repo to their syncing process."
	echo "<<<<<<<<<<<<<<<<<"
	
	sleep 2s
	echo "Building locally ..."
	go build -o ${BIN_NAME} pkg/main/main.go

run: build
	echo "Running locally ..."
	./${BIN_NAME}

clean:
	echo "Cleaning ..."
	go clean
	- rm ${BIN_NAME}
	# go clean -cache -modcache -i -r
	go mod tidy

upload: build
	echo "Start building on minikube ..."
	# echo "Get each of these packages in the Dockerfile"
	# rg --iglob "*.go" -o -I -N "[\"]github([^\"]+)[\"]" | sed '/^$/d' | sed 's/\"//g' | awk '{print "RUN go get " $0}'
	docker build --progress=plain --no-cache -t persistlogs -f ./Dockerfile .
	kubectl run persistlogs --image=persistlogs:latest --image-pull-policy=Never --restart=Never
