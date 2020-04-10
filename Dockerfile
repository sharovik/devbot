FROM golang:alpine

MAINTAINER Pavel Simzicov <sharovik89@ya.ru>

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64 \
    APP_PATH="/home/go/src/github.com/sharovik/devbot"

WORKDIR ${APP_PATH}

COPY . .

RUN apk add --no-cache bash && apk add --no-cache make && apk add build-base

RUN make vendor
RUN make build-project-for-current-system
RUN make install

# Command to run when starting the container
ENTRYPOINT ["./bin/current-system"]