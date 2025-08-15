# syntax=docker/dockerfile:1
FROM golang:1.22 as builder
WORKDIR /app
COPY . .
RUN go mod download && CGO_ENABLED=0 GOOS=linux go build -o /api ./cmd/api

FROM gcr.io/distroless/base-debian12
WORKDIR /
ENV API_PORT=3000
COPY --from=builder /api /api
COPY .env.example /.env
EXPOSE 3000 2112
ENTRYPOINT ["/api"]
