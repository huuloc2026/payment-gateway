# Stage 1: Build
FROM golang:1.22-alpine AS builder
WORKDIR /app

# Copy go.mod và go.sum trước để cache dependency
COPY go.mod ./
RUN go mod download


# Copy toàn bộ source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o api ./cmd/api

# Stage 2: Run
FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /app/api .
EXPOSE 8080
CMD ["./api"]
