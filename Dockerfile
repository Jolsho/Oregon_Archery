# Builder
FROM golang:1.26-alpine AS builder
WORKDIR /app/server

# Copy go.mod + go.sum
COPY server/go.mod server/go.sum ./
RUN go mod download

# Copy all Go source code for this module
COPY server/ ./

# Build binary
RUN go build -o server

# Runtime
FROM gcr.io/distroless/base-debian12
WORKDIR /app/server

# Copy binary
COPY --from=builder /app/server/server ./

# Copy ui subtree for static files
COPY ui/ /app/ui

ENV DST_DIR=/app/ui
CMD ["/app/server/server"]
