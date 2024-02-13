FROM --platform=linux/amd64 golang:alpine as base

MAINTAINER Pavel Simzicov <sharovik89@ya.ru>

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64 \
    APP_PATH="/home/go/src/github.com/sharovik/devbot"

WORKDIR ${APP_PATH}

#I am guessing you already already aware of distroless. It is a matter of developer taste, but distroless has been something I have fallen in love with due to security and simplicity.
COPY . .

RUN apk add --no-cache bash && apk add --no-cache make && apk add build-base && apk add --no-cache git && apk add --no-cache tzdata

RUN make build && make cleanup

FROM --platform=linux/amd64 alpine:latest as run
RUN apk --no-cache add ca-certificates

ENV APP_PATH="/home/go/src/github.com/sharovik/devbot"

WORKDIR ${APP_PATH}

COPY --from=base ${APP_PATH}/bin ${APP_PATH}/bin
COPY --from=base ${APP_PATH}/.env ${APP_PATH}/.env
COPY --from=base ${APP_PATH}/devbot.sqlite ${APP_PATH}/devbot.sqlite

# Command to run when starting the container
ENTRYPOINT ["./bin/devbot-current-system"]
