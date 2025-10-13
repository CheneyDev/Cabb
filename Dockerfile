## Multi-stage build for plane-integration
# Build stage
FROM golang:1.24-alpine AS build

WORKDIR /app

# Install build deps
RUN apk add --no-cache git ca-certificates tzdata && update-ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/server ./cmd/server

# Runtime stage
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata postgresql-client && update-ca-certificates \
    && adduser -D -u 10001 app

USER app

# App files
COPY --from=build /out/server /server
COPY db/migrations /app/db/migrations

ENV PORT=8080
EXPOSE 8080

ENTRYPOINT ["/server"]

