FROM golang:1.21.0-alpine AS builder

RUN apk add --update vips-dev
RUN apk add build-base
RUN apk add libwebp libwebp-tools

COPY ./wkhtmltopdf /usr/local/bin/
RUN rm ./wkhtmltopdf

COPY . /github.com/go-park-mail-ru/2024_1_scratch_senior_devs/
WORKDIR /github.com/go-park-mail-ru/2024_1_scratch_senior_devs/

RUN go mod download
RUN go clean --modcache
RUN CGO_ENABLED=1 GOOS=linux go build -mod=readonly -o ./.bin ./cmd/main/main.go

ENV TZ="Europe/Moscow"
ENV ZONEINFO=/zoneinfo.zip

EXPOSE 8080

ENTRYPOINT ["./.bin"]