# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git ca-certificates

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source files
COPY main.go ./
COPY templates/ ./templates/
COPY static/ ./static/

# Build the binary
RUN go build -o ts-hello-world main.go

# Final stage
FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/ts-hello-world .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static

# Create state directory
RUN mkdir -p /var/lib/ts-hello-world

# Set default environment variables
ENV TS_HOSTNAME=ts-hello-world \
    TS_STATE_DIR=/var/lib/ts-hello-world \
    TS_AUTHKEY="" \
    TS_ENABLE_SERVE=true \
    TS_ENABLE_FUNNEL=false

VOLUME ["/var/lib/ts-hello-world"]

EXPOSE 443

CMD ["./ts-hello-world"]