# ---------- Stage 1: Builder ----------
FROM golang:1.25-alpine AS builder

LABEL author="MOHSEN_SERAJIAN"

WORKDIR /app

# Install Git (needed for go mod download)
RUN apk add --no-cache git

# -------------------------------------
# Uncomment this block if your project
# requires CGO or C libraries
# RUN apk add --no-cache \
#     build-base \
#     gcc \
#     g++ \
#     musl-dev \
#     libc6-compat \
#     make \
#     pkgconfig \
#     bash
# -------------------------------------

# Copy Go module files and download deps
COPY go.mod go.sum ./
RUN go mod download

# Copy full source
COPY . .

# Build binary (static build for Distroless)
RUN CGO_ENABLED=0 go build -o app -ldflags="-s -w" ./main.go


# ---------- Stage 2: Runtime (Distroless) ----------
FROM gcr.io/distroless/static:nonroot

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/app ./app

# Copy runtime assets (migrations, public, configs)
COPY --from=builder /app/internal/infrastructure/db/migrations ./internal/infrastructure/db/migrations
COPY --from=builder /app/public ./public
# COPY config.yaml ./config.yaml  # optional

# Run as non-root (already enforced by Distroless)
USER nonroot:nonroot

EXPOSE 8080
ENTRYPOINT ["./app"]
