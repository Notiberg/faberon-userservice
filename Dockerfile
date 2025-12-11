# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Install PostgreSQL client for migrations
RUN apk add --no-cache postgresql-client

# Copy binary from builder
COPY --from=builder /app/main .

# Copy config file
COPY --from=builder /app/config.toml .

# Copy migrations
COPY --from=builder /app/migrations ./migrations

# Create logs directory
RUN mkdir -p /app/logs

# Expose port
EXPOSE 8080

# Create entrypoint script for migrations
RUN echo '#!/bin/sh\n\
for f in ./migrations/*.up.sql; do\n\
  echo "Running migration: $f"\n\
  psql "postgresql://$PGUSER:$PGPASSWORD@$PGHOST:$PGPORT/$POSTGRES_DB?sslmode=disable" -f "$f" || true\n\
done\n\
exec ./main\n\
' > /app/entrypoint.sh && chmod +x /app/entrypoint.sh

# Run the application with migrations
CMD ["/app/entrypoint.sh"]
