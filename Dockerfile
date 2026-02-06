FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o main cmd/main/main.go
RUN go build -o migrate cmd/migrator/main.go
RUN go build -o cli cmd/cli/main.go

FROM alpine:latest AS runtime

WORKDIR /app

RUN apk add --no-cache curl

COPY --from=builder /app/main .
COPY --from=builder /app/migrate .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/cli .

FROM runtime AS main
ENTRYPOINT [ "/app/main" ]

FROM runtime AS migrate
ENTRYPOINT [ "/app/migrate" ]

FROM runtime AS cli
ENTRYPOINT [ "/app/cli" ]