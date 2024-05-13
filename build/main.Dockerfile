FROM surnet/alpine-wkhtmltopdf:3.12-0.12.6-small as wkhtmltopdf

FROM golang:1.21.0-alpine AS builder

RUN apk add --update vips-dev
RUN apk add build-base
RUN apk add bash
RUN apk add libwebp libwebp-tools
RUN apk add libstdc++ libx11 libxrender libxext libssl1.1 fontconfig freetype ttf-dejavu ttf-droid ttf-freefont ttf-liberation && apk add --no-cache --virtual .build-deps msttcorefonts-installer && update-ms-fonts && fc-cache -f && rm -rf /tmp/* && apk del .build-deps

COPY --from=wkhtmltopdf /bin/wkhtmltopdf /bin/wkhtmltopdf

COPY . /github.com/go-park-mail-ru/2024_1_scratch_senior_devs/
WORKDIR /github.com/go-park-mail-ru/2024_1_scratch_senior_devs/

RUN go mod download
RUN go clean --modcache
RUN CGO_ENABLED=1 GOOS=linux go build -mod=readonly -o ./.bin ./cmd/main/main.go

ENV TZ="Europe/Moscow"
ENV ZONEINFO=/zoneinfo.zip

EXPOSE 8080

ENTRYPOINT ["./.bin"]