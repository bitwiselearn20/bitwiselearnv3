# syntax=docker/dockerfile:1

# ---- Build stage ----
FROM golang:1.26-alpine AS builder

WORKDIR /src

# Cache modules separately from source for fast rebuilds.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Static, stripped binary — no libc dependency so it runs on scratch/distroless.
# CGO_ENABLED=0 + -ldflags "-s -w" yields the smallest possible image.
ARG TARGETOS=linux
ARG TARGETARCH=amd64
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath -ldflags "-s -w" -o /out/api ./cmd/api && \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath -ldflags "-s -w" -o /out/worker ./cmd/worker

# ---- Runtime stage (distroless, non-root) ----
FROM gcr.io/distroless/static-debian12:nonroot AS api

WORKDIR /app
COPY --from=builder /out/api /app/api

# distroless:nonroot already runs as uid 65532.
EXPOSE 8080
ENV PORT=8080

# Liveness handled by Azure Container Apps probe against /health.
ENTRYPOINT ["/app/api"]

# ---- Worker image: `docker build --target worker` ----
FROM gcr.io/distroless/static-debian12:nonroot AS worker
WORKDIR /app
COPY --from=builder /out/worker /app/worker
ENTRYPOINT ["/app/worker"]
