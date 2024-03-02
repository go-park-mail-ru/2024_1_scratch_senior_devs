FROM golang:1.21-alpine as builder
WORKDIR /build
COPY go.mod .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /main ./cmd/main/main.go
FROM alpine:latest
COPY --from=builder main /bin/main
ENTRYPOINT ["/bin/main"]