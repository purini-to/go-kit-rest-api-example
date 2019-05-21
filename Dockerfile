# syntax = docker/dockerfile:experimental

FROM golang:1-stretch as builder

ENV GOOS=linux
ENV GOARCH=amd64
ENV GO111MODULE=auto

ENV APP_DIR /app/api
RUN mkdir -p $APP_DIR
WORKDIR $APP_DIR

RUN groupadd -r app && useradd --no-log-init -r -g app app

RUN \
    --mount=type=cache,target=/var/cache/apt/archives \
    apt-get update -y && apt-get install -y upx

RUN \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go \
    --mount=type=bind,target=. \
    CGO_ENABLED=0 go build -ldflags '-d -w -s' -o /bin/api cmd/api/main.go &&\
    upx /bin/api

FROM scratch as runner

ENV TZ=Asia/Tokyo
ENV LANG=ja_JP.UTF-8
ENV LANGUAGE=ja_JP.UTF-8
ENV LC_ALL=ja_JP.UTF-8

ENV APP_DIR /
WORKDIR $APP_DIR

COPY --from=builder /etc/passwd /etc/passwd

COPY --from=builder /bin/api $APP_DIR/api

USER app

EXPOSE 8080

ENTRYPOINT ["/api"]