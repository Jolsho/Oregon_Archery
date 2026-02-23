# Builder
FROM golang:1.26-alpine AS builder
WORKDIR /app/test_server

# Copy go.mod + go.sum
COPY ui/test_server/go.mod ui/test_server/go.sum ./
RUN go mod download

# Copy all Go source code for this module
COPY ui/test_server/ ./

# Build binary
RUN go build -o server

# Runtime
FROM gcr.io/distroless/base-debian12
WORKDIR /app/test_server

# Copy binary
COPY --from=builder /app/test_server/server ./

# Copy ui subtree for static files
COPY ui/ /app/ui

ENV DST_DIR=/app/ui
EXPOSE 8080
CMD ["/app/test_server/server"]
