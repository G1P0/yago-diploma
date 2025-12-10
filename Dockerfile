FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache build-base

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

RUN go build -o task-manager .

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/task-manager /app/task-manager
COPY web /app/web

RUN mkdir -p /data

ENV TODO_PORT=7540
ENV TODO_DBFILE=/data/scheduler.db
# ENV TODO_PASSWORD="12345"

EXPOSE ${TODO_PORT}

CMD ["/app/task-manager"]
