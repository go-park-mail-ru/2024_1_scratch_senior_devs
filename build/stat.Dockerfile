FROM golang:1.21.0-alpine AS builder

COPY . /github.com/go-park-mail-ru/2024_1_scratch_senior_devs/
WORKDIR /github.com/go-park-mail-ru/2024_1_scratch_senior_devs/

RUN go mod download
RUN go clean --modcache
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -o ./.bin ./cmd/stat/main.go

ENV TZ="Europe/Moscow"
ENV ZONEINFO=/zoneinfo.zip

EXPOSE 8083

ENTRYPOINT ["./.bin"]