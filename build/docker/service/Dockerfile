FROM golang:1.24.3-alpine AS builder

WORKDIR /app

ARG CMD_PATH=main.go
ARG SERVICE_NAME

COPY ./go.mod ./go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download -x && go mod verify

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s -extldflags '-static'" \
    -trimpath \
    -buildvcs=false \
    -o /app/service "$CMD_PATH"

FROM alpine:3.21

WORKDIR /app

COPY --from=builder /app/service /app/service
