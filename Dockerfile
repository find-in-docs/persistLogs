
# Base image
FROM golang:1.19-alpine

# Specify work directory on the image.
# All commands will refer to this work directory from now on below.
WORKDIR /app

# Copy local go.mod and go.sum files into the image
COPY go.mod ./
COPY go.sum ./

# Clean the modcache
# This is not required all the time. You should run this only
# when your modcache contains older versions that you cannot upgrade for some reason.
# RUN go clean -cache -modcache -i -r

# This should always be set. This way, instead of getting packages from Google servers,
# you get them directly from github. This way, there is no sync lag between
# your changes on github and your packages on Google servers. Sometimes, it takes
# more than a day for Google servers to catch up to your changes.
# RUN go env -w GOPROXY=direct 

# RUN apk update && \
#     apk add git

# Download required packages in the image
# RUN go mod download

# Copy source code into the image
COPY pkg/ /app/pkg/

RUN pwd
RUN ls -l pkg/main

# RUN go clean -modcache
# RUN go get github.com/spf13/viper
# RUN go get github.com/jackc/pgx/v4
# RUN go get github.com/find-in-docs/sidecar/protos/v1/messages
# RUN go get github.com/spf13/viper
# RUN go get github.com/find-in-docs/sidecar/pkg/client
# RUN go get github.com/find-in-docs/sidecar/pkg/utils
# RUN go get github.com/find-in-docs/sidecar/protos/v1/messages
# RUN go get github.com/spf13/viper
# RUN go get github.com/find-in-docs/persistLogs/pkg/config
# RUN go get github.com/find-in-docs/persistLogs/pkg/data
 
RUN go build -o persistlogs pkg/main/main.go

RUN ls -l /app