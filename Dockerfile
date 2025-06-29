# Stage 1: Build
FROM golang:1.24.4-alpine AS builder

# Ensure no root privilege escalation
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Add go mod files first for caching
COPY go.mod go.sum ./
RUN go mod download

# Then copy the source code
COPY . .

# Build statically-linked binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/port-service

# Stage 2: Minimal runtime image (distroless or scratch)
FROM gcr.io/distroless/static-debian11:nonroot

WORKDIR /

COPY --from=builder /app/main /
COPY --from=builder /app/static /static
COPY --from=builder /app/templates /templates

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/main"]
