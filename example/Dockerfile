# Build stage
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git upx

# Set working directory
WORKDIR /workspace

# Copy the entire workspace (including xcomp module)
COPY . .

# Change to example directory
WORKDIR /workspace/example

# Download dependencies
RUN go mod download

# Build optimized binary with UPX compression
ARG VERSION=docker
ARG BUILD_TIME
ARG GIT_COMMIT=unknown

# Build ultra-optimized static binary for Linux
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags "-s -w -extldflags=-static -X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.GitCommit=${GIT_COMMIT}" \
    -trimpath -a -installsuffix cgo \
    -o api-server .

# Compress with UPX for maximum size reduction
RUN upx --lzma --best api-server

# Final stage - distroless for security and minimal size
FROM gcr.io/distroless/static:nonroot

# Copy binary from builder
COPY --from=builder /workspace/example/api-server /api-server

# Copy config files
COPY --from=builder /workspace/example/config-prod.yaml /config-prod.yaml

# Use nonroot user for security
USER nonroot:nonroot

# Expose port
EXPOSE 3000

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/api-server", "health"]

# Default command
ENTRYPOINT ["/api-server"]
CMD ["serve", "--config", "/config-prod.yaml"]
