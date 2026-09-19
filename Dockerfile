# --- Base Stage ---
FROM golang:1.27-alpine AS base
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# --- Dev Stage (Con hot-reloading usando Air) ---
FROM base AS dev
RUN go install github.com/air-verse/air@latest
CMD ["air", "-c", ".air.toml"]

# --- Test Stage ---
FROM base AS test
CMD ["go", "test", "-v", "./..."]

# --- Prod Stage  ---
FROM base AS builder
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/main ./cmd/api/main.go

FROM alpine:latest AS prod
WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
